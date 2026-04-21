.PHONY: run build migrate-up migrate-down docker-up docker-down tidy

run:
	go run ./cmd/server

build:
	go build -o bin/server ./cmd/server

tidy:
	go mod tidy

docker-up:
	docker compose up -d

docker-down:
	docker compose down

migrate-up:
	go run ./scripts/migrate.go up

migrate-down:
	go run ./scripts/migrate.go down

seed:
	go run ./scripts/seed.go

lint:
	golangci-lint run ./...

test:
	go test ./... -v
