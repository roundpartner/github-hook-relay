package main

import (
	"bytes"
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"log"
)

type HookMessage struct {
	Subject           string                    `json:"Subject"`
	Message           string                    `json:"Message"`
	MessageAttributes map[string]HookAttributes `json:"MessageAttributes"`
	QueueUrl          string
	ReceiptHandle     *string
}

type HookAttributes struct {
	Type  string `json:"Type"`
	Value string `json:"Value"`
}

func ReadHookFromSqs(queue string) (*HookMessage, error) {
	message, err := ReadSqs(queue)
	if err != nil {
		return nil, err
	}
	if message == nil {
		return nil, nil
	}

	messageObject := &HookMessage{
		QueueUrl:      queue,
		ReceiptHandle: message.ReceiptHandle,
	}
	err = json.Unmarshal(bytes.NewBufferString(*message.Body).Bytes(), messageObject)
	if err != nil {
		return nil, err
	}
	return messageObject, nil
}

func ReadSqs(queue string) (*sqs.Message, error) {
	session := GetSession()
	sqsQueue := sqs.New(session)

	result, err := sqsQueue.ReceiveMessage(&sqs.ReceiveMessageInput{
		AttributeNames: []*string{
			aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
		},
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		QueueUrl:            &queue,
		MaxNumberOfMessages: aws.Int64(1),
		VisibilityTimeout:   aws.Int64(360),
		WaitTimeSeconds:     aws.Int64(20),
	})

	if err != nil {
		log.Printf("[ERROR] SQS Error: %s", err.Error())
		return nil, err
	}

	if len(result.Messages) == 0 {
		return nil, nil
	}

	return result.Messages[0], nil
}

func DeleteMessage(queue string, receipt *string) {
	session := GetSession()
	sqsQueue := sqs.New(session)
	sqsQueue.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      &queue,
		ReceiptHandle: receipt,
	})
}
