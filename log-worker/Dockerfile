FROM golang:1.22 AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o log-worker

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/log-worker .
CMD ["./log-worker"]
