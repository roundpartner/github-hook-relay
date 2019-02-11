package main

import (
	"log"
	"os"
	"time"
)

func main() {
	log.SetOutput(os.Stdout)
	queue := os.Getenv("AWS_SQS_QUEUE")
	for {
		RelayHook(queue)
		time.Sleep(time.Second)
	}
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
	if message == nil {
		return nil
	}
	DeleteMessage(message.QueueUrl, message.ReceiptHandle)
	return nil
}
