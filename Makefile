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
		echo "✳️ Compiling $$item..."; \
		if [[ "$$item" =~ "main.go" ]]; then \
          dir_name=$$(dirname "$$item"); \
          new_binary_name=$$(basename "$$dir_name"); \
		  echo "⚠️ Found 'main.go' in '$$item' changing name of bin. New bin name: $$new_binary_name"; \
          go build -o "bin/$$(basename "$$new_binary_name")" "./$$item"; \
		  echo "✅ Compiled"; \
        else \
          go build -o "bin/$$(basename "$$item")" "./$$item"; \
 		  echo "✅ Compiled"; \
	    fi \
	done
	@echo "All commands built successfully."

.PHONY: clean
clean: ## Target to clean up built binaries
	@echo "Cleaning up binaries..."
	@rm -rf bin
	@echo "Binaries removed."