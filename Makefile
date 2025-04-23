.PHONY: run, migrate, docker, docker-down, build

run:
	go run cmd/main.go

migrate:
	go run cmd/main.go -migrate

docker:
	docker compose up -d

docker-down:
	docker compose down

build:
	docker-compose up --build
