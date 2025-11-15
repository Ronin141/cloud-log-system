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
