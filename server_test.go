package main

import (
	"testing"
)

func TestNewRestServer(t *testing.T) {
	server := NewRestServer()
	if server == nil {
		t.Error("Rest server is not initialised")
	}
}
