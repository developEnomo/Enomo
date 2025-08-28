package main

import (
	"encoding/json"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})
	port := os.Getenv("PORT")
	if port == "" { port = "8080" }
	http.ListenAndServe(":"+port, nil)
}