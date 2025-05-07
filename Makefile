SHELL := /bin/bash

NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m

.PHONY: check-golangci lint

BASE := $(shell pwd)
BUILD_TIME := $(shell date +%FT%T%z)
# APP_VERSION := ${BUILD_TIME}.${GIT_COMMIT}

GOLANGCI_CMD := $(shell command -v golangci-lint 2> /dev/null)

check-golangci:
ifndef GOLANGCI_CMD
	$(error "Please install golangci linters from https://golangci-lint.run/usage/install/")
else
	@echo -e "$(OK_COLOR)golangci lint exists$(NO_COLOR)"
endif

CMD_MIGRATE := $(shell command -v goose 2> /dev/null)

check-migration:
ifndef CMD_MIGRATE
	$(error "goose (migration) is not installed, see: https://github.com/pressly/goose")
else
	@echo -e "$(OK_COLOR)goose (migration) exists$(NO_COLOR)"
endif

YQ_CMD := $(shell command -v yq 2> /dev/null)

check-yq:
ifndef YQ_CMD
	$(error "Please install yq YAML processor from https://github.com/mikefarah/yq/releases")
else
	@echo -e "$(OK_COLOR)yq exists$(NO_COLOR)"
endif

lint: check-golangci
	@echo -e "$(OK_COLOR)==> linting projects$(NO_COLOR)..."
	@golangci-lint run --fix
	@echo -e "$(OK_COLOR)==> done, all ok$(NO_COLOR)..."

# === get db cred
CONFIG_FILE := ./config.yml

MIGRATIONS_DIR := ./internal/database/migrations
SQLITE_PATH := $(shell echo $(shell yq ".db.sqlite_db_path" $(CONFIG_FILE)) | sed 's/^\.\.\///')

# Create a new SQL migration file
migrate-create:
ifndef NAME
	$(error Usage: make migrate-create NAME=your_migration_name)
endif
	goose -dir $(MIGRATIONS_DIR) create $(NAME) sql

# Apply all pending migrations
migrate-up:
	goose -dir $(MIGRATIONS_DIR) sqlite3 $(SQLITE_PATH) up

# Roll back the last migration
migrate-down:
	goose -dir $(MIGRATIONS_DIR) sqlite3 $(SQLITE_PATH) down

# Roll back all migrations
migrate-reset:
	goose -dir $(MIGRATIONS_DIR) sqlite3 $(SQLITE_PATH) reset

# Check migration status
migrate-status:
	goose -dir $(MIGRATIONS_DIR) sqlite3 $(SQLITE_PATH) status