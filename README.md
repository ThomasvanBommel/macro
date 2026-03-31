# Macro

Macro is a full-stack nutrition tracking app with a React frontend and a Go backend.

This repository is intended as a portfolio project that demonstrates end-to-end product delivery: planning the data model, implementing the API, building a usable frontend, shipping a containerized runtime, and running it live.

Live app: https://macro.cekeh.com

## Why This Project

Macro was built to show practical full-stack ownership, not just isolated feature work. The emphasis is on delivering a working vertical slice across the full stack and leaving a codebase that can be iterated on.

What this project demonstrates:

- API design and implementation in Go
- Session-based authentication and protected frontend routes
- Relational data modeling for diary and nutrition workflows
- Frontend state management and route-level UX flows
- Container-first local development and deployment packaging
- Ongoing iteration with explicit tradeoffs and roadmap gaps

## Current Features

- User registration and login
- Cookie-based session auth
- Session countdown timer in the UI
- Protected profile route
- Food catalog browsing
- Custom food creation
- Daily food-entry logging by meal (breakfast, lunch, dinner, snack)
- Date-based diary view with daily totals

## What Is Intentionally Not Finished Yet

This app is functional and deployed, but it is still under active iteration.

- UX polish is incomplete in several areas
- Some interactions still use alert-based error handling
- Test coverage is not yet in place
- Admin and advanced editing flows are still planned

## Tech Stack

### Frontend

- React 19
- Vite 8
- React Router 7
- Axios
- Pico CSS

### Backend

- Go 1.26
- Gin
- gin-contrib/sessions (cookie store)
- bcrypt password hashing

### Data and Tooling

- SQLite
- Goose migrations (embedded and run at startup)
- Docker and Docker Compose
- Make

## Architecture

```text
React + Vite (frontend)
        |
        | HTTP / JSON
        v
Go + Gin API
        |
        | SQL
        v
SQLite
```

### Runtime Modes

Development mode (two containers):

- Backend on :8080
- Frontend (Vite dev server) on :5173
- Vite proxies /api to backend

Full-stack mode (single container):

- Frontend is built to frontend/build
- Build artifacts are embedded in the Go binary
- Go serves API and static frontend on :8080

## Project Layout

```text
macro/
├── backend/
│   ├── app/           # Go application (API, DB layer, models, entrypoint)
│   ├── data/          # SQLite file location (macro.db)
│   └── migrations/    # Goose SQL migrations
├── frontend/
│   ├── src/           # React app code
│   ├── public/        # Public assets
│   ├── package.json
│   └── vite.config.js
├── docker-compose.yml
├── Dockerfile
├── Makefile
└── README.md
```

## API Endpoints

All endpoints are under /api.

- POST /register
- POST /login
- POST /logout
- POST /entries
- POST /entry
- GET /foods
- POST /food

Notes:

- Login can also be used as a session check when called with an existing valid session cookie.
- /food and /entry require an authenticated session.

## Database Notes

Schema is managed by Goose and currently includes:

- user
- session
- food
- meal
- entry

Precision strategy:

- food macros/calories and serving_count are stored as integers scaled by 100
- entry servings are stored as integers scaled by 100

This avoids floating-point precision drift and values are scaled back in frontend API helpers.

## Local Development

Prerequisites:

- Docker
- Docker Compose
- Make

Session secret:

- Docker Compose reads it from session_secret.txt and mounts it as a Docker secret.
- Backend fallback: SESSION_SECRET env var.

### Useful Commands

```bash
# List available targets
make

# Full-stack mode (:8080)
make fullstack

# Dev mode (backend + frontend)
make dev

# Start only backend container
make backend

# Start only frontend container
make frontend

# Stop and remove containers
make down
```

Helper shells:

```bash
# Go shell with backend mounted
make golang

# Node shell with frontend mounted
make npm

# SQLite shell for backend/data/macro.db
make sqlite

# Goose tool shell
make goose
```

Image workflow:

```bash
# Build production image (distroless target), tag as cekeh/macro:latest
make build-image

# Push image to Docker Hub
make push-image
```

## Deployment and Persistence

- The full-stack container expects writable database storage at /app/data.
- For persistence, mount a host or named volume to /app/data.
- If no persistent mount is used, SQLite data will be ephemeral.

## Current Gaps / Planned Improvements

- Replace alert-based notifications with a proper notification system
- Profile page cleanup and styling pass
- Entry and food editing/removal flows
- More backend logging and cleanup
- Automated tests (frontend and backend)
- Deployment automation and hardening

## Notes for Reviewers

This is an actively evolving project and intentionally pragmatic in parts.

If you are reviewing this for hiring:

- The strongest signal is breadth plus execution, not final polish
- The project is already usable end-to-end and publicly deployable
- The remaining work is mostly iterative hardening and product refinement

In short, this repo is meant to show the ability to ship across frontend, backend, data, and infrastructure, then continue improving from a solid baseline.
