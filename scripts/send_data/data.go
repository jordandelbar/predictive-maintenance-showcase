package main

type RequestPerSecond struct {
	rate      int
	rateBurst int
}

type SensorDataPayload struct {
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

// SensorData represents the data to be sent for prediction
type SensorData struct {
	SensorDataPayload
	MachineStatus string `json:"machine_status"`
}
