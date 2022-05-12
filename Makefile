cnf ?= ./deployments/.env
include $(cnf)
export $(shell sed 's/=.*//' $(cnf))

.PHONY: help

help:
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help

server: ## Run application server
	go run ./cmd/server/main.go

client: ## Run application client
	go run ./cmd/client/main.go

test: ## Run tests
	go test ./...
