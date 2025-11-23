migrate:
	goose -dir migrations postgres "$(DATABASE_URL)" up

run:
	CONFIG_PATH="$(CONFIG_PATH)" go run cmd/go-rest-api/main.go

build:
	mkdir -p bin
	go build -o bin/go-rest-api cmd/go-rest-api/main.go

docker-up:
	docker compose up --build -d

docker-down:
	docker compose down

test:
	go test ./tests/... -coverprofile=coverage.out

install:
	go mod download
