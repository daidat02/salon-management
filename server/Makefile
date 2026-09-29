APP_NAME = salon_server 

-include .env
DATABASE_URL = $(DB_URL)
GIN_MODE = $(GIN_MODE)

.PHONY: dev run build tidy test lint migrate-up migrate-down migrate-create

dev:
	air
run:
	GIN_MODE=GIN_MODE go run ./cmd

build:
	CGO_ENABLED=0 go build -o bin/server ./cmd

tidy:
	go mod tidy

test:
	go test -v ./internal/usecase/...

lint:
	go vet ./...
migrate-up:
	migrate -path migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path migrations -database "$(DATABASE_URL)" down 1
	
# 	ví dụ: make migrate-create name=name_table
migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)