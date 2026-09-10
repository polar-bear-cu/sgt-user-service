DB_URL ?= postgres://postgres:postgres@localhost:5434/users?sslmode=disable

.PHONY: run test lint format tidy compose-up compose-down migrate-up migrate-down migrate-create

run:
	go run .

test:
	go test ./... -cover

format:
	golangci-lint fmt

lint:
	golangci-lint run

tidy:
	go mod tidy

compose-up:
	docker compose up --build -d --wait

compose-down:
	docker compose down

migrate-up:
	migrate -path migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path migrations -database "$(DB_URL)" down 1

migrate-create:
	@test -n "$(name)" || (echo "usage: make migrate-create name=<snake>"; exit 1)
	migrate create -ext sql -dir migrations "$(name)"
