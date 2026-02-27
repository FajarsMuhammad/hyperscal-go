.PHONY: help run build clean test install migrate

help: ## Tampilkan bantuan
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

install: ## Install dependencies
	go mod download
	go mod tidy

run: ## Jalankan aplikasi
	go run main.go

build: ## Build aplikasi
	go build -o bin/hyperscal-go main.go

clean: ## Bersihkan build artifacts
	rm -rf bin/
	go clean

test: ## Jalankan tests
	go test -v ./...

dev: ## Jalankan dengan hot reload (requires air)
	air

docker-postgres: ## Jalankan PostgreSQL dengan Docker
	docker run --name hyperscal-postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=hyperscal_db -p 5432:5432 -d postgres:15

docker-postgres-stop: ## Stop PostgreSQL container
	docker stop hyperscal-postgres
	docker rm hyperscal-postgres
