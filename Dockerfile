# build
FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY backend/. .
# RUN go mod download
RUN go build -o macro .


# run
FROM alpine:latest
WORKDIR /root
COPY --from=builder /app/macro .
EXPOSE 8080
CMD ["./macro"]