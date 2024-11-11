package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type AppHandler func(http.ResponseWriter, *http.Request) (interface{}, error)

func (fn AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data, err := fn(w, r)
	if err != nil {
		// Log the error
		fmt.Println("Error while processing the request", err)
		http.Error(w, "Failed to process request", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
