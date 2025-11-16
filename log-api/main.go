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
	// default path for container
	path := "/data/log_queue.jsonl"
	dir := "/data"

	// Windows dev mode
	if os.PathSeparator == '\\' {
		dir = "../data"
		path = "../data/log_queue.jsonl"
	}

	// ALWAYS create directory
	os.MkdirAll(dir, 0755)

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
