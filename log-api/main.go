package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

type LogEntry struct {
	Service   string `json:"service"`
	Level     string `json:"level"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

func getQueuePath() string {
	// default container path (when running in docker)
	path := "/data/log_queue.jsonl"

	// for local dev on Windows (air), move file to project root: ../data/
	if os.PathSeparator == '\\' {
		// go up one directory (from log-api to project root)
		os.MkdirAll("../data", 0755)
		path = "../data/log_queue.jsonl"
	}

	return path
}

func logHandler(w http.ResponseWriter, r *http.Request) {
	var entry LogEntry

	if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if entry.Timestamp == "" {
		entry.Timestamp = time.Now().UTC().Format(time.RFC3339)
	}

	data, _ := json.Marshal(entry)

	path := getQueuePath()

	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Println("error writing queue:", err)
		http.Error(w, "server error", 500)
		return
	}
	defer f.Close()

	f.WriteString(string(data) + "\n")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"queued"}`))
}

func main() {
	http.HandleFunc("/logs", logHandler)
	log.Println("Log API running on :8080")
	http.ListenAndServe(":8080", nil)
}
