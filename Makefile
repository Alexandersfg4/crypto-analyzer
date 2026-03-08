APP_NAME ?= crypto-analyzer
BIN_DIR ?= bin
BIN_PATH := $(BIN_DIR)/$(APP_NAME)
INSTALL_DIR ?= /usr/local/bin

.PHONY: help run build install

help:
	@echo "Targets:"
	@echo "  run      Run the CLI with current flags/env"
	@echo "  build    Build the binary into $(BIN_PATH)"
	@echo "  install  Install the binary into $(INSTALL_DIR)"

run:
	go run cmd/main.go

build:
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_PATH) cmd/main.go

install: build
	install -m 0755 $(BIN_PATH) $(INSTALL_DIR)/$(APP_NAME)
