package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type RecResponse struct {
	Recommendations []string `json:"recommendations"`
	Error           string   `json:"error"`
}

func CallLocalModel(title string) (*RecResponse, error) {
	payload := map[string]string{"title": title}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Printf("[CallLocalModel] failed to marshal payload: %v", err)
		return nil, fmt.Errorf("marshal payload: %w", err)
	}

	log.Printf("[CallLocalModel] calling Flask recommender with title=%q", title)

	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("POST", "http://127.0.0.1:5000/predict", bytes.NewReader(jsonData))
	if err != nil {
		log.Printf("[CallLocalModel] failed to create HTTP request: %v", err)
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[CallLocalModel] HTTP call to recommender failed: %v", err)
		return nil, fmt.Errorf("call recommender: %w", err)
	}
	defer resp.Body.Close()

	var rec RecResponse
	if err := json.NewDecoder(resp.Body).Decode(&rec); err != nil {
		log.Printf("[CallLocalModel] failed to decode recommender response: %v", err)
		return nil, fmt.Errorf("decode recommender response: %w", err)
	}

	if resp.StatusCode >= 400 {
		log.Printf("[CallLocalModel] recommender returned HTTP %d and error=%q", resp.StatusCode, rec.Error)
		return nil, fmt.Errorf("recommender HTTP %d: %s", resp.StatusCode, rec.Error)
	}

	if rec.Error != "" {
		log.Printf("[CallLocalModel] predictor returned error field: %s", rec.Error)
	}

	log.Printf("[CallLocalModel] got %d recommendations", len(rec.Recommendations))
	return &rec, nil
}
