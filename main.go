package main

import (
	"log"
	"os"
	"time"
)

func main() {
	log.SetOutput(os.Stdout)
	queue, exists := os.LookupEnv("AWS_SQS_QUEUE")
	if !exists {
		log.Printf("AWS_SQS_QUEUE must be set")
		os.Exit(0)
	}
	multiply := time.Duration(1)
	for {
		result, err := RelayHook(queue)
		if err != nil {
			log.Printf("[ERROR] %s", err)
		}
		if result {
			multiply = time.Duration(1)
		} else if multiply < 32 {
			multiply = multiply * 2
			log.Printf("[INFO] Sleeping for %d seconds", multiply)
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
