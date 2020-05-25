package pubsub

// A Publisher publishes messages/events of arbitrary type.
// to its exhange under the provided routing key.
type Publisher interface {
	Closer
	// Publish publishes an event to the `Publisher`'s exchange.
	// Events are marshaled into `JSON` before publishing.
	Publish(msg interface{}, routingKey string) error
}

// A Subscriber listens for messages/events.
type Subscriber interface {
	Closer
	// Subscribe declares the queue and binds the queue/routing key pair to an exchange.
	//
	// Use the messageCallback function to handle and ack incoming messages.
	Subscribe(queue, routingKey string, prefetchCount int, messageCallback func(body []byte, ack func(ok bool, requeue bool))) error
}

// Closer closes a pub/sub connection. Safe to call multiple times.
type Closer interface {
	// Close closes the AMQP connection and channel(s).
	//
	// This SHOULD be called, AMQP 0-9-1 asks for gracefully closing connections.
	// https://www.rabbitmq.com/tutorials/amqp-concepts.html#amqp-connections
	// Channels close with their connection.
	//
	// TODO(ilazakis): test we close on cleanup (worker, progress, api)
	Close() error
}
