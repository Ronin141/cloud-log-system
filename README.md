# 🌥️ Cloud Logging Platform (Go Version)

A cloud-native, event-driven logging system built on **Azure Container
Apps**, designed for high‑volume log ingestion, background processing,
and end‑to‑end observability --- fully rewritten for **Golang API +
Golang Worker**.

------------------------------------------------------------------------

## 📌 Overview

This platform consists of two main Go-based services:

-   **API Service (Go Fiber / net/http)** --- receives logs from
    clients, validates data, and pushes them into Azure Storage Queue.
-   **Worker Service (Go)** --- continuously consumes queue messages,
    processes logs, transforms data, and sends them into Azure Log
    Analytics.

Architecture focuses on **scalability, performance, and
cost-efficiency** using Go's lightweight runtime.

------------------------------------------------------------------------

## 🏗️ Architecture

Client → Go API → Storage Queue → Go Worker → Log Analytics

    ┌────────────┐      POST /logs      ┌─────────────────────────┐
    │   Client   │ ───────────────────▶ │     Go API Container     │
    └────────────┘                      └───────────┬─────────────┘
                                                    │
                                                    ▼
                                          ┌───────────────────────┐
                                          │  Azure Storage Queue  │
                                          └───────────┬───────────┘
                                                      │
                                                      ▼
                                         ┌────────────────────────┐
                                         │    Go Worker (ACA)     │
                                         │  process + transform   │
                                         └───────────┬────────────┘
                                                     │
                                                     ▼
                                       ┌─────────────────────────────┐
                                       │ Log Analytics Workspace     │
                                       │ Query + Dashboard + Alerts  │
                                       └─────────────────────────────┘

------------------------------------------------------------------------

## 🚀 Features

### ✔ Built with Go (Fast, Lightweight, Low Memory)

Perfect for cloud workloads and high‑throughput logging.

### ✔ Event‑Driven

API pushes logs to queue instantly → zero backpressure.

### ✔ Auto‑scaling

-   API scales by HTTP load
-   Worker scales by queue depth

### ✔ Secure by Design

-   Managed Identity
-   No secrets stored in code
-   Optional Entra ID auth

### ✔ Deep Observability

-   Log Analytics (KQL)
-   Alerts
-   Dashboarding

------------------------------------------------------------------------

## 🛠 Technologies

  Component   Tech
  ----------- ----------------------------
  Compute     Azure Container Apps
  Queue       Azure Storage Queue
  Logging     Log Analytics Workspace
  API         Go (Fiber / net/http)
  Worker      Go
  IaC         Bicep / Terraform
  CI/CD       GitHub Actions / GitLab CI

------------------------------------------------------------------------

## 📁 Project Structure

    cloud-log-system/
    │
    ├── api/
    │   ├── main.go
    │   ├── Dockerfile
    │
    ├── worker/
    │   ├── worker.go
    │   ├── Dockerfile
    │
    └── README.md

------------------------------------------------------------------------

## ⚙️ API Specification

### **POST /logs**

#### Request Example

``` json
{
  "service": "inventory-api",
  "level": "error",
  "message": "database connection failed",
  "meta": { "retry": 3 }
}
```

### Response

``` json
{
  "status": "queued",
  "queueMessageId": "256d8d3c-45f9"
}
```

------------------------------------------------------------------------

## 🔧 Go Worker Logic

### Pseudo‑code (Go)

``` go
for {
    msg, err := queue.ReceiveMessage()
    if err != nil {
        log.Println("Queue error:", err)
        continue
    }

    var logData LogModel
    if err := json.Unmarshal([]byte(msg.Text), &logData); err != nil {
        log.Println("JSON error:", err)
        continue
    }

    if err := SendToLogAnalytics(logData); err != nil {
        log.Println("Push error:", err)
        continue
    }

    queue.DeleteMessage(msg)
}
```

### Key Behavior

-   Infinite worker loop (no HTTP server)
-   No port exposed
-   Pure background processing

------------------------------------------------------------------------

## 🐳 Local Development (Go)

### API

``` bash
docker build -t go-log-api ./api
docker run -p 3000:3000 go-log-api
```

### Worker

``` bash
docker build -t go-log-worker ./worker
docker run go-log-worker
```

------------------------------------------------------------------------

## ☁️ Deployment to Azure

### 1. Provision infra

``` bash
az deployment group create   --resource-group cloudlog-rg   --template-file main.bicep
```

### 2. Deploy Go API

``` bash
az containerapp up   --name go-log-api   --image registry/go-log-api:latest   --env-vars QUEUE_URL=...
```

### 3. Deploy Go Worker

``` bash
az containerapp up   --name go-log-worker   --image registry/go-log-worker:latest   --revision-suffix v1
```

------------------------------------------------------------------------

## 📊 Monitoring (KQL)

### Error Logs

``` kql
CloudLog_CL
| where Level == "error"
| order by TimeGenerated desc
```

### Log Volume per Service

``` kql
CloudLog_CL
| summarize count() by Service, bin(TimeGenerated, 1h)
```

------------------------------------------------------------------------

## 🪲 Troubleshooting

### ❌ Worker keeps restarting

Most common causes (Go version): - ❌ Worker accidentally starts an HTTP
server\
- ❌ Wrong health probe type → must be **"liveness: exec"** or disabled\
- ❌ Queue endpoint unreachable\
- ❌ Missing Managed Identity permissions

### Check logs

``` bash
az containerapp logs show   -n go-log-worker   -g cloudlog-rg   --type system --follow
```

------------------------------------------------------------------------

## 📌 Future Enhancements

-   Add DLQ
-   OpenTelemetry
-   Retry policies
-   Batch ingestion to Log Analytics\
-   Blob archival

------------------------------------------------------------------------

## 👤 Author
Tatsn Limsodsai
Cloud Engineer / DevOps Engineer
Golang-based Logging Architecture
Designed, built, and deployed by : Tatsn Limsodsai
