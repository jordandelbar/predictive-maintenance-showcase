package config

import (
	"fmt"
	"time"
)

type Config struct {
	Env  string
	Port int
	Db   struct {
		Host         string
		Port         string
		Username     string
		Password     string
		DatabaseName string
		MaxOpenConns int
		MaxIdleConns int
		MaxIdleTime  time.Duration
	}
	Rdb struct {
		Host string
		Port string
	}
	MlService struct {
		Host string
		Port string
	}
}

func (c Config) DbDsn() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.Db.Username, c.Db.Password, c.Db.Host, c.Db.Port, c.Db.DatabaseName)
}

func (c Config) RdbDsn() string {
	return fmt.Sprintf("%s:%s", c.Rdb.Host, c.Rdb.Port)
}

func (c Config) MlServiceDsn() string {
	return fmt.Sprintf("http://%s:%s", c.MlService.Host, c.MlService.Port)
}
