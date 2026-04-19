package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetRate_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"base":"USD","target":"EUR","rate":0.9}`))
	}))
	defer server.Close()

	service := NewExchangeService(server.URL)

	rate, err := service.GetRate("USD", "EUR")

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if rate != 0.9 {
		t.Errorf("expected 0.9, got %f", rate)
	}
}

func TestGetRate_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"invalid currency pair"}`))
	}))
	defer server.Close()

	service := NewExchangeService(server.URL)

	_, err := service.GetRate("USD", "XXX")

	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestGetRate_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`invalid json`))
	}))
	defer server.Close()

	service := NewExchangeService(server.URL)

	_, err := service.GetRate("USD", "EUR")

	if err == nil {
		t.Errorf("expected decode error, got nil")
	}
}

func TestGetRate_EmptyBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

	}))
	defer server.Close()

	service := NewExchangeService(server.URL)

	_, err := service.GetRate("USD", "EUR")

	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestGetRate_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`internal server error`))
	}))
	defer server.Close()

	service := NewExchangeService(server.URL)

	_, err := service.GetRate("USD", "EUR")

	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestGetRate_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(6 * time.Second)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"rate":1.0}`))
	}))
	defer server.Close()

	service := NewExchangeService(server.URL)

	_, err := service.GetRate("USD", "EUR")

	if err == nil {
		t.Errorf("expected timeout error, got nil")
	}
}
