BIN := ./bin/app

.PHONY: run
run: ## Run the application
	clear
	go run cmd/main.go

.PHONY: build
build: ## Build the binary
	go build -o $(BIN) cmd/main.go

.PHONY: fmt
fmt: ## Format code
	go fmt ./...
	go mod tidy

.PHONY: proto
proto: ## Generate protobuf and gRPC code
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
