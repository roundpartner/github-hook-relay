package main

import (
	"os"
	"testing"
)

func TestReadSqs(t *testing.T) {
	queue := os.Getenv("AWS_SQS_QUEUE")

	_, err := ReadSqs(queue)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
}
