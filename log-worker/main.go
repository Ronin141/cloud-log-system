package main

import (
	"bufio"
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
	path := "/data/log_queue.jsonl"
	dir := "/data"

	if os.PathSeparator == '\\' {
		dir = "../data"
		path = "../data/log_queue.jsonl"
	}

	os.MkdirAll(dir, 0755)
	return path
}

func processQueue() {
	filePath := getQueuePath()

	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Println("worker cannot open file:", err)
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var entries []LogEntry

	for scanner.Scan() {
		line := scanner.Text()

		var entry LogEntry
		if err := json.Unmarshal([]byte(line), &entry); err == nil {
			entries = append(entries, entry)
		}
	}

	if len(entries) > 0 {
		for _, e := range entries {
			log.Printf("[PROCESS] %s | %s | %s\n", e.Timestamp, e.Service, e.Message)
		}
	}

	f.Truncate(0)
	f.Seek(0, 0)
}

func main() {
	log.Println("Worker started...")

	// 🔥 Health Check Endpoint (เพื่อให้ Startup Probe ผ่าน)
	go func() {
		http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			w.Write([]byte("ok"))
		})
		http.ListenAndServe(":8080", nil)
	}()

	// 🔥 Worker loop
	for {
		processQueue()
		time.Sleep(2 * time.Second)
	}
}
