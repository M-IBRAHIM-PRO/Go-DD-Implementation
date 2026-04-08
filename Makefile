# Makefile
APP_NAME=library_system
BINARY=bin/$(APP_NAME)

DB_DSN ?= $(shell grep '^DB_DSN=' .env | cut -d '=' -f2-)

# ==============================
# Core Commands
# ==============================

run: build
	bash scripts/run.sh

build:
	bash scripts/build.sh

# ==============================
# Migrations
# ==============================

migrate-up:
	DB_DSN="$(DB_DSN)" bash scripts/migrate.sh up

migrate-down:
	DB_DSN="$(DB_DSN)" bash scripts/migrate.sh down

migrate-force:
	DB_DSN="$(DB_DSN)" bash scripts/migrate.sh force $(VERSION)

# ==============================
# Dev Tools
# ==============================

lint:
	golangci-lint run

test:
	go test ./...

# ==============================
# Utilities
# ==============================

clean:
	rm -rf bin/

.PHONY: run build migrate-up migrate-down migrate-force lint test clean