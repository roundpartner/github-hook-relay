package main

import (
	"log"
	"os"
)

func main() {
	log.SetOutput(os.Stdout)
	queue := os.Getenv("AWS_SQS_QUEUE")
	RelayHook(queue)
}

func RelayHook(queue string) error {
	message, err := ReadHookFromSqs(queue)
	if err != nil {
		return err
	}
	err = RelayRequestHook(message)
	if err != nil {
		return err
	}
	DeleteMessage(message.QueueUrl, message.ReceiptHandle)
	return nil
}
