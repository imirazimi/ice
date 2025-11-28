up:
	docker compose up -d

down:
	docker compose down

migrate:
	go run ./cmd/main.go -migrate

run:
	go run ./cmd/main.go

run-dev:
	go run ./cmd/main.go -dev

swagger:
	swag init -g cmd/main.go -o docs

test:
	go test ./...

benchmark:
	go test -bench=. ./...
