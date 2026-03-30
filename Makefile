.PHONY: help backend frontend goose npm golang fmt clean

help:
	@echo "Available targets:"
	@echo "  fullstack - Build and run the entire application"
	@echo "  devmode   - Build and run both backend and frontend in development mode"

	@echo "  backend   - Build and run the backend server in development mode"
	@echo "  golang    - Run a shell in the golang container for development"
	@echo "  goose     - Run goose CLI for database migrations"
	@echo "  migration - Create a new database migration with goose (usage: make migration " \
	                    "name=your_migration_name)"
	@echo "  sqlite    - Run a shell in the sqlite container for database management"

	@echo "  frontend  - Build and run the frontend development server"
	@echo "  npm       - Run npm CLI for frontend package management"
	@echo "  add       - Add all changes to git staging area and show status"

	@echo "  down      - Stop and remove all containers, networks, and volumes"

fullstack:
	docker compose up --build fullstack

devmode:
	docker compose up --build backend frontend

# --

backend:
	docker compose up --build backend

golang:
	docker compose run -it --rm golang sh

goose:
	docker compose run -it --rm goose sh

migration:
	docker compose run --rm goose sh -c "goose create $(name) sql && \
		chown -R 1000:1000 /app/migrations"

sqlite:
	docker compose run -it --rm sqlite

# --

frontend:
	docker compose up --build frontend

npm:
	docker compose run -it --rm npm

add:
	git add -A
	git status

docker-build:
	docker build -t macro-fullstack:latest -f ./Dockerfile.fullstack .
	docker images macro-fullstack

docker-tag:
	docker build -t macro-fullstack:latest -f ./Dockerfile.fullstack .
	docker tag macro-fullstack:latest cekeh/macro:latest
	docker images

docker-push:
	docker push cekeh/macro:latest

# --

down:
	docker compose down --remove-orphans