package main

import (
	"flag"
	"ml_facade/app"
	"ml_facade/config"
	"os"
	"time"
)

func main() {
	var cfg config.Config

	flag.IntVar(&cfg.Port, "port", 4000, "API server port")
	flag.StringVar(&cfg.Env, "env", os.Getenv("ENVIRONMENT"), "Environment (development|staging|production)")
	flag.StringVar(&cfg.Db.Host, "db-host", os.Getenv("MONITORING_DB_HOST"), "PostgreSQL Host")
	flag.StringVar(&cfg.Db.Port, "db-port", os.Getenv("MONITORING_DB_PORT"), "PostgreSQL Port")
	flag.StringVar(&cfg.Db.Username, "db-username", os.Getenv("MONITORING_DB_USERNAME"), "PostgreSQL Username")
	flag.StringVar(&cfg.Db.Password, "db-password", os.Getenv("MONITORING_DB_PASSWORD"), "PostgreSQL Password")
	flag.StringVar(&cfg.Db.DatabaseName, "db-dbname", os.Getenv("MONITORING_DB_NAME"), "PostgreSQL Name")
	flag.StringVar(&cfg.Rdb.Host, "rdb-host", os.Getenv("REDIS_DB_HOST"), "Redis Host")
	flag.StringVar(&cfg.Rdb.Port, "rdb-port", os.Getenv("REDIS_DB_PORT"), "Redis Port")
	flag.StringVar(&cfg.MlService.Host, "ml-service-host", os.Getenv("ML_SERVICE_HOST"), "ML Service Host")
	flag.StringVar(&cfg.MlService.Port, "ml-service-port", os.Getenv("ML_SERVICE_PORT"), "ML Service Port")
	flag.IntVar(&cfg.Db.MaxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.Db.MaxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.DurationVar(&cfg.Db.MaxIdleTime, "db-max-idle-time", 15*time.Minute, "PostgreSQL max connection idle time")

	flag.Parse()

	app.StartApp(cfg)
	os.Exit(1)
}
