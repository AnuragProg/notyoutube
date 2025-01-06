/*
MQ Manager has been created for the sole purpose of providing consistent codec and type safe way to handle mq messages
Will be using protocol buffers as our protocol

Will be using service name as the group id so as to make it part of the horizontally scaled file-services cluster
*/

package mq

import (
	"github.com/anuragprog/notyoutube/file-service/configs"
	mqTypes "github.com/anuragprog/notyoutube/file-service/types/mq"
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

func (mqm *MessageQueueManager) PublishToRawVideoTopic(message *mqTypes.RawVideoMetadata) error {
	encodedMessage, err := proto.Marshal(message)
	if err != nil {
		return err
	}
	return mqm.mq.Publish(configs.MQ_TOPIC_RAW_VIDEO, uuid.New().String(), encodedMessage)
}
