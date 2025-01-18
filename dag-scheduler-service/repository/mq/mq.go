package mq

import (
	"context"
	"io"
)

// INFO: Not to be used directly for publishing and subscribing but rather use MessageQueueManager for doing the same
type MessageQueue interface {
	io.Closer
	Publish(topic, key string, message []byte) error
	Subscribe(ctx context.Context, topics []string, groupID string, messageHandler func([]byte) error) <-chan error 
}
