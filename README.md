# macro
simple full stack macro tracking application

image now available @ https://hub.docker.com/repository/docker/cekeh/macro

it's live!
https://macro-292620794122.us-east1.run.app

## using
- [make](https://www.gnu.org/software/make/manual/make.html)
- [docker](https://www.docker.com/)
  - [node:20-alpine](https://hub.docker.com/_/node)
  - [golang:1.26-alpine](https://hub.docker.com/_/golang)
  - [alpine:latest](https://hub.docker.com/_/alpine)
  - [alpine/sqlite:latest](https://hub.docker.com/r/alpine/sqlite)
- backend
  - [go](https://go.dev/)
    - [goose](https://pkg.go.dev/github.com/pressly/goose/v3)
    - [sqlite](https://pkg.go.dev/modernc.org/sqlite)
    - [gin](https://gin-gonic.com)
      - [sessions](https://github.com/gin-contrib/sessions)
- frontend
  - [react](https://react.dev)
    - [vite](https://vite.dev)
    - [react-router](https://reactrouter.com)
    - [axios](https://axios-http.com/)

# Make Commands
```bash
# lists these commands and their descriptions
make help # or
make
```

```sh
# build and run the entire application
# exposed on :8080
make fullstack
```

```bash
# build and run both backend and frontend in parallel
# enables frontend's live edit
# exposed on :5173, backend proxied to :8080 via Vite config
make devmode
```

```bash
# build and run or restart only the backend
# exposed on :8080
make backend
```

```bash
# run shell in golang container for backend management
# copies backend project files you can execute go commands on
make golang
```

```bash
# run goose cli in container for database migrations
# copies database and migration files for modification
make goose
```

```bash
# create a migration using goose container
# also changes owner of migration to 1000:1000
make migration name="schema_name"
```

```bash
# manage database via sqlite container
make sqlite
```

```bash
# build and run or refresh frontend service
# enables live edit
# exposed on :5173
make frontend
```

```bash
# run npm cli for frontend package management
make npm
```

```bash
# stop and remove all containers, networks, and volumes
make down
```