FROM golang:1.20-alpine AS builder
ENV GO111MODULE=on CGO_ENABLED=0
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY main.go .
COPY calc ./calc
COPY flags ./flags
RUN go vet && gofmt -s -w . && go build -o memory-calculator

FROM eclipse-temurin:17-jre-alpine
COPY --from=builder /app/memory-calculator /usr/local/bin/memory-calculator
COPY entrypoint.sh /tmp/entrypoint.sh
