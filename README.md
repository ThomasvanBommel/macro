# macro
simple full stack macro tracking application

image now available @ https://hub.docker.com/repository/docker/cekeh/macro

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
  - javascript
  - css
  - html
  - [react](https://react.dev)
    - [vite](https://vite.dev)
    - [react-router](https://reactrouter.com)
    - [axios](https://axios-http.com/)