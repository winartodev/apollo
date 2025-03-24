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

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  run-http        Run the application"
	@echo "  migrate-create  Create a new migration (use NAME=<migration_name>)"
	@echo "  help            Display this help message"