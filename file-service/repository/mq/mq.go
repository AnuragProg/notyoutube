package mq

import (
	"context"
	"io"
)

type MessageQueue interface {
	io.Closer
	SendMessage(topic, key string, message []byte) error
	ListenMessages(ctx context.Context, topic []string, groupID string, handler func([]byte)) error
}
