# Macro

Macro is a full-stack nutrition tracker built with React + Go.

Live app: https://macro.cekeh.com

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
- Docker + Compose + Make: reproducible dev environment and deploy path with minimal setup friction

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

## Run Locally

Prerequisites: Docker, Docker Compose, Make

```bash
# list targets
make

# run tests
make test

# dev mode: backend (:8080) + frontend (:5173)
make dev

# single-container full-stack mode (:8080)
make fullstack

# stop everything
make down
```

Notes:

- `SESSION_SECRET` defaults to a dev value in compose.
- SQLite data is stored in `backend/data` when using the provided compose setup.

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
- Deployment hardening and operational cleanup

## Reviewer Context

This project is intentionally presented as a working, evolving product rather than a finished demo. The goal is to show ability to ship a real vertical slice, then iterate pragmatically.
