.PHONY: help backend frontend goose npm golang fmt clean

help:
	@echo "Available targets:"
	@echo "  backend   - Build and run the backend service"
	@echo "  frontend  - Build and run the frontend service"
	@echo "  goose     - Run a shell in the goose container for database migrations"
	@echo "  npm       - Run a shell in a Node.js container with the frontend code mounted"
	@echo "  golang    - Run a shell in a Golang container with the backend code mounted"
	@echo "  fmt       - Format the Go code using gofmt"
	@echo "  clean     - Stop and remove all containers and networks created by docker-compose"

backend:
	docker compose up --build backend

frontend:
	docker compose up --build frontend

goose:
	docker compose run -it --rm goose sh

npm:
	docker run -it --rm \
		-v "$(PWD)/frontend":/app \
		-w /app node:20-alpine sh

golang:
	docker run -it --rm \
		-v "$(PWD)/backend/app":/app \
		-w /app golang:1.26-alpine sh

fmt:
	docker run --rm \
		-v "$(PWD)/backend/app":/app \
		-w /app golang:1.26-alpine \
		gofmt -s -w .

clean:
	docker compose down --remove-orphans