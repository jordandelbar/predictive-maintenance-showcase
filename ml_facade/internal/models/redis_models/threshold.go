package redis_models

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

type Threshold struct {
	MachineID int     `json:"machine_id" redis:"machine_id"`
	Threshold float64 `json:"threshold" redis:"threshold"`
}

type ThresholdModel struct {
	RedisDB *redis.Pool
}

func (t *ThresholdModel) Insert(threshold Threshold) error {
	conn := t.RedisDB.Get()
	defer conn.Close()

	// Initialize the threshold value
	_, err := conn.Do("SET", fmt.Sprintf("threshold:%v", threshold.MachineID), threshold.Threshold)
	if err != nil {
		return err
	}

	// Initialize the anomaly counter
	_, err = conn.Do("SET", fmt.Sprintf("anomaly_counter:%v", threshold.MachineID), 0)
	if err != nil {
		return err
	}
	return nil
}

func (t *ThresholdModel) Get(id int) (float64, error) {
	conn := t.RedisDB.Get()
	defer conn.Close()

	threshold, err := redis.Float64(conn.Do("GET", fmt.Sprintf("threshold:%v", id)))
	if err != nil {
		return 0., err
	}
	return threshold, nil
}

func (t *ThresholdModel) Increment(id int) (int, error) {
	conn := t.RedisDB.Get()
	defer conn.Close()

	script := redis.NewScript(1, `
		local key = KEYS[1]
		local currentCounter = tonumber(redis.call("GET", key) or "0")
		if currentCounter < 20 then
			currentCounter = redis.call("INCR", key)
		end
		return currentCounter
	`)

	anomalyCounter, err := redis.Int(script.Do(conn, fmt.Sprintf("anomaly_counter:%v", id)))
	if err != nil {
		return 0, err
	}

	return anomalyCounter, nil
}

func (t *ThresholdModel) Decrement(id int) (int, error) {
	conn := t.RedisDB.Get()
	defer conn.Close()

	// Lua script to decrement the counter if it's above 0
	script := redis.NewScript(1, `
		local key = KEYS[1]
		local currentCounter = tonumber(redis.call("GET", key) or "0")
		if currentCounter > 0 then
			currentCounter = redis.call("DECR", key)
		end
		return currentCounter
	`)

	anomalyCounter, err := redis.Int(script.Do(conn, fmt.Sprintf("anomaly_counter:%v", id)))
	if err != nil {
		return 0, err
	}

	return anomalyCounter, nil
}
