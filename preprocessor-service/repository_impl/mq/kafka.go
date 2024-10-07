package mq

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
	"os"

	"github.com/IBM/sarama"
)


func init(){
	sarama.Logger = log.New(
		os.Stdout, "Custom Logger:",
		log.Ldate|log.Ltime|log.Lshortfile|log.Llongfile|log.LstdFlags,
	)
}


type KafkaConsumerGroupHandler struct {
	messageHandler func([]byte) error
}

func (kcgh KafkaConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error { return nil }

func (kcgh KafkaConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (kcgh KafkaConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claims sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claims.Messages():
			if !ok {
				return nil
			}
			fmt.Println("Received message on preprocessor side", string(message.Value))
			if err := kcgh.messageHandler(message.Value); err != nil {
				return err
			}

			// Auto commit flag whether set or not, this will handle the offsetting properly
			session.MarkMessage(message, "")
			session.Commit()

		// return when all the things are done
		case <-session.Context().Done():
			return nil
		}
	}
}

type KafkaQueue struct {
	brokers   []string
	producer  sarama.SyncProducer

	// Set of ConsumerGroup 
	consumers *sync.Map
}

func getDefaultKafkaConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Metadata.Retry.Max = 3
	config.Metadata.Retry.Backoff = time.Second * 2
	config.Metadata.RefreshFrequency = 5 * time.Second
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Group.Rebalance.Timeout = 10 * time.Second
	return config
}

func NewKafkaQueue(brokers []string) (*KafkaQueue, error) {
	producer, err := sarama.NewSyncProducer(brokers, getDefaultKafkaConfig())
	if err != nil {
		return nil, err
	}

	return &KafkaQueue{brokers: brokers, producer: producer, consumers: new(sync.Map)}, nil
}

func MustNewKafkaQueue(brokers []string) *KafkaQueue {
	producer, err := sarama.NewSyncProducer(brokers, getDefaultKafkaConfig())
	if err != nil {
		panic(err)
	}

	return &KafkaQueue{brokers: brokers, producer: producer, consumers: new(sync.Map)}
}

// sends the message over the given topic, key for partitioning
func (kc *KafkaQueue) Publish(topic, key string, message []byte) error {
	_, _, err := kc.producer.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(message),
	})
	return err
}

// creates a new consumer group and listens to the topic, returned error chan is for immediate & non immediate error reporting
func (kc *KafkaQueue) Subscribe(ctx context.Context, topics []string, groupID string, messageHandler func([]byte) error) <-chan error {
	errChan := make(chan error)

	// putting everything in a goroutine to ensure responsibility of closing err chan to the 
	// ListenMessages function only
	go func() {
		defer close(errChan)

		group, err := sarama.NewConsumerGroup(kc.brokers, groupID, getDefaultKafkaConfig())
		if err != nil {
			errChan<- err
			return
		}

		kc.consumers.Store(group, struct{}{})
		defer func(){
			group.Close() // close the consumer group
			kc.consumers.Delete(group) // delete the consumer group
		}()

		consumer := KafkaConsumerGroupHandler{
			messageHandler: messageHandler,
		}

		for {
			if err := group.Consume(ctx, topics, &consumer); err != nil {
				if errors.Is(err, sarama.ErrClosedConsumerGroup) {
					return
				}
				errChan<- err
				return
			}
			// cancelled using context
			if ctx.Err() != nil {
				return
			}
		}
	}()

	return errChan
}

func (kc *KafkaQueue) Close() error {
	err := kc.producer.Close()
	if err != nil {
		return err
	}

	kc.consumers.Range(func(key, value any) bool {
		consumer, ok := key.(sarama.ConsumerGroup)
		if !ok {
			panic("consumer not of type KafkaConsumerGroupHandler")
		}
		consumer.Close()
		return true
	})
	return nil
}
