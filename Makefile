.PHONY: build run test docker-build docker-run docker-stop clean bdd help up down logs clean-volumes

# Variables
IMAGE_NAME := go-demo
CONTAINER_NAME := go-demo-dev
HOST_PORT := 8080
CONTAINER_PORT := 80

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Build Go binary (via Docker)
	docker run --rm -v $(PWD):/app -w /app golang:1.23-alpine go build -o tmp/main ./cmd/server

run: ## Run server locally (via Docker, no hot-reload)
	docker run --rm -v $(PWD):/app -w /app -p $(HOST_PORT):$(CONTAINER_PORT) golang:1.23-alpine go run ./cmd/server

test: ## Run Go tests (via Docker)
	docker run --rm -v $(PWD):/app -w /app golang:1.23-alpine go test -v ./...

docker-build: ## Build Docker image
	docker build -t $(IMAGE_NAME) .

docker-run: docker-build ## Run in Docker with hot-reload
	docker run --rm -it \
		--name $(CONTAINER_NAME) \
		-v $(PWD):/app \
		-p $(HOST_PORT):$(CONTAINER_PORT) \
		$(IMAGE_NAME)

docker-stop: ## Stop running container
	docker stop $(CONTAINER_NAME) || true

clean: ## Remove build artifacts and Docker image
	rm -rf tmp/
	docker rmi $(IMAGE_NAME) || true

bdd: ## Run BDD tests (via Docker)
	docker run --rm -v $(PWD):/app -w /app golang:1.23-alpine go test -v ./features/...

bdd-build: ## Install BDD dependencies
	docker run --rm -v $(PWD):/app -w /app golang:1.23-alpine go get github.com/cucumber/godog/cmd/godog@latest

# Docker Compose targets
up: ## Start all services with docker-compose
	docker-compose up --build

up-d: ## Start all services in background
	docker-compose up --build -d

down: ## Stop all services
	docker-compose down

logs: ## Follow logs from all services
	docker-compose logs -f

clean-volumes: ## Stop services and remove data volumes
	docker-compose down -v
