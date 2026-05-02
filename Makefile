include .env
export

export PROJECT_ROOT=$(shell pwd)

run:
	clear
	go run cmd/main.go

fmt:
	go fmt ./...
	go mod tidy

proto:
	rm -rf pb/*
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
		--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
		proto/*.proto
	@for f in pb/*.go; do \
		awk '{ \
			if (/^package / && prev !~ /^\/\/ Package pb/) { print "// Package pb - сгенерированные прото файлы" } \
			if (/^package /) { print "package pb" } else { print } \
			prev=$$0 \
		}' "$$f" > "$$f.tmp" && mv "$$f.tmp" "$$f"; \
	done

env-up:
	@docker compose up -d wallet-flow-postgres

env-down:
	@docker compose down wallet-flow-postgres

env-cleanup:
	@read -p "Очистить все данные окружения? (y/N): " ans; \
	if [ "$$ans" = "y" ] || [ "$$ans" = "Y" ]; then \
		docker compose down wallet-flow-postgres && \
		rm -rf out/pgdata && \
		echo "Файлы окружения удалены"; \
	else \
		echo "Очистка окружения отменена"; \
	fi

env-port-forward:
	@docker compose up -d port-forwarder

env-port-close:
	@docker compose down port-forwarder

migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "Отсутсвует необходимый параметр seq. Пример: make migrate-create seq=1"; \
		exit 1; \
	fi;\
	docker compose run --rm wallet-flow-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"

migrate-up:
	@make migrate-action action=up

migrate-down:
	@make migrate-action action=down

migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "Отсутсвует необходимый параметр action. Пример: make migrate-action action=up"; \
		exit 1; \
	fi; \
	docker compose run --rm wallet-flow-migrate \
		-path /migrations \
		-database postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@wallet-flow-postgres:5432/$(POSTGRES_DB)?sslmode=disable \
		"$(action)"