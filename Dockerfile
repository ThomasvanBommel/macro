# frontend
FROM node:20-alpine AS frontend-builder
WORKDIR /app
COPY frontend/package*.json ./
RUN npm install
COPY frontend/ ./
RUN npm run build

# backend
FROM golang:1.26-alpine AS backend-builder
WORKDIR /app
COPY backend/go.mod backend/go.su[m] ./
RUN go mod download
COPY backend/ ./
COPY --from=frontend-builder /app/build ./frontend/build
RUN go build -o macro .

# production
FROM alpine:latest
WORKDIR /root
COPY --from=backend-builder /app/macro .
EXPOSE 8080
RUN adduser -D macrouser
USER macrouser
CMD ["./macro"]