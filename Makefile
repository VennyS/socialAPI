.PHONY: run, migrate, docker, docker-down, build

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
