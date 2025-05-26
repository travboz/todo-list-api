include .env # read from .env file

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

OUTPUT_BINARY=todo-list-api
OUTPUT_DIR=./bin
ENTRY_DIR=./cmd/api
SCRIPTS_DIR=./scripts

.PHONY: build run clean

## build: build application into executable in a 'bin' directory
build:
	@echo "Building..."
	@mkdir -p $(OUTPUT_DIR)
	@go build -o $(OUTPUT_DIR)/$(OUTPUT_BINARY) $(ENTRY_DIR)

## clean: delete all files within output directory
clean:
	@echo "Cleaning files..."
	@rm -rf $(OUTPUT_DIR)

seed:
	@echo "Seeding database..."
	@${SCRIPTS_DIR}/seed_users.sh
	@echo "Seeding complete."

## run: run the binary/server
run: build
	@$(OUTPUT_DIR)/$(OUTPUT_BINARY)


# Docker commands
.PHONY: compose/up compose/down connect-shell

## compose/up: start any docker containers using docker compose
compose/up:	
	@echo "Starting containers..."
	docker compose up -d

## compose/down: drop any docker containers using docker compose
compose/down:
	@echo "Stopping containers..."
	@docker compose down -v

## connect-shell: connect to the docker container which is running the mongo database
connect-shell:
	@echo "Connecting to postgres container via shell.."
	@docker exec -it ${DB_CONTAINER_NAME} /bin/bash

## db/access-port: shows the access port for the mongo db port
db/access-port:
	@echo ${DB_ACCESS_PORT}

# Helpers

## test: run all tests
.PHONY: test
test:
	go test -v -race -buildvcs ./...

## tidy: tidy modfiles and format .go files
.PHONY: tidy
tidy:
	go mod tidy -v
	go fmt ./...
