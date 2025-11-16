package main

import (
	"bufio"
	"encoding/json"
	"log"
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

func main() {
	log.Println("Worker started...")

	for {
		filePath := getQueuePath()

		f, err := os.Open(filePath)
		if err != nil {
			time.Sleep(3 * time.Second)
			continue
		}

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()

			var entry LogEntry
			json.Unmarshal([]byte(line), &entry)

			// Processing logic (สัปดาห์หน้าเก็บ CosmosDB)
			log.Printf("[PROCESS] %s | %s | %s\n", entry.Timestamp, entry.Service, entry.Message)
		}

		f.Close()

		// clear file
		os.Remove(filePath)

		time.Sleep(3 * time.Second)
	}
}
