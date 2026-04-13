# base image
FROM golang:1.26-alpine AS go-builder
FROM node:24-alpine AS base
COPY --from=go-builder /usr/local/go /usr/local/go
ENV PATH="/root/go/bin:/usr/local/go/bin:${PATH}" PORT=8080 CGO_ENABLED=0 GOOS=linux
RUN go install github.com/air-verse/air@latest
WORKDIR /repo

# goose (migration tool)
FROM golang:1.26-alpine AS goose
RUN go install github.com/pressly/goose/v3/cmd/goose@latest