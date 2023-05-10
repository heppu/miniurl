.DEFAULT_GOAL := help

# Common build variables
TARGET_DIR	:= target
TOOLS_DIR	:= .tools
GO_VERSION	:= 1.20.3
GO 			:= ${TOOLS_DIR}/go/go${GO_VERSION}

# Variables for text transformations.
COMMA := ,
SPACE := $(subst ,, )

${GO}: # Install required Go version
	@GOBIN=${PWD}/$(dir ${GO}) go install -mod=readonly golang.org/dl/go${GO_VERSION}@latest
	${GO} download

.PHONY: help tidy unit-test integration-test benchmark fuzz build image go-version clean clean-tools

help: ## Show help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<recipe>\033[0m\n\nRecipes:\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-22s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

tidy: ${GO} ## Tidy Go modules
	${GO} mod tidy

unit-test: ${GO} ## Run unit tests
	mkdir -p ${TARGET_DIR}
	${GO} test -v -cover -coverprofile=${TARGET_DIR}/cover.out ./...

integration-test: ${GO} ## Run integration tests
	@echo TODO

benchmark: ${GO} ## Run benchmarks
	@echo TODO

fuzz: ${GO} ## Run fuzzy tests
	@echo TODO

build: ${GO} ## Build binary
	@echo TODO

image: ## Build image
	@echo TODO

go-version: ## Print Go version
	@echo ${GO_VERSION}

clean: ## Remove all build and test artifacts
	rm -r ${TARGET_DIR}

clean-tools: ## Remove all tools
	rm -r ${TOOLS_DIR}
