.PHONY: help dev run build clean

help:
	@echo "Available targets:"
	@echo "  dev            - Start the development environment using Docker Compose."
	@echo "  run            - Run the production Docker container."
	@echo "  build          - Build the production Docker image, macro-production."
	@echo "  frontend/build - Build the frontend assets using a Node.js Docker container."
	@echo "  cleanup        - Clean up build artifacts and stop Docker containers."

dev:
	docker compose up --build

fmt:
	docker run --rm \
		-v "$(PWD)/backend":/app \
		-w /app golang:1.26-alpine \
		gofmt -s -w .

run: build
	docker run -it --rm -p 8080:8080 macro-production

build:
	docker build -t macro-production .

frontend/build:
	docker run -it --rm \
		-v "$(PWD)/frontend":/app \
		-w /app node:20-alpine \
		sh -c "npm install && npm run build && chown -R $(shell id -u):$(shell id -g) ."

clean:
	rm -rf frontend/build frontend/node_modules
	docker compose down --remove-orphans