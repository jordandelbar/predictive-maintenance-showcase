package config

import "time"

type Config struct {
	Env  string
	Port int
	Db   struct {
		Dsn          string
		MaxOpenConns int
		MaxIdleConns int
		MaxIdleTime  time.Duration
	}
	Rdb struct {
		Uri string
	}
}
