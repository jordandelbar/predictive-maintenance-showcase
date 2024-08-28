package postgres_models

import (
	"time"
)

type MlServiceResponse struct {
	ReconstructionErrors []float64 `json:"reconstruction_errors"`
}

type Record struct {
	ID                  int64
	CreatedAt           time.Time
	SensorData          Sensor
	ReconstructionError float64
	Anomaly             bool
	AnomalyCounter      int
	Origin              string
}
