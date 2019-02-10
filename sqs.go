package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"log"
)

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
		log.Printf("[ERROR] SQS Error: %s\n", err.Error())
		return nil, err
	}

	if len(result.Messages) == 0 {
		return nil, nil
	}

	return result.Messages[0], nil
}
