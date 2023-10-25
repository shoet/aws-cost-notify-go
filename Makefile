.PHONY: help build generate
.DEFAULT_GOAL := help

// TODO: 修正
build: ## Build app
	cd src && \
	GOOS=linux GOARCH=amd64 go build -o handler && \
	zip handler.zip handler

generate: ## Generate codes
	go generate ./...

help: ## Show options
	@grep -E '^[a-zA-Z_]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
