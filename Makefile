create_env:
	cp .env.example .env

run:
	docker compose up --build

stop:
	docker compose down

unit_test:
	go test ./internal/... -coverprofile=coverage.out
	go tool cover -func=coverage.out

integration_test:
	go test ./integration_test

debug_run:
	go run cmd/server/main.go