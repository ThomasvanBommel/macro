# Macro

Macro is a full-stack nutrition tracker built with React + Go.

Live app: https://macro.cekeh.com

Deployment: automatic (production deploys are triggered by repository updates; no manual deploy command is required in this repo).

Status: Active work in progress. Core functionality is usable end-to-end; polish and hardening are ongoing.

## What This Shows

- Full-stack ownership from schema design to deployable runtime
- Session-based authentication and protected application flows
- Practical API and data modeling for a diary-style product
- Containerized local development and production packaging

## Current Features

- Register, login, logout
- Cookie-backed sessions with session validation endpoint
- Browse foods and create custom foods
- Create diary entries by meal and date
- View date-specific entries and daily totals
- Protected profile route and session timer in UI

## Stack and Why

- React 19 + React Router 7: fast iteration, clean component model, straightforward route-level auth UX
- Vite 8: very fast local feedback loop and simple production build pipeline
- Go 1.26 + Gin: small, explicit backend with low overhead and predictable request handling
- gin-contrib/sessions (cookie store): simple server-controlled session model without token complexity
- SQLite + Goose migrations: lightweight persistence with schema versioning and easy portability
- Docker + Compose + Make: reproducible local development workflow with minimal setup friction

## Architecture

```text
React (Vite) -> /api -> Go (Gin) -> SQLite
```

All API routes are under `/api`.

Primary endpoints:

- POST /register
- POST /login
- POST /logout
- POST /validate-session
- POST /foods
- POST /food
- POST /entries
- POST /entry
- GET /health

## Local Workflow

Prerequisites: Docker, Docker Compose, Make

```bash
# list targets and helper containers
make

# run backend tests in a disposable container
make test

# start local development stack
# backend:  http://localhost:8080
# frontend: http://localhost:5173
make dev

# stop and remove containers
make down
```

Notes:

- `make dev` runs `make test` first, then starts `backend` and `frontend` via Compose.
- `SESSION_SECRET` defaults to a dev-only value in `docker-compose.yml`.
- SQLite data is persisted in `backend/data`.

Optional helper shells:

- `make golang` for a Go shell with the repo mounted.
- `make npm` for a Node shell in the frontend container.
- `make sqlite` for direct SQLite access to `backend/data/macro.db`.
- `make goose` for migration operations.

## Deployment

Production is auto-deployed. This repository no longer includes a manual deploy command or a dedicated production Compose flow.

What is in-repo today:

- `docker-compose.yml` is focused on local development services (`backend`, `frontend`, `goose`).
- `Dockerfile` includes trimmed build targets for local/dev workflows.
- Runtime behavior supports trusted proxy handling in managed environments (for example when `K_SERVICE` is set).

## Repo Layout

```text
backend/app      Go API and server
backend/data     SQLite database files
backend/migrations  Goose SQL migrations
frontend/src     React application
```

## Known Gaps (In Progress)

- More complete automated test coverage
- UX polish and notification handling improvements
- Editing/deleting existing food and entry records
- Additional deployment hardening and operational cleanup

## Reviewer Context

This project is intentionally presented as a working, evolving product rather than a finished demo. The goal is to show ability to ship a real vertical slice, then iterate pragmatically.
