
BIN_NAME ?= ./build/app

.PHONY: all
all: build run


.PHONY: build
build:
	@echo "File Service Building..."
	@go build -o $(BIN_NAME) ./cmd/api/main.go

.PHONY: run
run:
	@echo "File Service Running..."
	@$(BIN_NAME)
