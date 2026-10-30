.PHONY: up down seed server web lint

up: ## Start Meilisearch
	docker compose up -d --wait

down: ## Stop Meilisearch and delete its data
	docker compose down -v

seed: ## Load 50,000 dummy books
	cd backend && go run ./cmd/seed

server: ## Start the API server on :8080
	cd backend && go run ./cmd/server

web: ## Start the frontend on :5173
	cd frontend && npm install && npm run dev

lint: ## Lint and type-check everything
	cd backend && go vet ./... && go tool declscope shrink ./... && go tool declscope ./...
	cd frontend && npm run typecheck
