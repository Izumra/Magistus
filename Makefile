include .env

SHELL:=/usr/bin/bash

migrations_status:
	goose -dir=$(MIGRATIONS_DIR) status

migrations_up:
	goose -dir=$(MIGRATIONS_DIR) up

migrations_reset:
	goose -dir=$(MIGRATIONS_DIR) reset

add_sql_migration:
	goose -dir=${MIGRATIONS_DIR} create ${MIGRATION_NAME} sql

run: 
	go run bot/cmd/bot/main.go --config=bot/config/local.yaml