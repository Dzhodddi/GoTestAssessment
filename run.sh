#!/bin/bash
set -e

mkdir -p bin
go build -o ./bin/main ./cmd/api

docker run -d --name my-postgres -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=adminadmin -e POSTGRES_DB=testGo -p 5432:5432 postgres:latest
docker run -d -p 6379:6379 --name my-redis redis:latest

echo "Setting up docker containers..."
sleep 5

go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
go install github.com/swaggo/swag/cmd/swag@latest

swag init -g api/main.go -d cmd,internal && swag fmt

migrate -path internal/db/migrations -database "postgresql://postgres:adminadmin@localhost:5432/testGo?sslmode=disable" up

./bin/main
