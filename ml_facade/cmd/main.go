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
	flag.StringVar(&cfg.Env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.Db.Dsn, "db-dsn", os.Getenv("MONITORING_DB_DSN"), "PostgreSQL DSN")
	flag.IntVar(&cfg.Db.MaxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.Db.MaxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.DurationVar(&cfg.Db.MaxIdleTime, "db-max-idle-time", 15*time.Minute, "PostgreSQL max connection idle time")

	flag.Parse()

	app.StartApp(cfg)
	os.Exit(1)
}
