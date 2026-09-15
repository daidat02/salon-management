APP_NAME = salon_server 

-include .env

.PHONY: dev run build tidy test lint

dev:
	air
run:
	GIN_MODE=release go run ./cmd/server

build:
	CGO_ENABLED=0 go build -o bin/server ./cmd/server

tidy:
	go mod tidy

test:
	go test ./...

lint:
	go vet ./...