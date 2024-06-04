package data

import "time"

type ModelResponse struct {
	ReconstructionError float64 `json:"reconstruction_error"`
}

type Record struct {
	ID            int64
	CreatedAt     time.Time
	SensorData    Sensor
	ModelResponse ModelResponse
	Anomaly       bool
}
