## This is a self-documented Makefile. For usage information, run `make help`:
##
## For more information, refer to https://suva.sh/posts/well-documented-makefiles/

include .bingo/Variables.mk

.PHONY: all build-server run help


##@ Building
run: $(BRA) ## Build and run web server on filesystem changes.
	$(BRA) run

build: build-server build-gencode ## build

build-server: ## Build FHUB server.
	@echo "build server"
	GOOS=linux GOARCH=amd64 CGO_ENABLED=1 $(GO) build -o ./bin/fhub-rest ./cmd/rest

build-gencode: ## Build FHUB gencode.
	@echo "build gencode"
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 $(GO) build -o ./bin/fhub-gencode ./cmd/gencode


##@ Helpers
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
