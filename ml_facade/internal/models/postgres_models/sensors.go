package postgres_models

import (
	"database/sql"
)

type Sensor struct {
	MachineID int     `json:"machine_id"`
	Sensor00  float64 `json:"sensor_00"`
	Sensor01  float64 `json:"sensor_01"`
	Sensor02  float64 `json:"sensor_02"`
	Sensor03  float64 `json:"sensor_03"`
	Sensor04  float64 `json:"sensor_04"`
	Sensor05  float64 `json:"sensor_05"`
	Sensor06  float64 `json:"sensor_06"`
	Sensor07  float64 `json:"sensor_07"`
	Sensor08  float64 `json:"sensor_08"`
	Sensor09  float64 `json:"sensor_09"`
	Sensor10  float64 `json:"sensor_10"`
	Sensor11  float64 `json:"sensor_11"`
	Sensor12  float64 `json:"sensor_12"`
	Sensor13  float64 `json:"sensor_13"`
	Sensor14  float64 `json:"sensor_14"`
	Sensor15  float64 `json:"sensor_15"`
	Sensor16  float64 `json:"sensor_16"`
	Sensor17  float64 `json:"sensor_17"`
	Sensor18  float64 `json:"sensor_18"`
	Sensor19  float64 `json:"sensor_19"`
	Sensor20  float64 `json:"sensor_20"`
	Sensor21  float64 `json:"sensor_21"`
	Sensor22  float64 `json:"sensor_22"`
	Sensor23  float64 `json:"sensor_23"`
	Sensor24  float64 `json:"sensor_24"`
	Sensor25  float64 `json:"sensor_25"`
	Sensor26  float64 `json:"sensor_26"`
	Sensor27  float64 `json:"sensor_27"`
	Sensor28  float64 `json:"sensor_28"`
	Sensor29  float64 `json:"sensor_29"`
	Sensor30  float64 `json:"sensor_30"`
	Sensor31  float64 `json:"sensor_31"`
	Sensor32  float64 `json:"sensor_32"`
	Sensor33  float64 `json:"sensor_33"`
	Sensor34  float64 `json:"sensor_34"`
	Sensor35  float64 `json:"sensor_35"`
	Sensor36  float64 `json:"sensor_36"`
	Sensor37  float64 `json:"sensor_37"`
	Sensor38  float64 `json:"sensor_38"`
	Sensor39  float64 `json:"sensor_39"`
	Sensor40  float64 `json:"sensor_40"`
	Sensor41  float64 `json:"sensor_41"`
	Sensor42  float64 `json:"sensor_42"`
	Sensor43  float64 `json:"sensor_43"`
	Sensor44  float64 `json:"sensor_44"`
	Sensor45  float64 `json:"sensor_45"`
	Sensor46  float64 `json:"sensor_46"`
	Sensor47  float64 `json:"sensor_47"`
	Sensor48  float64 `json:"sensor_48"`
	Sensor49  float64 `json:"sensor_49"`
	Sensor50  float64 `json:"sensor_50"`
	Sensor51  float64 `json:"sensor_51"`
}

type SensorModel struct {
	PostgresDB *sql.DB
}

func (s *SensorModel) Insert(record *Record) error {
	query := `
		INSERT INTO monitoring (
			machine_id, sensor_00, sensor_01, sensor_02, sensor_03, sensor_04, sensor_05, sensor_06, sensor_07,
			sensor_08, sensor_09, sensor_10, sensor_11, sensor_12, sensor_13, sensor_14, sensor_15, sensor_16,
			sensor_17, sensor_18, sensor_19, sensor_20, sensor_21, sensor_22, sensor_23, sensor_24, sensor_25,
			sensor_26, sensor_27, sensor_28, sensor_29, sensor_30, sensor_31, sensor_32, sensor_33, sensor_34,
			sensor_35, sensor_36, sensor_37, sensor_38, sensor_39, sensor_40, sensor_41, sensor_42, sensor_43,
			sensor_44, sensor_45, sensor_46, sensor_47, sensor_48, sensor_49, sensor_50, sensor_51,
			reconstruction_error, anomaly, anomaly_counter, origin
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20,
			$21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38,
			$39, $40, $41, $42, $43, $44, $45, $46, $47, $48, $49, $50, $51, $52, $53, $54, $55, $56,
			$57
		) RETURNING id, created_at`

	args := []any{
		record.SensorData.MachineID, record.SensorData.Sensor00, record.SensorData.Sensor01, record.SensorData.Sensor02,
		record.SensorData.Sensor03, record.SensorData.Sensor04, record.SensorData.Sensor05, record.SensorData.Sensor06,
		record.SensorData.Sensor07, record.SensorData.Sensor08, record.SensorData.Sensor09, record.SensorData.Sensor10,
		record.SensorData.Sensor11, record.SensorData.Sensor12, record.SensorData.Sensor13, record.SensorData.Sensor14,
		record.SensorData.Sensor15, record.SensorData.Sensor16, record.SensorData.Sensor17, record.SensorData.Sensor18,
		record.SensorData.Sensor19, record.SensorData.Sensor20, record.SensorData.Sensor21, record.SensorData.Sensor22,
		record.SensorData.Sensor23, record.SensorData.Sensor24, record.SensorData.Sensor25, record.SensorData.Sensor26,
		record.SensorData.Sensor27, record.SensorData.Sensor28, record.SensorData.Sensor29, record.SensorData.Sensor30,
		record.SensorData.Sensor31, record.SensorData.Sensor32, record.SensorData.Sensor33, record.SensorData.Sensor34,
		record.SensorData.Sensor35, record.SensorData.Sensor36, record.SensorData.Sensor37, record.SensorData.Sensor38,
		record.SensorData.Sensor39, record.SensorData.Sensor40, record.SensorData.Sensor41, record.SensorData.Sensor42,
		record.SensorData.Sensor43, record.SensorData.Sensor44, record.SensorData.Sensor45, record.SensorData.Sensor46,
		record.SensorData.Sensor47, record.SensorData.Sensor48, record.SensorData.Sensor49, record.SensorData.Sensor50,
		record.SensorData.Sensor51, record.ModelResponse.ReconstructionError, record.Anomaly, record.AnomalyCounter,
		record.Origin,
	}
	return s.PostgresDB.QueryRow(query, args...).Scan(&record.ID, &record.CreatedAt)
}
