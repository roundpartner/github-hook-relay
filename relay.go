package main

import (
	"bytes"
	"errors"
	"log"
	"net/http"
	"os"
)

func RelayRequestHook(hook *HookMessage) error {
	event := hook.MessageAttributes["X-Github-Event"].Value
	if event == "" {
		return errors.New("event not set")
	}
	signature := hook.MessageAttributes["X-Hub-Signature"].Value
	if signature == "" {
		return errors.New("signature not set")
	}
	return RelayRequest(bytes.NewBufferString(hook.Message), "application/json", signature, event)
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
	if httpResponse.StatusCode != 200 {
		buf := new(bytes.Buffer)
		buf.ReadFrom(httpResponse.Body)
		log.Printf("[INFO] Response: %s", buf.String())
	}

	return nil
}
