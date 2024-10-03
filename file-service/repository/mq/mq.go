package mq

import (
	"context"
	"io"
)

type MessageQueue interface {
	io.Closer
	SendMessage(topic, key string, message []byte) error
	ListenMessages(ctx context.Context, topics []string, groupID string, messageHandler func([]byte) error) <-chan error 
}
