# TODO

## Current Sprint (Highest Hiring Signal)

### 1) Frontend Test Baseline (Top Priority)
- [ ] add Vitest + React Testing Library setup
- [ ] add tests for API helpers in frontend/src/api.js
- [ ] add login success/failure flow tests
- [ ] add diary create/edit/delete flow tests
- [ ] run frontend tests in CI

### 2) Reliability and Product Polish
- [ ] handle diary fetch errors in frontend/src/Components/Diary/Diary.jsx
- [ ] improve session refresh strategy (avoid aggressive session checks)
- [ ] add backend tests for unauthorized/invalid edit and delete entry cases

### 3) Database Maintainability
- [ ] split backend/db/db.go into focused files by domain
- [ ] add DB-layer tests for ownership checks and food search ordering

## Next Sprint

### Product Depth
- [ ] add ability to edit foods (backend + frontend + tests)
- [ ] profile page cleanup and UX pass

### CI/CD Quality Bar
- [ ] add lint steps to CI for frontend and backend
- [ ] harden CI dependency setup to avoid future breakage
- [ ] add post-deploy smoke checks (health + one protected endpoint)

## Documentation (Interview Leverage)

- [ ] add Architecture Decision Notes to README
- [ ] add one simple architecture/data-flow diagram to README
- [ ] add "how I would scale this" section (auth, DB, observability)

## Backlog (Lower Hiring Signal)

- [ ] github-like year-at-a-glance diary view
- [ ] create administration area