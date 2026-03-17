FROM golang:1.26-alpine AS builder
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

FROM alpine:latest
COPY --from=builder /go/bin/goose /usr/local/bin/goose
WORKDIR /app
# ENTRYPOINT ["goose"]