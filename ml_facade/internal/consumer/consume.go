package consumer

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
)

func (c *RabbitMQConsumer) Consume(ctx context.Context) error {
	messages, err := c.channel.Consume(
		c.config.RabbitMQConsumer.Queue,
		"",
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		return err
	}

	closeCh := make(chan *amqp.Error)
	c.connection.NotifyClose(closeCh)

	semaphore := make(chan struct{}, c.numWorkers)

	for {
		select {
		case msg := <-messages:
			if msg.Body != nil {
				semaphore <- struct{}{}

				go c.processBatch(msg, semaphore)
			}
		case <-ctx.Done():
			c.logger.Warn("context canceled, stopping RabbitMQConsumer")
			return nil
		case err := <-closeCh:
			return err
		}
	}
}

func (c *RabbitMQConsumer) processBatch(msg amqp.Delivery, semaphore chan struct{}) {
	defer func() {
		<-semaphore
	}()

	c.wg.Add(1)
	_, _, err := c.service.HandleMlServiceRequest(msg, "rabbitmq")
	if err != nil {
		c.logger.Error(fmt.Sprintf("error handling ml service request: %v", err))
	}
}
