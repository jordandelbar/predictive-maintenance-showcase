package consumer

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
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

	const batchSize = 5
	const batchTimeout = 50 * time.Millisecond
	buffer := make([]amqp.Delivery, 0, batchSize)
	timer := time.NewTimer(batchTimeout)
	defer timer.Stop()

	for {
		select {
		case msg, ok := <-messages:
			if !ok {
				if len(buffer) > 0 {
					if err := c.dispatchBatch(buffer, semaphore); err != nil {
						return err
					}
				}
				return nil
			}

			if msg.Body != nil {
				buffer = append(buffer, msg)
				if len(buffer) == 1 {
					if !timer.Stop() {
						<-timer.C
					}
					timer.Reset(batchTimeout)
				}
				if len(buffer) >= batchSize {
					if err := c.dispatchBatch(buffer, semaphore); err != nil {
						return err
					}
					buffer = make([]amqp.Delivery, 0, batchSize)
					if !timer.Stop() {
						<-timer.C
					}
					timer.Reset(batchTimeout)
				}
			}

		case <-timer.C:
			if len(buffer) > 0 {
				if err := c.dispatchBatch(buffer, semaphore); err != nil {
					return err
				}
				buffer = make([]amqp.Delivery, 0, batchSize)
			}
			timer.Reset(batchTimeout)

		case <-ctx.Done():
			c.logger.Warn("context canceled, stopping RabbitMQConsumer")
			if len(buffer) > 0 {
				if err := c.dispatchBatch(buffer, semaphore); err != nil {
					return err
				}
			}
			return nil
		case err := <-closeCh:
			if len(buffer) > 0 {
				if dispatchErr := c.dispatchBatch(buffer, semaphore); dispatchErr != nil {
					c.logger.Error(fmt.Sprintf("error dispatching final batch: %v", dispatchErr))
				}
			}
			return err
		}
	}
}

// dispatchBatch handles the dispatching of a batch of messages.
// It acquires a semaphore slot, starts a goroutine for processing,
// and ensures the semaphore is released after processing.
func (c *RabbitMQConsumer) dispatchBatch(batch []amqp.Delivery, semaphore chan struct{}) error {
	// Acquire semaphore slot
	semaphore <- struct{}{}

	// Make a copy of the batch to avoid data races
	batchCopy := make([]amqp.Delivery, len(batch))
	copy(batchCopy, batch)

	go c.processBatch(batchCopy, semaphore)
	return nil
}

func (c *RabbitMQConsumer) processBatch(msgs []amqp.Delivery, semaphore chan struct{}) {
	defer func() {
		<-semaphore
	}()

	c.wg.Add(1)
	_, _, err := c.service.HandleMlServiceRequest(msgs, "rabbitmq")
	if err != nil {
		c.logger.Error(fmt.Sprintf("error handling ml service request: %v", err))
	}
}
