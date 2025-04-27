.PHONY: run, migrate, docker, docker-down, build, mockery, test, coverage, clean

COVERAGE_FILE = coverage.out
COVERAGE_HTML = coverage.html

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

clean-cache:
	go clean -testcache

test: clean clean-cache
	go test ./... -coverprofile=$(COVERAGE_FILE)

coverage: test
	go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)

clean:
	rm -f $(COVERAGE_FILE) $(COVERAGE_HTML)
