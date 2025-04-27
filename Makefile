.PHONY: run, migrate, docker, docker-down, build, mockery, test

run:
	docker-compose up db redis -d
	go run cmd/main.go

migrate:
	go run cmd/main.go -migrate

docker:
	docker compose up -d

docker-down:
	docker compose down

build:
	docker-compose up --build

mocks:
	mockery --all --output=./internal/mocks

test:
	go test ./...