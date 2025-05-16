#.DEFAULT_GOAL := start
SHELL=/bin/bash
NOW = $(shell date +"%Y%m%d%H%M%S")
UID = $(shell id -u)
PWD = $(shell pwd)

.PHONY: help
help: ## prints all targets available and their description
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'


.PHONY: build
build: ## Target to build all commands in cmd/*/*
	@echo "Building all commands..."
	@mkdir -p bin
	@for item in cmd/*/*; do \
		echo "Compiling $$item..."; \
		go build -o "bin/$$(basename "$$item")" "./$$item"; \
	done
	@echo "All commands built successfully."

.PHONY: clean
clean: ## Target to clean up built binaries
	@echo "Cleaning up binaries..."
	@rm -rf bin
	@echo "Binaries removed."