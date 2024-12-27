
BIN_NAME ?= ./build/app

.PHONY: all
all: build run

.PHONY: build
build:
	@echo "Preprocessor Service Building..."
	@go build -o $(BIN_NAME) ./cmd/api/main.go

.PHONY: run
run:
	@echo "Preprocessor Service Running..."
	@$(BIN_NAME)

.PHONY: help
help:
	@echo "build: Build preprocessor service"
	@echo "run:   Run preprocessor service"
