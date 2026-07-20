set dotenv-load

database_url := "postgres://postgres:postgres@localhost:5432/finance_go?sslmode=disable"

run:
    go run ./cmd/api

docs:
    go run github.com/swaggo/swag/cmd/swag@v1.16.6 init -g cmd/api/main.go -o docs

test:
    go test ./...

vet:
    go vet ./...

migrate-up:
    go run ./cmd/migrate -direction up

migrate-down:
    go run ./cmd/migrate -direction down -steps 1
