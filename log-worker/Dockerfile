FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o log-worker

FROM alpine:3.19
WORKDIR /app

COPY --from=builder /app/log-worker .

EXPOSE 8080   

CMD ["./log-worker"]
