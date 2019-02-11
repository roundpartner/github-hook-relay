package main

import (
	"log"
	"os"
	"time"
)

func main() {
	log.SetOutput(os.Stdout)
	queue := os.Getenv("AWS_SQS_QUEUE")
	multiply := time.Duration(1)
	for {
		result, err := RelayHook(queue)
		if err != nil {
			log.Printf("[Error] %s", err)
			multiply = multiply * 2
		}
		if result {
			multiply = time.Duration(1)
		} else if multiply < 32 {
			multiply = multiply * 2
		}
		time.Sleep(time.Second * multiply)
	}
}

func RelayHook(queue string) (bool, error) {
	message, err := ReadHookFromSqs(queue)
	if err != nil {
		return false, err
	}
	if message == nil {
		return false, nil
	}
	err = RelayRequestHook(message)
	if err != nil {
		return true, err
	}
	DeleteMessage(message.QueueUrl, message.ReceiptHandle)
	return true, nil
}
