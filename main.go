package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type badge struct {
	SchemaVersion int    `json:"schemaVersion"`
	Label         string `json:"label"`
	Message       string `json:"message"`
	Color         string `json:"color"`
}

type haState struct {
	State      string `json:"state"`
	Attributes struct {
		UnitOfMeasurement string `json:"unit_of_measurement"`
	} `json:"attributes"`
}

var client = &http.Client{Timeout: 5 * time.Second}

func badgeHandler(haURL, haToken, sensorID string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/states/%s", haURL, sensorID), nil)
		if err != nil {
			http.Error(w, "failed to create request", http.StatusInternalServerError)
			return
		}
		req.Header.Set("Authorization", "Bearer "+haToken)

		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, "failed to reach home assistant", http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "failed to read response", http.StatusInternalServerError)
			return
		}

		var state haState
		if err := json.Unmarshal(body, &state); err != nil {
			http.Error(w, "failed to parse response", http.StatusInternalServerError)
			return
		}

		unit := state.Attributes.UnitOfMeasurement
		if unit == "" {
			unit = "W"
		}

		b := badge{
			SchemaVersion: 1,
			Label:         "power",
			Message:       state.State + " " + unit,
			Color:         "green",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(b)
	}
}

func main() {
	haURL := os.Getenv("HA_URL")
	haToken := os.Getenv("HA_TOKEN")
	sensorID := os.Getenv("HA_SENSOR_ID")
	port := os.Getenv("PORT")

	if haURL == "" || haToken == "" || sensorID == "" {
		log.Fatal("HA_URL, HA_TOKEN, HA_SENSOR_ID are required")
	}

	if port == "" {
		port = "80"
	}

	http.HandleFunc("/", badgeHandler(haURL, haToken, sensorID))

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	log.Printf("power-badge listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
