.PHONY: help full backend frontend dev down build-image push-image add commit golang npm sqlite \
        goose

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  Compose:"
	@echo "    full           Build and run both backend and frontend"
	@echo "    backend        Build and run only the backend in development mode"
	@echo "    frontend       Build and run only the frontend in development mode"
	@echo "    dev            Build and run backend and frontend in development mode"
	@echo "    down           Stop and remove all containers"
	@echo ""
	@echo "  Image:"
	@echo "    build-image    Build the Docker image for production"
	@echo "    push-image     Push the Docker image to Docker Hub"
	@echo ""
	@echo "  Git:"
	@echo "    add            Stage all changes for commit"
	@echo "    commit         Commit staged changes with a message"
	@echo ""
	@echo "  Helpers:"
	@echo "    golang         Start a Golang shell with the backend code mounted"
	@echo "    npm            Start a Node.js shell with the frontend code mounted"
	@echo "    sqlite         Start a SQLite shell with the database file mounted"
	@echo "    goose          Start a shell with Goose installed for database migrations"
	@echo ""


## Docker Compose

full:
	docker compose up --build fullstack

backend:
	docker compose up --build backend

frontend:
	docker compose up --build frontend

dev:
	docker compose up --build backend frontend

down:
	docker compose down --remove-orphans


## Docker image

build-image:
	docker build -t macro-distroless:latest --target distroless .
	docker tag macro-distroless:latest cekeh/macro:latest
	docker images macro-distroless

push-image: build-image
	docker push cekeh/macro:latest


## Git

add:
	git add -A
	git status

commit: add
	@read -p "Enter commit message: " msg && git commit -m "$$msg"


## Helpers

golang:
	docker run -it --rm -v ./backend:/app -w /app golang:1.26-alpine sh

npm:
	docker run -it --rm \
		-v ./frontend/src:/app/src \
		-v ./frontend/public:/app/public \
		-v ./frontend/index.html:/app/index.html \
		-v ./frontend/package.json:/app/package.json \
		-v ./frontend/vite.config.js:/app/vite.config.js \
		-w /app node:20-alpine sh

sqlite:
	docker run -it --rm -v ./backend/data:/app -w /app alpine/sqlite:latest macro.db

goose:
	docker compose run -it --rm goose sh