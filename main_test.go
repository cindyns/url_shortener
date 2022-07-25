package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestShorten(t *testing.T) {
	Urls = map[string]*UrlDatabase{}

	// Test empty url
	newReq, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/shorten?url=%s", "https://localhost:8090", ""), nil)
	rec := httptest.NewRecorder()

	shorten(rec, newReq)
	res := rec.Result()

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("got %v, expected %v", res.StatusCode, http.StatusOK)
	}

	// Test invalid url
	newReq, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("%s/shorten?url=%s", "https://localhost:8090", "http:/www.google.com"), nil)
	rec = httptest.NewRecorder()

	shorten(rec, newReq)
	res = rec.Result()

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("got %v, expected %v", res.StatusCode, http.StatusOK)
	}

	// Test valid url
	newReq, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("%s/shorten?url=%s", "https://localhost:8090", "https://www.google.com"), nil)
	rec = httptest.NewRecorder()

	shorten(rec, newReq)
	res = rec.Result()

	if res.StatusCode != http.StatusAccepted {
		t.Errorf("got %v, expected %v", res.StatusCode, http.StatusOK)
	}

	// Test duplicate url
	newReq, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("%s/shorten?url=%s", "https://localhost:8090", "https://www.google.com"), nil)
	rec = httptest.NewRecorder()

	shorten(rec, newReq)
	res = rec.Result()

	if res.StatusCode != http.StatusAccepted {
		t.Errorf("got %v, expected %v", res.StatusCode, http.StatusOK)
	}
}
