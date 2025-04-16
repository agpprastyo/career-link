# Migration directory
MIGRATION_DIR=migrations

# Load environment variables from .env
include .env
export

# Default migration name if not specified
NAME?=migration

.PHONY: migrate-new migrate-up migrate-down migrate-down-all migrate-force migrate-version migrate-create-dir

# Create migrations directory if it doesn't exist
migrate-create-dir:
	@mkdir -p $(MIGRATION_DIR)

# Create a new migration file
migrate-new: migrate-create-dir
	@echo "Creating migration files for $(NAME)..."
	@migrate create -ext sql -dir $(MIGRATION_DIR) -seq $(NAME)

# Apply all pending migrations
migrate-up:
	@echo "Applying migrations..."
	@migrate -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)" -path $(MIGRATION_DIR) up

# Rollback the last migration
migrate-down:
	@echo "Rolling back last migration..."
	@migrate -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)" -path $(MIGRATION_DIR) down 1

# Rollback all migrations
migrate-down-all:
	@echo "Rolling back all migrations..."
	@migrate -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)" -path $(MIGRATION_DIR) down

# Force to specific version
migrate-force:
	@read -p "Force version to: " VERSION; \
	migrate -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)" -path $(MIGRATION_DIR) force $$VERSION

# Show current migration version
migrate-version:
	@echo "Current migration version:"
	@migrate -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)" -path $(MIGRATION_DIR) version


# Database seeders
db-seed:
	@echo "Seeding database..."
	@go run cmd/seeder/main.go


# Start monitoring stack
start-monitoring:
	docker-compose -f docker-compose.monitoring.yml up -d

# Stop monitoring stack
stop-monitoring:
	docker-compose -f docker-compose.monitoring.yml down

.PHONY: start-monitoring stop-monitoring
