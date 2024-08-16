package app

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gomodule/redigo/redis"
	_ "github.com/lib/pq"
	"log/slog"
	"ml_facade/config"
	"ml_facade/internal/api"
	"ml_facade/internal/consumer"
	"ml_facade/internal/models/postgres_models"
	"ml_facade/internal/models/redis_models"
	"ml_facade/internal/service"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const version = "0.0.1"

type application struct {
	config         config.Config
	logger         *slog.Logger
	server         *api.Server
	rabbitConsumer *consumer.RabbitMQConsumer
}

func StartApp(cfg config.Config) {
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	pdb, err := postgresDB(cfg)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to connect to postgres_models: %s", err))
		os.Exit(1)
	}
	logger.Info("postgresql database connection pool established")

	rdb, err := redisDB(cfg)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to connect to redis database: %s", err))
		os.Exit(1)
	}
	logger.Info("redis database connection pool established")

	thresholdModel := redis_models.ThresholdModel{RedisDB: rdb}
	sensorModel := postgres_models.SensorModel{PostgresDB: pdb}

	mlService, err := service.NewMlService(cfg, logger, &sensorModel, &thresholdModel, &wg)
	if err != nil {
		logger.Error(fmt.Sprintf("error creating ml service client: %v", err))
		os.Exit(1)
	}
	logger.Info("ml service client successfully initialized")

	server := api.NewApiServer(cfg, logger, mlService, &thresholdModel, version, &wg)
	rabbitmqConsumer := consumer.NewRabbitMQConsumer(cfg, logger, mlService, &wg)

	app := &application{
		config:         cfg,
		logger:         logger,
		server:         server,
		rabbitConsumer: rabbitmqConsumer,
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := app.server.Serve(ctx); err != nil {
			app.logger.Error("api server error: %v", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := app.rabbitConsumer.Serve(ctx); err != nil {
			app.logger.Error("rabbit consumer error: %v", err)
		}
	}()

	sig := <-signalChan
	logger.Info(fmt.Sprintf("received signal %v. Initiating shutdown...", sig))

	cancel()

	wg.Wait()

	logger.Info("application shutdown completed")
}

func postgresDB(cfg config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.PostgresDBDsn())
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.PostgresDB.MaxOpenConns)
	db.SetMaxIdleConns(cfg.PostgresDB.MaxIdleConns)
	db.SetConnMaxIdleTime(cfg.PostgresDB.MaxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		defer db.Close()
		return nil, err
	}

	return db, nil
}

func redisDB(cfg config.Config) (*redis.Pool, error) {
	rdb := &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", cfg.RedisDBDsn())
		},
	}
	conn := rdb.Get()
	defer conn.Close()

	_, err := redis.String(conn.Do("PING"))
	if err != nil {
		return nil, err
	}
	return rdb, nil
}
