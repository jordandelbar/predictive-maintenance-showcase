package data

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

type Threshold struct {
	MachineID int     `json:"machine_id" redis:"machine_id"`
	Threshold float64 `json:"threshold" redis:"threshold"`
}

type ThresholdModel struct {
	Rdb *redis.Pool
}

func (t ThresholdModel) Insert(threshold Threshold) error {
	conn := t.Rdb.Get()
	defer conn.Close()

	_, err := conn.Do("SET", fmt.Sprintf("threshold:%v", threshold.MachineID), threshold.Threshold)
	if err != nil {
		return err
	}
	return nil
}

func (t ThresholdModel) Get(id int) (float64, error) {
	conn := t.Rdb.Get()
	defer conn.Close()

	threshold, err := redis.Float64(conn.Do("GET", fmt.Sprintf("threshold:%v", id)))
	if err != nil {
		return 0., err
	}
	return threshold, nil
}

func (t ThresholdModel) Increment(id int) (int, error) {
	conn := t.Rdb.Get()
	defer conn.Close()

	counter, err := redis.Int(conn.Do("INCR", fmt.Sprintf("anomaly_counter:%v", id)))
	if err != nil {
		return 0, err
	}
	return counter, nil
}

func (t ThresholdModel) Decrement(id int) (int, error) {
	conn := t.Rdb.Get()
	defer conn.Close()

	currentCounter, err := redis.Int(conn.Do("GET", fmt.Sprintf("anomaly_counter:%v", id)))
	if err != nil {
		return 0, nil
	}

	var anomalyCounter int

	if currentCounter > 0 {
		anomalyCounter, err = redis.Int(conn.Do("DECR", fmt.Sprintf("anomaly_counter:%v", id)))
		if err != nil {
			return 0, err
		}
	} else {
		anomalyCounter = currentCounter
	}

	return anomalyCounter, nil
}
