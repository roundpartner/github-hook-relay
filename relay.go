package main

import (
	"bytes"
	"errors"
	"log"
	"net/http"
	"os"
)

func RelayRequestHook(hook *HookMessage) error {
	return RelayRequest(bytes.NewBufferString(hook.Message), "application/json", hook.MessageAttributes["X-Hub-Signature"].Value, hook.MessageAttributes["X-GitHub-Event"].Value)
}

func RelayRequest(payload *bytes.Buffer, contentType string, signature string, event string) error {
	endpoint := os.Getenv("JENKINS_ENDPOINT")
	if endpoint == "" {
		return errors.New("jenkins endpoint not set")
	}
	httpRequest, err := http.NewRequest("POST", endpoint, payload)
	if err != nil {
		log.Printf("[ERROR] %s", err.Error())
		return err
	}
	httpRequest.Header.Add("Content-Type", contentType)
	httpRequest.Header.Add("User-Agent", "GitHub Relay")
	httpRequest.Header.Add("X-Hub-Signature", signature)
	httpRequest.Header.Add("X-GitHub-Event", event)
	httpClient := http.Client{}
	httpResponse, err := httpClient.Do(httpRequest)
	if err != nil {
		log.Printf("[ERROR] %s", err.Error())
		return err
	}
	log.Printf("[INFO] Http Status: %d", httpResponse.StatusCode)

	return nil
}
