package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"os"
	"testing"
)

func TestReadSqs(t *testing.T) {
	queue := os.Getenv("AWS_SQS_QUEUE")
	pushSqs(queue, "Hello world")

	_, err := ReadSqs(queue)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
}

func pushSqs(queue string, content string) {
	session := GetSession()
	sqsQueue := sqs.New(session)

	sqsQueue.SendMessage(&sqs.SendMessageInput{
		DelaySeconds: aws.Int64(0),
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"Type": {
				DataType:    aws.String("String"),
				StringValue: aws.String("message"),
			},
		},
		MessageBody: aws.String(content),
		QueueUrl:    &queue,
	})
}
