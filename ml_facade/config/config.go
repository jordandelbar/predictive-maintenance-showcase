package config

import (
	"fmt"
	"time"
)

type CfgRabbitMQConsumer struct {
	URI          string
	Queue        string
	NumWorkers   int
	BatchSize    int
	BatchTimeout time.Duration
}

type Config struct {
	Env     string
	Port    int
	Limiter struct {
		Rps     int
		Burst   int
		Enabled bool
	}
	PostgresDB struct {
		Host         string
		Port         string
		Username     string
		Password     string
		DatabaseName string
		MaxOpenConns int
		MaxIdleConns int
		MaxIdleTime  time.Duration
	}
	RedisDB struct {
		Host string
		Port string
	}
	MlService struct {
		Host string
		Port string
	}
	CfgRabbitMQConsumer
}

func (c Config) PostgresDBDsn() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.PostgresDB.Username, c.PostgresDB.Password, c.PostgresDB.Host, c.PostgresDB.Port, c.PostgresDB.DatabaseName)
}

func (c Config) RedisDBDsn() string {
	return fmt.Sprintf("%s:%s", c.RedisDB.Host, c.RedisDB.Port)
}

func (c Config) MlServiceUri() string {
	return fmt.Sprintf("http://%s:%s", c.MlService.Host, c.MlService.Port)
}
