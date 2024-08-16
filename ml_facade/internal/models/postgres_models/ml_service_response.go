package postgres_models

import (
	"time"
)

type MlServiceResponse struct {
	ReconstructionError float64 `json:"reconstruction_error"`
}

type Record struct {
	ID             int64
	CreatedAt      time.Time
	SensorData     Sensor
	ModelResponse  MlServiceResponse
	Anomaly        bool
	AnomalyCounter int
	Origin         string
}
