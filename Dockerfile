# frontend
FROM node:24-alpine AS frontend
WORKDIR /frontend
COPY frontend/package*.json ./
RUN npm install
COPY frontend/ .
RUN npm run build

# backend
FROM golang:1.26-alpine AS backend
RUN go install github.com/air-verse/air@latest
WORKDIR /repo
ENV PORT=8080 CGO_ENABLED=0 GOOS=linux
CMD ["air"]

# goose (migration tool)
FROM golang:1.26-alpine AS goose
RUN go install github.com/pressly/goose/v3/cmd/goose@latest