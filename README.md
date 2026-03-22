# Macro

Macro is a full-stack nutrition tracking application built to demonstrate end-to-end product development: frontend implementation, backend API design, authentication, database modeling, containerized workflows, and live deployment.

Live app: https://macro.cekeh.com

## Why This Project Exists

I built Macro as a portfolio project to show that I can take an idea from concept to a working deployed application.

I built the first version in about a week of focused work, so the goal was not to produce a polished, production-ready product. The goal was to demonstrate full-stack ownership by handling the entire workflow:

- designing the schema
- building the API
- implementing authentication and sessions
- creating the React UI
- wiring frontend and backend together
- containerizing the project
- deploying a live version

It is still unfinished, and I plan to keep improving it over time. That context matters: this repository is best read as a strong working foundation and a demonstration of engineering range, not as a finished commercial product.

## What It Does

Macro lets a user:

- register and log in
- maintain a session with cookie-based auth
- view a protected profile page
- browse existing foods
- add custom foods with nutrition values
- create food entries for a specific day
- organize entries by meal
- view daily calorie and macro totals

## What This Demonstrates

The value of this project is that it required decisions across the full application surface rather than in just one layer.

### Frontend

- React application built with Vite
- client-side routing with protected and unprotected routes
- form handling for login, registration, food creation, and entry creation
- API integration against a Go backend
- session-aware UI flow

### Backend

- HTTP API built with Go and Gin
- cookie-backed session handling
- password hashing with bcrypt
- embedded frontend assets for single-app deployment
- database access layer for users, sessions, foods, and entries

### Data and Infrastructure

- SQLite persistence
- schema migrations with Goose
- Docker and Docker Compose workflow
- separate development and full-stack runtime modes
- live deployment so the app can be reviewed without local setup

## Stack

### Frontend

- React 19
- Vite
- React Router
- Axios
- Pico CSS

### Backend

- Go 1.26
- Gin
- gin-contrib/sessions
- bcrypt

### Data / Tooling

- SQLite
- Goose migrations
- Docker
- Docker Compose
- Make

## Architecture Overview

```text
React + Vite UI
    |
    | HTTP / JSON
    v
Go + Gin API
    |
    | SQL
    v
SQLite database

Dev mode:
- frontend runs on :5173
- backend runs on :8080
- Vite proxies /api to the backend

Full-stack mode:
- frontend build is embedded into the Go binary/runtime image
- backend serves both the API and the compiled frontend on :8080
```

## Project Structure

```text
macro/
├── backend/
│   ├── app/            # Go API, DB layer, models, app entrypoint
│   ├── data/           # SQLite database files
│   └── migrations/     # Goose SQL migrations
├── frontend/
│   ├── src/            # React app
│   ├── public/         # Static assets
│   └── package.json    # Frontend dependencies and scripts
├── docker-compose.yml  # Local orchestration
├── Makefile            # Common development commands
└── READMEV2.md
```

## Notable Implementation Details

- Sessions are stored server-side and identified via a cookie.
- Passwords are hashed with bcrypt before storage.
- Database migrations run automatically when the backend starts.
- Session cleanup runs on a ticker to remove expired sessions.
- In development, the frontend runs through Vite; in full-stack mode, the Go app serves the compiled frontend.
- The project is separated into frontend and backend concerns while still supporting a single-container deployment path.

## Running the Project

### Prerequisites

- Docker
- Docker Compose
- Make

### Main Commands

```bash
# Show available commands
make
```

```bash
# Run the full stack on :8080
make fullstack
```

```bash
# Run frontend and backend separately for development
# frontend: :5173
# backend:  :8080
make devmode
```

```bash
# Backend only
make backend
```

```bash
# Frontend only
make frontend
```

```bash
# Open a Go container for backend work
make golang
```

```bash
# Open a Goose container for migration work
make goose
```

```bash
# Create a new migration
make migration name="your_migration_name"
```

```bash
# Open SQLite against the local database volume
make sqlite
```

```bash
# Open an npm container for frontend package work
make npm
```

```bash
# Stop and remove containers, networks, and orphans
make down
```

## Development Notes

- In development mode, the frontend uses Vite and proxies API requests to the backend container.
- In full-stack mode, the backend serves both the compiled frontend build and the API.
- The live deployment should be treated as a working demo, not a production-hardened environment.
- Because the project is still evolving, data models and UI behavior may change over time.

## Current Limitations

This project is intentionally honest about where it stands today.

- The homepage is still minimal.
- The UI needs another styling pass.
- Notifications still use simple alerts in some places.
- Automated test coverage has not been added yet.
- Logging and production hardening are still limited.
- The live demo should not be treated as secure for sensitive personal data.

Those gaps are real, but they also reflect the prioritization tradeoff: I focused first on building a complete working vertical slice.

## What I Would Improve Next

- Add backend and frontend test coverage.
- Improve validation and error handling.
- Replace alert-based UX with a proper notification system.
- Refine the profile workflow and homepage.
- Expand logging and observability.
- Harden deployment and security practices.
- Automate deployment.

## Why I’d Include This In A Job Application

This project is still a work in progress. I reached a good stopping point for the application, but not a final one. What I wanted this repository to show is the kind of work I can do independently within a limited timeframe:

- Translate an idea into a real application.
- Make reasonable technology choices.
- Move between UI, API, data, and deployment work.
- Deliver something live and reviewable.
- Leave a codebase in a state that can continue evolving.

I ran out of time to push it further before submitting applications, but I plan to keep improving it. If you are reviewing this repository for a role, the main signal I want it to send is straightforward: I can build across the stack, deliver a working product, and continue iterating from a solid foundation.
