package data

import "database/sql"

type Models struct {
	Sensor SensorModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Sensor: SensorModel{DB: db},
	}
}
