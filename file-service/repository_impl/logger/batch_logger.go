package logger

import (
	"bytes"
	"io"
	"sync"
	"time"
)

type BatchLogger struct {
	writer io.WriteCloser

	buffer    *bytes.Buffer
	bufferMut *sync.Mutex

	flushThreshold int
	flushInterval  time.Duration

	quitChan <-chan struct{}
}

func NewBatchLogger(writer io.WriteCloser, flushThreshold int, flushInterval time.Duration) *BatchLogger {
	batchLogger := &BatchLogger{
		writer: writer,

		buffer:    new(bytes.Buffer),
		bufferMut: new(sync.Mutex),

		flushThreshold: flushThreshold,
		flushInterval:  flushInterval,

		quitChan: make(chan struct{}),
	}

	go batchLogger.periodicFlush()

	return batchLogger
}

// Periodically flush buffer contents
func (logger *BatchLogger) periodicFlush() {
	ticker := time.NewTicker(logger.flushInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			logger.Flush()
		case <-logger.quitChan:
			return
		}
	}
}

// Flush buffer content to writer
func (logger *BatchLogger) Flush() error {
	logger.bufferMut.Lock()
	defer logger.bufferMut.Unlock()

	if logger.buffer.Len() == 0 {
		return nil
	}

	if _, err := logger.writer.Write(logger.buffer.Bytes()); err != nil {
		return err
	}

	logger.buffer.Reset()

	return nil
}

// This method can be called concurrently
func (logger *BatchLogger) Write(p []byte) (n int, err error) {

	// write to the buffer
	logger.bufferMut.Lock()
	n, err = logger.buffer.Write(p)
	shouldFlush := logger.buffer.Len() >= logger.flushThreshold
	logger.bufferMut.Unlock()

	// check for flush on threshold
	if !shouldFlush {
		return
	}

	// flush
	logger.Flush()

	return
}

func (logger *BatchLogger) Close() error {
	return logger.writer.Close()
}
