.PHONY: help backend frontend goose npm golang fmt clean

help:
	@echo "Available targets:"
	@echo "  backend   - Build and run the backend server (also builds frontend)"
	@echo "  golang    - Run a shell in the golang container for development"
	@echo "  goose     - Run goose CLI for database migrations"
	@echo "  fmt       - Format Go code using gofmt"

	@echo "  frontend  - Build and run the frontend development server"
	@echo "  npm       - Run npm CLI for frontend package management"

	@echo "  clean     - Stop and remove all containers, networks, and volumes"

# --

backend:
	docker compose up --build backend

golang:
	docker compose run -it --rm golang sh

goose:
	docker compose run -it --rm goose

fmt:
	docker compose run -it --rm golang gofmt -s -w .

# --

frontend:
	docker compose up --build frontend

npm:
	docker compose run -it --rm npm

# --

clean:
	docker compose down --remove-orphans