FROM golang:1.24 AS builder

WORKDIR /app

# ปิด CGO และระบุ target OS/ARCH
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o log-worker

FROM alpine:3.19
WORKDIR /app

COPY --from=builder /app/log-worker .

EXPOSE 8080

CMD ["./log-api"]
