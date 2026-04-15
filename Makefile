.PHONY: help test dev down sh sqlite goose

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "    test           Run tests in a temporary container"
	@echo "    dev            Build and run backend and frontend in development mode"
	@echo "    down           Stop and remove all containers"
	@echo ""
	@echo "    sh             Start a shell in the backend container"
	@echo "    sqlite         Start a SQLite shell with the database file mounted"
	@echo "    goose          Start a shell with Goose installed for database migrations"
	@echo ""


## Docker Compose

test-backend:
	@echo "\nBackend testing..."
	docker compose run --rm backend sh -c " \
		npm run build --prefix frontend && \
		cd backend && \
		go test -tags='no_postgres no_clickhouse no_mssql no_mysql prod' -v ./..."

test-frontend:
	@echo "\nFrontend testing..."
	docker compose run --rm frontend sh -c "CI=true npm test"

test:
	@echo "\nRebuild base, audit and build npm..."
	docker compose run --rm --build frontend sh -c " \
		npm i && \
		npm audit --audit-level=high && \
		npm run build"
	
	$(MAKE) -j2 test-backend test-frontend

dev:
	docker compose up --build backend frontend

down:
	docker compose down --remove-orphans


## Helpers

sh:
	docker compose run -it --rm backend sh

sqlite:
	docker run -it --rm -v ./backend/data:/app -w /app alpine/sqlite:latest macro.db

goose:
	docker compose run -it --rm goose sh