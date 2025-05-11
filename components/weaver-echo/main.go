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
	fmt.Printf("Echo received: %v\n", e.Data)
	// no further re-emit
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	fmt.Println("Starting server on port 8080")
	http.ListenAndServe(":8080", nil)
}
