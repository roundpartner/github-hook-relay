package main

import (
	"bytes"
	"gopkg.in/h2non/gock.v1"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

func TestRelayHooks(t *testing.T) {
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

	defer gock.Off()
	gock.New("http://localhost:8080").
		Post("/github-webhook").
		Reply(http.StatusOK)

	gock.New("https://sqs.eu-west-2.amazonaws.com").
		EnableNetworking()

	os.Setenv("JENKINS_ENDPOINT", "http://localhost:8080/github-webhook")

	_, err = RelayHook(queue)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
}
