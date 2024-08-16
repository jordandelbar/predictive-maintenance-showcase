package consumer

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log/slog"
	"ml_facade/config"
	"ml_facade/internal/service"
	"sync"
	"time"
)

type RabbitMQConsumer struct {
	logger     *slog.Logger
	config     config.Config
	connection *amqp.Connection
	channel    *amqp.Channel
	connected  chan bool
	service    *service.MlService
	numWorkers int
	wg         *sync.WaitGroup
}

func NewRabbitMQConsumer(cfg config.Config, logger *slog.Logger, mlService *service.MlService, wg *sync.WaitGroup) *RabbitMQConsumer {
	return &RabbitMQConsumer{
		logger:     logger,
		config:     cfg,
		connected:  make(chan bool),
		service:    mlService,
		numWorkers: cfg.RabbitMQConsumer.NumWorkers,
		wg:         wg,
	}
}

func (c *RabbitMQConsumer) connect() error {
	conn, err := amqp.Dial(c.config.RabbitMQConsumer.URI)
	if err != nil {
		return err
	}
	c.connection = conn

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	c.channel = ch

	_, err = ch.QueueDeclare(
		c.config.RabbitMQConsumer.Queue,
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		return err
	}

	c.connected <- true
	return nil
}

func (c *RabbitMQConsumer) retryConnect(ctx context.Context) {
	initialBackoff := time.Second * 5
	maxBackoff := time.Minute

	for {
		select {
		case <-ctx.Done():
			c.logger.Info("rabbitmq retry connect received context cancellation")
			return
		default:
			for {
				backoff := initialBackoff
				for {
					if err := c.connect(); err == nil {
						return
					}
					c.logger.Warn(fmt.Sprintf("could not connect to RabbitMQ. Retrying in %v", backoff))
					time.Sleep(backoff)
					if backoff < maxBackoff {
						backoff *= 2
					}
				}
			}
		}
	}
}

func (c *RabbitMQConsumer) Close() {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.connection != nil {
		c.connection.Close()
	}
}

func (c *RabbitMQConsumer) Serve(ctx context.Context) error {

	go func() {
		for {
			select {
			case isConnected := <-c.connected:
				if isConnected {
					c.logger.Info("rabbitmq consumer is running")
					go c.consumeWrapper(ctx)
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	go c.retryConnect(ctx)
	return nil
}

func (c *RabbitMQConsumer) consumeWrapper(ctx context.Context) {
	if err := c.Consume(ctx); err != nil {
		c.logger.Error(fmt.Sprintf("rabbitmq consumer error: %v", err))
		c.Close()
		c.retryConnect(ctx)
	}
}
