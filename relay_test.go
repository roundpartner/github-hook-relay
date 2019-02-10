package main

import (
	"bytes"
	"gopkg.in/h2non/gock.v1"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

func TestRelayRequest(t *testing.T) {
	content, err := ioutil.ReadFile("tests/valid-hook.req")
	if err != nil {
		t.Skipf("%s", err.Error())
	}
	defer gock.Off()
	gock.New("http://localhost:8080").
		Post("/github-webhook").
		Reply(http.StatusOK)

	os.Setenv("JENKINS_ENDPOINT", "http://localhost:8080/github-webhook")
	buffer := bytes.NewBuffer(content)

	err = RelayRequest(buffer, "application/json", "sha1=0000000000000000000000000000000000000000", "push")
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
}
