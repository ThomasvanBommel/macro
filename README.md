# Macro

Macro is a full-stack nutrition tracker built with React and Go.

Live app: https://macro.cekeh.com

## Why This Project Exists

This is a portfolio project focused on practical product delivery, not a toy demo.

It demonstrates:
- End-to-end full-stack ownership
- Session-based authentication and protected UI flows
- API + relational data modeling for a real user workflow
- CI/CD automation and containerized local development

## Current Feature Set

- Register, login, logout, and session validation
- Create and search foods
- Create, edit, and delete diary entries
- View diary entries by date and meal
- Daily macro totals and macro percentage split

## Architecture

```text
React (Vite) -> /api -> Go (Gin) -> SQLite
```

Key endpoints:
- POST /api/register
- POST /api/login
- POST /api/logout
- POST /api/validate-session
- POST /api/food
- POST /api/foods
- POST /api/food/search
- POST /api/entry
- POST /api/entry/edit
- POST /api/entry/delete
- POST /api/diary
- GET /api/health

## Tech Stack

- Frontend: React 19, React Router 7, Vite 8
- Backend: Go 1.26, Gin
- Data: SQLite, Goose migrations
- Dev workflow: Docker, Docker Compose, Make
- Auth/session: cookie-backed server sessions

## Local Development

Prerequisites: Docker, Docker Compose, Make

```bash
# inspect available targets
make

# run checks/build/tests in containers
make test

# run backend + frontend for development
make dev

# stop and clean containers
make down
```

Default local URLs:
- Frontend: http://localhost:5173
- Backend API: http://localhost:8080

Notes:
- SQLite data is persisted in backend/data
- docker-compose.yml sets a dev-only SESSION_SECRET for local use

## CI/CD

GitHub Actions runs on each push:
- Frontend dependency audit and build
- Backend tests and binary build
- Artifact packaging

On pushes to main:
- Build and push Docker image
- Deploy a new revision to Cloud Run

## Roadmap

High-priority next steps:
- Frontend test baseline (API helpers + core UI flows)
- Database layer refactor + targeted DB tests
- Profile UX polish and food editing

## Reviewer Notes

This project is intentionally iterative. The goal is to show the ability to ship a complete vertical slice, then improve reliability, test coverage, and maintainability over time.
