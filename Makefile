DB_URL=postgres://postgres:secret-key-for-db@localhost:5432/todo-api?sslmode=disable
BINARY_NAME=TODO-API

migrate-up:
	migrate -path migrations -database $(DB_URL) -verbose up

migrate-down:
	migrate -path migrations -database $(DB_URL) -verbose down

migrate-force:
	migrate -path migrations -database $(DB_URL) -verbose force

run:
	go run cmd/main.go

build:
	go build -o $(BINARY_NAME) cmd/main.go

.PHONY: migrate-up migrate-down migrate-force migrate-version migrate-create run build