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

	// ใช้ CREATE เพื่อไม่ให้ error ตอนไฟล์ยังไม่เกิดขึ้น
	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Println("worker: cannot open queue:", err)
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

	// Clear file โดยไม่ลบไฟล์
	f.Truncate(0)
	f.Seek(0, 0)
}

func main() {
	log.Println("Worker started...")

	// Loop ตลอดแบบไม่ทำให้ probe ตีเป็น unhealthy
	for {
		processQueue()
		time.Sleep(2 * time.Second)
	}
}
