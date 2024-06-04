package data

import (
	"database/sql"
	"github.com/redis/go-redis/v9"
)

type Models struct {
	Sensor    SensorModel
	Threshold ThresholdModel
}

func NewModels(db *sql.DB, rdb *redis.Client) Models {
	return Models{
		Sensor:    SensorModel{DB: db},
		Threshold: ThresholdModel{Rdb: rdb},
	}
}
