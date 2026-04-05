.PHONY: help test dev down golang npm sqlite goose

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  Compose:"
	@echo "    test           Run tests in a temporary container"
	@echo "    dev            Build and run backend and frontend in development mode"
	@echo "    down           Stop and remove all containers"
	@echo ""
	@echo "  Helpers:"
	@echo "    golang         Start a Golang shell with the backend code mounted"
	@echo "    npm            Start a Node.js shell with the frontend code mounted"
	@echo "    sqlite         Start a SQLite shell with the database file mounted"
	@echo "    goose          Start a shell with Goose installed for database migrations"
	@echo ""


## Docker Compose

test:
	docker compose run --rm backend sh -c "cd backend && go test -v ./..."

dev: test
	docker compose up --build backend frontend

down:
	docker compose down --remove-orphans


## Helpers

golang:
	docker run -it --rm -v .:/app -w /app golang:1.26-alpine sh

npm:
	docker compose run -it --rm frontend sh

sqlite:
	docker run -it --rm -v ./backend/data:/app -w /app alpine/sqlite:latest macro.db

goose:
	docker compose run -it --rm goose sh