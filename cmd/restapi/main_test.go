package main

import (
	"github.com/olivernadj/post-proc/internal/api"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAddAction(t *testing.T)  {
	b := strings.NewReader("{ \"action\": \"test\", \"state\": \"new\"}")
	req, err := http.NewRequest("POST", "http://localhost:8080/v1/action", b)
	if err != nil {
		t.Fatalf("Could not created request: %v", err)
	}
	req.Header.Add("Source-Type", "client")
	rec := httptest.NewRecorder()
	api.AddAction(rec, req)

	res := rec.Result()
	if res.StatusCode != http.StatusCreated {
		t.Errorf("Expected status code 201, got %d", res.StatusCode)
	}
}

func TestAddActionMissingHeader(t *testing.T)  {
	b := strings.NewReader("{ \"action\": \"test\", \"state\": \"new\"}")
	req, err := http.NewRequest("POST", "http://localhost:8080/v1/action", b)
	if err != nil {
		t.Fatalf("Could not created request: %v", err)
	}
	rec := httptest.NewRecorder()
	api.AddAction(rec, req)

	res := rec.Result()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code 400, got %d", res.StatusCode)
	}
}
