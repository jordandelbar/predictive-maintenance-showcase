package api

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"ml_facade/config"
	"ml_facade/internal/models/redis_models"
	"ml_facade/internal/service"
	"net/http"
	"sync"
	"time"
)

type Server struct {
	config     config.Config
	logger     *slog.Logger
	service    *service.MlService
	redisModel *redis_models.ThresholdModel
	version    string
	wg         *sync.WaitGroup
}

func NewApiServer(config config.Config, logger *slog.Logger, service *service.MlService, redisModel *redis_models.ThresholdModel, version string, wg *sync.WaitGroup) *Server {
	return &Server{
		config:     config,
		logger:     logger,
		service:    service,
		redisModel: redisModel,
		version:    version,
		wg:         wg,
	}
}

func (a *Server) Serve(ctx context.Context) error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", a.config.ApiServer.Port),
		Handler:      a.routes(),
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  30 * time.Second,
		ErrorLog:     slog.NewLogLogger(a.logger.Handler(), slog.LevelError),
	}

	go func() {
		a.logger.Info("starting api server", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.logger.Error("error api server: %v", err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return srv.Shutdown(shutdownCtx)
}
