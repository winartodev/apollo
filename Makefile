MIGRATION_EXT=sql
MIGRATION_PATH=core/database/migration/

DB_DRIVER=postgresql
DB_HOST=127.0.0.1
DB_PORT=5432
DB_NAME=apollo_db
DB_USER=apollo_user
DB_PASSWORD=apollo
DB_SSL=disable

run-http:
	go run cmd/http.go

download:
	go mod download

test:
	go test -v ./...

test-cover:
	go test -v -cover -race ./...

compose-up:
	docker-compose up -d

# create new migrations
migrate-create:
	@echo "Creating migration with NAME=${NAME}"
	migrate create -ext ${MIGRATION_EXT} -dir ${MIGRATION_PATH} ${NAME}

migrate-down:
	migrate -path ${MIGRATION_PATH} -database "$(DB_DRIVER)://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL)" -verbose down ${STEPS}

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  run-http        Run the application"
	@echo "  migrate-create  Create a new migration (use NAME=<migration_name>)"
	@echo "  help            Display this help message"
