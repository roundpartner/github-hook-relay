package main

import "testing"

func TestGetSession(t *testing.T) {
	session := GetSession()
	if session == nil {
		t.FailNow()
	}
}
