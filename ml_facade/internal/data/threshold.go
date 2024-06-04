package data

import (
	"context"
	"github.com/redis/go-redis/v9"
	"strconv"
)

type Threshold struct {
	MachineID int     `json:"machine_id"`
	Threshold float64 `json:"threshold"`
}

type ThresholdModel struct {
	Rdb *redis.Client
}

func (t ThresholdModel) Insert(threshold Threshold) error {
	var ctx = context.Background()

	return t.Rdb.Set(ctx, strconv.Itoa(threshold.MachineID), threshold.Threshold, 0).Err()
}

func (t ThresholdModel) Get(id int) (float64, error) {
	var ctx = context.Background()

	thresholdStr, err := t.Rdb.Get(ctx, strconv.Itoa(id)).Result()
	threshold, err := strconv.ParseFloat(thresholdStr, 32)
	return threshold, err
}
