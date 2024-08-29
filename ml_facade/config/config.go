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

type CfgLimiter struct {
	Rps     int
	Burst   int
	Enabled bool
}

type CfgPostgresDB struct {
	Host         string
	Port         string
	Username     string
	Password     string
	DatabaseName string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  time.Duration
}

type CfgRedisDB struct {
	Host string
	Port string
}

type CfgMlService struct {
	Host string
	Port string
}

type Config struct {
	Env              string
	Port             int
	Limiter          CfgLimiter
	PostgresDB       CfgPostgresDB
	RedisDB          CfgRedisDB
	MlService        CfgMlService
	RabbitMQConsumer CfgRabbitMQConsumer
}

func (c CfgPostgresDB) PostgresDBDsn() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.Username, c.Password, c.Host, c.Port, c.DatabaseName)
}

func (c CfgRedisDB) RedisDBDsn() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

func (c CfgMlService) MlServiceUri() string {
	return fmt.Sprintf("http://%s:%s", c.Host, c.Port)
}
