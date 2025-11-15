FROM golang:1.22 AS builder

WORKDIR /app

# copy go.mod
COPY go.mod .

# copy go.sum only if exists
RUN test -f go.sum || touch go.sum
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o log-worker

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/log-worker .
CMD ["./log-api"]
