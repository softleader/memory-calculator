##@ General

help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-25s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

bootstrap:  ## Bootstrap project
	go mod download

govet:	## Run go vet
	go vet

gofmt:	## Run gofmt
	gofmt -s -w .

test: ## Run tests
	go test ./...

build: bootstrap govet gofmt ## Build app
	 go build -o ./build/memory-calculator

