package main

import (
	"flag"
	"ml_facade/cmd/app"
	"ml_facade/config"
	"os"
	"time"
)

func main() {
	var cfg config.Config

	flag.StringVar(&cfg.Env, "env", os.Getenv("ENVIRONMENT"), "Environment (development|staging|production)")
	flag.IntVar(&cfg.ApiServer.Port, "port", 4000, "API server port")
	flag.IntVar(&cfg.ApiServer.Limiter.Rps, "rate-limiter", 500, "Rate limiter")
	flag.IntVar(&cfg.ApiServer.Limiter.Burst, "rate-limiter-burst", 20, "Rate limiter burst")
	flag.BoolVar(&cfg.ApiServer.Limiter.Enabled, "rate-limiter-enabled", true, "Enable rate limiter")
	flag.StringVar(&cfg.PostgresDB.Host, "db-host", os.Getenv("MONITORING_DB_HOST"), "PostgreSQL Host")
	flag.StringVar(&cfg.PostgresDB.Port, "db-port", os.Getenv("MONITORING_DB_PORT"), "PostgreSQL Port")
	flag.StringVar(&cfg.PostgresDB.Username, "db-username", os.Getenv("MONITORING_DB_USERNAME"), "PostgreSQL Username")
	flag.StringVar(&cfg.PostgresDB.Password, "db-password", os.Getenv("MONITORING_DB_PASSWORD"), "PostgreSQL Password")
	flag.StringVar(&cfg.PostgresDB.DatabaseName, "db-dbname", os.Getenv("MONITORING_DB_NAME"), "PostgreSQL Name")
	flag.IntVar(&cfg.PostgresDB.MaxOpenConns, "db-max-open-conns", 75, "PostgreSQL max open connections")
	flag.IntVar(&cfg.PostgresDB.MaxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.DurationVar(&cfg.PostgresDB.MaxIdleTime, "db-max-idle-time", 15*time.Minute, "PostgreSQL max connection idle time")
	flag.StringVar(&cfg.RedisDB.Host, "rdb-host", os.Getenv("REDIS_DB_HOST"), "Redis Host")
	flag.StringVar(&cfg.RedisDB.Port, "rdb-port", os.Getenv("REDIS_DB_PORT"), "Redis Port")
	flag.StringVar(&cfg.MlService.Host, "ml-service-host", os.Getenv("ML_SERVICE_HOST"), "ML Service Host")
	flag.StringVar(&cfg.MlService.Port, "ml-service-port", os.Getenv("ML_SERVICE_PORT"), "ML Service Port")
	flag.StringVar(&cfg.RabbitMQConsumer.URI, "rabbitmq-uri", os.Getenv("RABBITMQ_URI"), "RabbitMQ URI")
	flag.StringVar(&cfg.RabbitMQConsumer.Queue, "rabbitmq-queue", os.Getenv("RABBITMQ_QUEUE"), "RabbitMQ Queue")
	flag.IntVar(&cfg.RabbitMQConsumer.NumWorkers, "rabbitmq-workers", 100, "RabbitMQ number of workers")
	flag.IntVar(&cfg.RabbitMQConsumer.BatchSize, "rabbitmq-batchsize", 20, "RabbitMQ batch size")
	flag.DurationVar(&cfg.RabbitMQConsumer.BatchTimeout, "rabbitmq-batchtimeout", 50*time.Millisecond, "RabbitMQ batch timeout")

	flag.Parse()

	app.StartApp(cfg)
	os.Exit(1)
}
