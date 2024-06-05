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
	_, err := conn.Do("HMSET", fmt.Sprintf("threshold:%v", threshold.MachineID), "threshold", threshold.Threshold)
	if err != nil {
		return err
	}
	return nil
}

func (t ThresholdModel) Get(id int) (float64, error) {
	conn := t.Rdb.Get()
	threshold, err := redis.Float64(conn.Do("HGETALL", fmt.Sprintf("threshold:%v", id)))
	if err != nil {
		return 0., err
	}
	return threshold, nil
}
