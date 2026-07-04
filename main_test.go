package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthz(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/healthz", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestBadge(t *testing.T) {
	// mock HA server
	ha := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Error("missing or wrong auth header")
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		json.NewEncoder(w).Encode(haState{
			State: "84.3",
			Attributes: struct {
				UnitOfMeasurement string `json:"unit_of_measurement"`
			}{UnitOfMeasurement: "W"},
		})
	}))
	defer ha.Close()

	// build handler that talks to mock HA
	mux := http.NewServeMux()
	mux.HandleFunc("/", badgeHandler(ha.URL, "test-token", "sensor.power"))

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var b badge
	if err := json.NewDecoder(w.Body).Decode(&b); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if b.SchemaVersion != 1 {
		t.Errorf("expected schemaVersion 1, got %d", b.SchemaVersion)
	}
	if b.Label != "power" {
		t.Errorf("expected label power, got %s", b.Label)
	}
	if b.Message != "84.3 W" {
		t.Errorf("expected message '84.3 W', got '%s'", b.Message)
	}
	if b.Color != "brightgreen" {
		t.Errorf("expected color brightgreen, got %s", b.Color)
	}
}

func TestBadgeHADown(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", badgeHandler("http://localhost:1", "token", "sensor.power"))

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadGateway {
		t.Errorf("expected 502, got %d", w.Code)
	}
}

func TestBadgeDefaultUnit(t *testing.T) {
	ha := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// return state with no unit
		json.NewEncoder(w).Encode(haState{State: "100"})
	}))
	defer ha.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/", badgeHandler(ha.URL, "token", "sensor.power"))

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	var b badge
	json.NewDecoder(w.Body).Decode(&b)

	if b.Message != "100 W" {
		t.Errorf("expected '100 W', got '%s'", b.Message)
	}
}
