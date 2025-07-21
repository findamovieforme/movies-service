package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
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

	log.Printf("[CallLocalModel] calling predictor.py with title=%q", title)

	cmd := exec.Command("python3", "helpers/predictor.py")
	cmd.Stdin = bytes.NewBuffer(jsonData)

	var out, errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut

	err = cmd.Run()
	if err != nil {
		exitCode := -1
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		}
		log.Printf("[CallLocalModel] predictor.py failed: exit_code=%d err=%v", exitCode, err)
		log.Printf("[CallLocalModel] predictor.py stdout: %s", out.String())
		log.Printf("[CallLocalModel] predictor.py stderr: %s", errOut.String())
		return nil, fmt.Errorf("predictor.py exit %v: %w (stderr: %s)", err, err, errOut.String())
	}

	log.Printf("[CallLocalModel] predictor.py stdout (raw): %s", out.String())

	var rec RecResponse
	if err := json.Unmarshal(out.Bytes(), &rec); err != nil {
		log.Printf("[CallLocalModel] failed to unmarshal response: %v; raw stdout: %s", err, out.String())
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	if rec.Error != "" {
		log.Printf("[CallLocalModel] predictor returned error field: %s", rec.Error)
	}

	log.Printf("[CallLocalModel] got %d recommendations", len(rec.Recommendations))
	return &rec, nil
}
