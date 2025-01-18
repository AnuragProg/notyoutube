package mq

import (
	"context"
)

type NoopQueue struct{}

func NewNoopQueue() *NoopQueue {
	return &NoopQueue{}
}

func (nq *NoopQueue) Publish(string,  string, []byte) error { return nil }
func (nq *NoopQueue) Subscribe(context.Context, []string, string, func([]byte) error) <-chan error  {
	errChan := make(chan error)
	defer close(errChan)
	return errChan
}
func (nq *NoopQueue) Close() error { return nil }
