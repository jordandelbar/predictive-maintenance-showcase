package service

import (
	"errors"
	"fmt"
	retryablehttp "github.com/hashicorp/go-retryablehttp"
	"log/slog"
	"ml_facade/config"
	"ml_facade/internal/models/postgres_models"
	"ml_facade/internal/models/redis_models"
	"net/http"
	"sync"
	"time"
)

type cacheEntry struct {
	value      float64
	expiration time.Time
}

type MlService struct {
	client         *retryablehttp.Client
	sensorModel    *postgres_models.SensorModel
	thresholdModel *redis_models.ThresholdModel
	config         config.Config
	logger         *slog.Logger
	thresholdCache sync.Map
	wg             *sync.WaitGroup
	clientReady    chan struct{}
	clientErr      chan error
}

func NewMlService(
	cfg config.Config,
	logger *slog.Logger,
	sensorModel *postgres_models.SensorModel,
	thresholdModel *redis_models.ThresholdModel,
	wg *sync.WaitGroup) (*MlService, error) {
	mlService := &MlService{
		sensorModel:    sensorModel,
		thresholdModel: thresholdModel,
		config:         cfg,
		logger:         logger,
		wg:             wg,
		clientReady:    make(chan struct{}),
		clientErr:      make(chan error, 1),
	}

	go mlService.initClientWithRetry()

	select {
	case <-mlService.clientReady:
		return mlService, nil
	case err := <-mlService.clientErr:
		return nil, err
	}
}

func (m *MlService) initClientWithRetry() {
	ErrMlServiceNotRunning := errors.New("ml service is not healthy")

	client := retryablehttp.NewClient()
	client.RetryMax = 5
	client.RetryWaitMin = 10 * time.Millisecond
	client.RetryWaitMax = 100 * time.Millisecond
	client.Logger = m.logger

	maxRetries := 3
	initialBackoff := time.Second

	for i := 0; i < maxRetries; i++ {
		resp, err := client.Get(m.config.MlServiceUri() + "/health")
		if err == nil && resp.StatusCode == http.StatusOK {
			// Successfully connected, set the client and signal readiness
			m.client = client
			close(m.clientReady)
			m.logger.Info("connected to ml service successfully")
			return
		}

		if err != nil {
			m.logger.Error(fmt.Sprintf("failed to connect to ml service: %v", err.Error()))
		} else {
			m.logger.Error(fmt.Sprintf("ml service returned non-OK status: %v", resp.StatusCode))
		}

		backoff := initialBackoff * time.Duration(1<<i)
		m.logger.Warn(fmt.Sprintf("retrying connection to ml service after %v", backoff))

		time.Sleep(backoff)
	}

	m.logger.Error("failed to connect to ml service after maximum retries")
	m.clientErr <- ErrMlServiceNotRunning
	close(m.clientReady)
}

// WaitForClientReady allows other parts of the app to wait for the ML service client to be ready
func (m *MlService) WaitForClientReady() {
	<-m.clientReady
}
