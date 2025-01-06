
BIN_NAME ?= ./build/app
SQLC_CONFIG ?= ./repository_impl/database/sqlc.yaml
SQLC_GENERATED ?= ./repository_impl/database/postgres/*.go

.PHONY: help
help:
	@echo "build: Build preprocessor service"
	@echo "run: Run preprocessor service"
	@echo "sqlc-clean: Cleanup Sqlc generated schemas and queries"
	@echo "sqlc-generate: Generate Sqlc schemas and queries"

.PHONY: build
build:
	@echo "Preprocessor Service Building..."
	@go build -o $(BIN_NAME) ./cmd/api/main.go

.PHONY: run
run:
	@echo "Preprocessor Service Running..."
	@$(BIN_NAME)

.PHONY: sqlc-generate
sqlc-generate:
	@echo "Generating sqlc schemas and queries..."
	@sqlc generate -f $(SQLC_CONFIG)
	@echo "Generated sqlc schemas and queries"

.PHONY: sqlc-clean
sqlc-clean:
	@echo "Cleaning sqlc generated schemas and queries..."
	@rm $(SQLC_GENERATED)
	@echo "Cleanedup sqlc generated schemas and queries"

