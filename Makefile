.PHONY: help test-audit test-backend test-frontend test dev down sh sqlite goose

# backend test tags
BTT ?= 

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "    help           Show this help message"
	@echo ""
	@echo "    test-audit     Run npm audit and build frontend"
	@echo "    test-backend   Run backend tests in dev mode"
	@echo "    test-frontend  Run frontend tests"
	@echo ""
	@echo "    test           Run all tests (audit, backend dev, backend prod, frontend)"
	@echo "    dev            Build and run backend and frontend in development mode"
	@echo "    down           Stop and remove all containers"
	@echo ""
	@echo "    sh             Start a shell in the backend container"
	@echo "    sqlite         Start a SQLite shell with the database file mounted"
	@echo "    goose          Start a shell with Goose installed for database migrations"
	@echo ""


## Docker Compose

test-audit:
	@echo "\nAudit and build npm..."
	docker compose run --rm frontend sh -c " \
		npm i && \
		npm audit --audit-level=high && \
		npm run build"

test-backend:
	@echo "\nBackend testing..."
	docker compose run -w /repo/backend --rm backend sh -c " \
		go test -tags='no_postgres no_clickhouse no_mssql no_mysql $(BTT)' -v ./..."

test-frontend:
	@echo "\nFrontend testing..."
	docker compose run --rm frontend sh -c "CI=true npm test"

test:
	$(MAKE) -j2 test-backend test-audit
	$(MAKE) BTT=prod -j2 test-backend test-frontend

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