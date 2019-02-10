package main

import (
	"bytes"
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
	"io/ioutil"
	"os"
	"testing"
)

func TestReadSqs(t *testing.T) {
	queue := os.Getenv("AWS_SQS_QUEUE")
	topic := os.Getenv("AWS_SNS_TOPIC")
	content, err := ioutil.ReadFile("tests/valid-hook.req")
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	err = pushSns(bytes.NewBuffer(content).String(), topic)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}

	message, err := ReadSqs(queue)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}

	DeleteMessage(queue, message.ReceiptHandle)

	messageObject := &HookMessage{}
	err = json.Unmarshal(bytes.NewBufferString(*message.Body).Bytes(), messageObject)
	if err != nil {
		t.Errorf("message: %v", *message.Body)
		t.Fatalf("%s", err.Error())
	}

	if _, ok := messageObject.MessageAttributes["X-Github-Event"]; ok == false {
		t.Fatalf("Attribute not set")
	}

	if messageObject.MessageAttributes["X-Github-Event"].Value == "" {
		t.Fatalf("Event value is not set")
	}
}

func TestReadHookFromSqs(t *testing.T) {
	queue := os.Getenv("AWS_SQS_QUEUE")
	topic := os.Getenv("AWS_SNS_TOPIC")
	content, err := ioutil.ReadFile("tests/valid-hook.req")
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	err = pushSns(bytes.NewBuffer(content).String(), topic)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}

	messageObject, err := ReadHookFromSqs(queue)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}

	DeleteMessage(messageObject.QueueUrl, messageObject.ReceiptHandle)

	if _, ok := messageObject.MessageAttributes["X-Github-Event"]; ok == false {
		t.Fatalf("Attribute not set")
	}
	if messageObject.MessageAttributes["X-Github-Event"].Value == "" {
		t.Fatalf("Event value is not set")
	}
}

func pushSns(content string, topic string) error {
	params := &sns.PublishInput{
		Subject:  aws.String("GitHub Web Hook"),
		Message:  aws.String(content),
		TopicArn: aws.String(topic),
	}
	attr := map[string]*sns.MessageAttributeValue{
		"X-Github-Event": {
			DataType:    aws.String("String"),
			StringValue: aws.String("push"),
		},
		"X-Hub-Signature": {
			DataType:    aws.String("String"),
			StringValue: aws.String("sha1=0000000000000000000000000000000000000000"),
		},
	}
	params.SetMessageAttributes(attr)

	session := GetSession()
	snsService := sns.New(session)
	_, err := snsService.Publish(params)
	return err
}
