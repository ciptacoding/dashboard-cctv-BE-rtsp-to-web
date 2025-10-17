# Makefile untuk CCTV Monitoring Backend

.PHONY: help build run stop restart logs clean test

# Default target
help:
	@echo "Available commands:"
	@echo "  make build      - Build Docker images"
	@echo "  make run        - Start all services"
	@echo "  make stop       - Stop all services"
	@echo "  make restart    - Restart all services"
	@echo "  make logs       - View logs"
	@echo "  make clean      - Stop and remove volumes"
	@echo "  make test       - Test API endpoints"
	@echo "  make db         - Connect to database"
	@echo "  make migrate    - Run migrations"

# Build Docker images
build:
	docker-compose build --no-cache

# Start all services
run:
	docker-compose up -d
	@echo "✓ Services started!"
	@echo "  Backend API: http://localhost:8080"
	@echo "  RTSPtoWeb: http://localhost:8083"
	@echo "  PostgreSQL: localhost:5432"

# Stop all services
stop:
	docker-compose down
	@echo "✓ Services stopped!"

# Restart all services
restart:
	docker-compose restart
	@echo "✓ Services restarted!"

# View logs
logs:
	docker-compose logs -f

# View backend logs only
logs-backend:
	docker-compose logs -f backend

# Stop and remove volumes (WARNING: data will be lost!)
clean:
	docker-compose down -v
	@echo "✓ Services stopped and volumes removed!"

# Connect to PostgreSQL
db:
	docker exec -it cctv_postgres psql -U cctv_user -d cctv_monitoring

# Run migrations manually
migrate:
	docker-compose exec backend go run cmd/api/main.go

# Test health check
test-health:
	@echo "Testing health endpoint..."
	@curl -s http://localhost:8080/health | python -m json.tool

# Test login
test-login:
	@echo "Testing login..."
	@curl -s -X POST http://localhost:8080/api/v1/auth/login \
		-H "Content-Type: application/json" \
		-d '{"username":"admin","password":"admin123"}' | python -m json.tool

# Development mode - run without Docker
dev:
	go run cmd/api/main.go

# Install Go dependencies
deps:
	go mod download
	go mod tidy

# Format code
fmt:
	go fmt ./...