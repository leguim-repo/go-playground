#.DEFAULT_GOAL := start
SHELL=/bin/bash
NOW = $(shell date +"%Y%m%d%H%M%S")
UID = $(shell id -u)
PWD = $(shell pwd)
BUILD_FOLDER=build
BUILD_FOLDER_LINUX=build-linux
BUILD_FOLDER_RPI=build-rpi
BUILD_FOLDER_MACOS_INTEL=build-macos-intel
BUILD_FOLDER_MACOS_SILICON=build-macos-silicon

DOCKER_COMPOSE_DIR = TheDocker
DOCKER_COMPOSE = docker compose -f $(DOCKER_COMPOSE_DIR)/docker-compose.yml


.PHONY: help
help: ## Prints all targets available and their description
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'


define build-target
	@echo "Building all commands for $(1)..."
	@mkdir -p $(2)
	@for item in cmd/*/*; do \
		echo "✳️ Compiling $$item..."; \
		if [[ "$$item" =~ "main.go" ]]; then \
			dir_name=$$(dirname "$$item"); \
			new_binary_name=$$(basename "$$dir_name"); \
			echo "⚠️ Found 'main.go' in '$$item', changing bin name to: $$new_binary_name"; \
			$(3) go build -o "$(2)/$$(basename $$new_binary_name)" "./$$item"; \
			echo "✅ Compiled"; \
		else \
			$(3) go build -o "$(2)/$$(basename $$item)" "./$$item"; \
			echo "✅ Compiled"; \
		fi \
	done
	@echo "All commands built successfully for $(1)."
	@echo "Build folder: $(2)"
	@echo ""

endef


.PHONY: build
build: ## Target to build all commands in cmd/*/* for current system
	$(call build-target,local,${BUILD_FOLDER},)


.PHONY:	build-linux
build-linux: ## Target to build all commands in cmd/*/* for linux
	$(call build-target,linux,${BUILD_FOLDER_LINUX},GOOS=linux GOARCH=arm64)


.PHONY:	build-rpi
build-rpi: ## Target to build all commands in cmd/*/* for rpi
	$(call build-target,rpi,${BUILD_FOLDER_RPI},GOOS=linux GOARCH=arm)


.PHONY:	build-macos-intel
build-macos-intel: ## Target to build all commands in cmd/*/* for macos intel
	$(call build-target,macos-intel,${BUILD_FOLDER_MACOS_INTEL},GOOS=darwin GOARCH=amd64)


.PHONY:	build-macos-silicon
build-macos-silicon: ## Target to build all commands in cmd/*/* for macos apple silicon
	$(call build-target,macos-silicon,${BUILD_FOLDER_MACOS_SILICON},GOOS=darwin GOARCH=arm64)


.PHONY: build-all-arch
build-all-arch: build build-linux build-macos-intel build-macos-silicon build-rpi ## Target to build all architectures

.PHONY: clean
clean: ## Target to clean up built binaries
	@echo "Cleaning up binaries..."
	@rm -rf ${BUILD_FOLDER}
	@rm -rf ${BUILD_FOLDER_LINUX}
	@rm -rf ${BUILD_FOLDER_RPI}
	@rm -rf ${BUILD_FOLDER_MACOS_INTEL}
	@rm -rf ${BUILD_FOLDER_MACOS_SILICON}
	@echo "Build directory removed."

.PHONY: unit-tests
unit-tests: ## Launch all unit tests found in modules
	@go test ./... -v

.PHONY: docker-up
docker-up: ## Up The Docker
	@$(DOCKER_COMPOSE) up -d


.PHONY: docker-down
docker-down: ## Down The Docker
	@$(DOCKER_COMPOSE) down

.PHONY: docker-restart
docker-restart: ## Restart The Docker
	@$(DOCKER_COMPOSE) restart

.PHONY: docker-logs
docker-logs: ## View The Docker logs
	@$(DOCKER_COMPOSE) logs -f

.PHONY: docker-ps
docker-ps: ## View state of The Docker
	@$(DOCKER_COMPOSE) ps

.PHONY: docker-clean
docker-clean: ## Clean The Docker
	@$(DOCKER_COMPOSE) down -v --remove-orphans

.PHONY: docker-rebuild
docker-rebuild: ## Rebuild The Docker
	@$(DOCKER_COMPOSE) up -d --build

.PHONY: docker-reset
docker-reset: docker-clean docker-rebuild ## Reset The Docker
