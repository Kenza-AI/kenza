package pubsub

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/kenza-ai/kenza/logutil"
	"github.com/streadway/amqp"
)

// RabbitMQ `Client` implementation.
type RabbitMQ struct {
	exchange   string
	channel    *amqp.Channel
	connection *amqp.Connection
}

// NewRabbitMQ - RabbitMQ client initializer.
func NewRabbitMQ(exchange, user, pass, host string, port int64, retryInterval time.Duration) (*RabbitMQ, error) {
	rabbitMQ := &RabbitMQ{}
	err := rabbitMQ.init(exchange, user, pass, host, port, retryInterval)
	return rabbitMQ, err
}

func (c *RabbitMQ) init(exchange, user, pass, host string, port int64, retryInterval time.Duration) error {
	connectionErrors := make(chan *amqp.Error)
	go func() {
		err := <-connectionErrors
		if err != nil {
			logutil.Info("rabbitmq reconnecting, reason '%s'", err)
			c.init(exchange, user, pass, host, port, retryInterval)
		}
	}()

	var err error
	var connection *amqp.Connection
	connectionString := fmt.Sprintf("amqp://%s:%s@%s:%d", user, pass, host, port)

	ticker := time.NewTicker(time.Second * retryInterval)
	defer ticker.Stop()

	for range ticker.C {
		if connection, err = amqp.Dial(connectionString); err == nil {
			break
		}
		logutil.Info("dialing exchange '%s' failed '%s', retrying in %d\"", exchange, err, retryInterval)
	}
	connection.NotifyClose(connectionErrors)

	channel, err := connection.Channel()
	if err != nil {
		return err
	}

	if err := declareExchange(exchange, channel); err != nil {
		return err
	}

	c.exchange = exchange
	c.connection = connection
	c.channel = channel

	return nil
}

// Subscribe declares the queue and binds the queue/routing key pair to an exchange.
//
// Use the messageCallback function to handle and ack the message.
func (c *RabbitMQ) Subscribe(queue, routingKey string, prefetchCount int, messageCallback func(body []byte, ack func(ok bool, requeue bool))) error {
	q, err := declareQueue(queue, c.channel)
	if err != nil {
		return err
	}

	if err := c.channel.QueueBind(q.Name, routingKey, c.exchange, false, nil); err != nil {
		return err
	}

	if err := c.channel.Qos(prefetchCount, 0, false); err != nil {
		return err
	}

	msgs, err := c.channel.Consume(
		queue, // queue
		"",    // consumer
		false, // auto ack
		false, // exclusive
		false, // no local
		false, // no wait
		nil,   // args
	)

	for delivery := range msgs {
		messageCallback(delivery.Body, func(ack bool, requeue bool) {
			if ack {
				delivery.Ack(requeue)
			} else {
				delivery.Nack(false, requeue)
			}
		})
	}

	return err
}

// Publish publishes a message to `Client`'s exchange.
func (c *RabbitMQ) Publish(msg interface{}, routingKey string) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	if err := c.channel.Publish(
		c.exchange, // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         []byte(body),
		}); err != nil {
		return err
	}

	return nil
}

// Close closes the Rabbit MQ connection and channel(s).
//
// This SHOULD be called, AMQP 0-9-1 asks for gracefully closing connections.
// https://www.rabbitmq.com/tutorials/amqp-concepts.html#amqp-connections
// Channels close with their connection.
//
// TODO(ilazakis): test to verify we close on cleanup.
func (c *RabbitMQ) Close() error {
	if c.connection.IsClosed() {
		return nil
	}
	return c.connection.Close()
}

func declareExchange(exchange string, ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		exchange, // name
		"topic",  // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
}

func declareQueue(queue string, ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		queue, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
}
