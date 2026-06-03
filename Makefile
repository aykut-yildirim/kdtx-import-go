SHELL := /bin/bash
PROJECT_NAME := kx_import
BUILD_TIMESTAMP := $(shell date +%Y-%m-%dT%H:%M:%S)
PROTO_SRC := ../ktx-grpc-server/proto
.PHONY: run build test clean format lint type-check setup proto deps tidy check \
        local-serve local-api local-worker prod-api prod-worker prod-stop local-test

format:
	@echo "Formatting Go code..."
	@go fmt ./...
	@goimports -w .

lint:
	@echo "Running golangci-lint..."
	@golangci-lint run ./...

type-check:
	@echo "Type checking..."
	@go build ./...

build:
	@echo "Building $(PROJECT_NAME)..."
	@go build -o bin/$(PROJECT_NAME) ./cmd/api

run:
	@echo "Running API..."
	@go run ./cmd/api

test:
	@echo "Running tests..."
	@go test ./... -v

clean:
	@echo "Cleaning..."
	@rm -rf bin/

setup:
	@echo "Setting up Go modules..."
	@go mod tidy

deps:
	@echo "Downloading dependencies..."
	@go mod download

tidy:
	@go mod tidy

check:
	@echo "Full check pipeline..."
	@go fmt ./...
	@golangci-lint run ./...
	@go test ./...

proto:
	@echo "Generating Go gRPC files..."
	@rm -rf internal/proto
	@mkdir -p internal/proto

	protoc -I=$(PROTO_SRC) \
		--go_out=internal/proto --go_opt=paths=source_relative \
		--go-grpc_out=internal/proto --go-grpc_opt=paths=source_relative \
		$(PROTO_SRC)/kudtax/v1/import_worker.proto \
		$(PROTO_SRC)/kudtax/v1/portal.proto \
		$(PROTO_SRC)/kudtax/v1/common.proto
	@echo "Proto generation completed."

local-serve:
	@echo "Starting local services... ($(BUILD_TIMESTAMP))"
	@ENV=local docker compose -f docker-compose_serve.yml up --build -d

local-api:
	@echo "Starting local API... ($(BUILD_TIMESTAMP))"
	@ENV=local docker compose -f docker-compose_api.yml up --build -d

local-worker:
	@echo "Starting local Worker... ($(BUILD_TIMESTAMP))"
	@ENV=local docker compose -f docker-compose_worker.yml up --build -d

prod-api:
	@echo "Starting PROD API... ($(BUILD_TIMESTAMP))"
	@ENV=prod docker compose -f docker-compose_api.yml up --build -d

prod-worker:
	@echo "Starting PROD Worker... ($(BUILD_TIMESTAMP))"
	@ENV=prod docker compose -f docker-compose_worker.yml up --build -d

prod-stop:
	@echo "Stopping PROD services..."
	@ENV=prod docker compose down

local-test:
	@echo "Running local tests..."
	@ENV=local go run ./test/test_file.go