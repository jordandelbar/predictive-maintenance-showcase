package app

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"ml_facade/config"
	"ml_facade/internal/data"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

const version = "0.0.1"

type application struct {
	config        config.Config
	logger        *slog.Logger
	models        data.Models
	mlModelClient http.Client
}

func StartApp(cfg config.Config) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	db, err := openDB(cfg)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	transport := &http.Transport{
		MaxIdleConns:        50,
		IdleConnTimeout:     10 * time.Second,
		MaxIdleConnsPerHost: 10,
	}

	client := http.Client{Transport: transport, Timeout: time.Second * 10}

	logger.Info("database connection pool established")

	app := &application{
		config:        cfg,
		logger:        logger,
		models:        data.NewModels(db),
		mlModelClient: client,
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.Info("starting server", "addr", srv.Addr, "env", cfg.Env)

	err = srv.ListenAndServe()
	logger.Error(err.Error())
}

func openDB(cfg config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.Db.Dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.Db.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Db.MaxIdleConns)
	db.SetConnMaxIdleTime(cfg.Db.MaxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
