package data

import (
	"database/sql"

	"github.com/gomodule/redigo/redis"
)

type Models struct {
	Sensor    SensorModel
	Threshold ThresholdModel
}

func NewModels(db *sql.DB, rdb *redis.Pool) Models {
	return Models{
		Sensor:    SensorModel{DB: db},
		Threshold: ThresholdModel{Rdb: rdb},
	}
}
