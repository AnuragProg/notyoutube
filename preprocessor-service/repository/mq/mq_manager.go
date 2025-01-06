/*
MQ Manager has been created for the sole purpose of providing consistent codec and type safe way to handle mq messages
Will be using protocol buffers as our protocol

Will be using service name as the group id so as to make it part of the horizontally scaled preprocessor-services cluster
*/

package mq

import (
	"context"

	"github.com/anuragprog/notyoutube/preprocessor-service/configs"
	mqTypes "github.com/anuragprog/notyoutube/preprocessor-service/types/mq"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

type MessageQueueManager struct {
	mq MessageQueue
}

func NewMessageQueueManager(mq MessageQueue) *MessageQueueManager {
	return &MessageQueueManager{
		mq: mq,
	}
}

func (mqm *MessageQueueManager) SubscribeToRawVideoTopic(ctx context.Context, messageHandler func(*mqTypes.RawVideoMetadata) error) <-chan error {
	messageHandlerWrapper := func(message []byte) error {
		var decodedMessage mqTypes.RawVideoMetadata
		err := proto.Unmarshal(message, &decodedMessage)
		if err != nil {
			return err
		}
		return messageHandler(&decodedMessage)
	}
	return mqm.mq.Subscribe(ctx, []string{configs.MQ_TOPIC_RAW_VIDEO}, configs.SERVICE_NAME, messageHandlerWrapper)
}

func (mqm *MessageQueueManager) PublishToDAGTopic(message *mqTypes.DAG) error {
	encodedMessage, err := proto.Marshal(message)
	if err != nil {
		return err
	}
	return mqm.mq.Publish(configs.MQ_TOPIC_DAG, uuid.New().String(), encodedMessage)
}
