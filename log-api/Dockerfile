FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o log-api

FROM alpine:3.19
WORKDIR /app

COPY --from=builder /app/log-api .

EXPOSE 8080

CMD ["./log-api"]
