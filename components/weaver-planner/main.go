package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Event struct {
	Data map[string]interface{} `json:"data"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	var e Event
	json.NewDecoder(r.Body).Decode(&e)
	// log and re-emit “planned” event
	fmt.Printf("Planner received: %v\n", e.Data)
	resp := Event{Data: map[string]interface{}{"plannedFrom": e.Data}}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	fmt.Println("Starting server on port 8080")
	http.ListenAndServe(":8080", nil)
}
