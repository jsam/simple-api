package main

import "testing"

func TestNewAPI(t *testing.T) {
	api := NewAPI("0.0.0.0", 1234)
	if api == nil {
		t.Fatal("Creation of API service is broken.")
	}
}
