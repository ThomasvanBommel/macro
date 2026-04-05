# frontend
FROM node:24-alpine AS frontend-builder
WORKDIR /frontend
COPY frontend/package*.json ./
RUN npm install
COPY frontend/ .
RUN npm run build

# backend
FROM golang:1.26-alpine AS backend-builder
WORKDIR /repo
RUN go env -w GOCACHE=/go-cache GOMODCACHE=/go-mod-cache
COPY ./backend/go.mod ./backend/go.sum ./
RUN --mount=type=cache,target=/go-mod-cache go mod download
COPY . .
COPY --from=frontend-builder /frontend/build ./backend/frontend/build
ENV CGO_ENABLED=0 GOOS=linux
RUN --mount=type=cache,target=/go-cache \
    --mount=type=cache,target=/go-mod-cache \
    cd backend && go build \
    -ldflags="-s -w" \
    -tags='no_postgres no_clickhouse no_mssql no_mysql prod' \
    -o /macro

# fullstack (production)
FROM gcr.io/distroless/static-debian12:latest AS distroless
COPY --from=backend-builder /macro ./macro
ENV GIN_MODE=release PORT=8080
EXPOSE 8080
CMD ["/macro"]

# backend dev
FROM golang:1.26-alpine AS backend-dev
RUN go install github.com/air-verse/air@latest
WORKDIR /repo
ENV PORT=8080 CGO_ENABLED=0 GOOS=linux
CMD ["air"]

# goose (migration tool)
FROM golang:1.26-alpine AS goose
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# production
FROM gcr.io/distroless/static-debian12:latest AS production
COPY ./backend/macro /macro
ENV GIN_MODE=release PORT=8080
EXPOSE 8080
CMD ["/macro"]