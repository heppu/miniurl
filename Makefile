.DEFAULT_GOAL := help

# Variables for text transformations.
COMMA := ,
SPACE := $(subst ,, )

# Tool related definitions.
GO_VERSION		:= 1.20.3
TOOLS_DIR		:= .tools
GO 				:= ${TOOLS_DIR}/go/go${GO_VERSION}
GOLANGCI_LINT	:= ${TOOLS_DIR}/github.com/golangci/golangci-lint/cmd/golangci-lint@v1.52.2

# Tool installation helpers.
MAJOR_VER	= $(firstword $(subst ., ,$(lastword $(subst @, ,${@}))))
LAST_PART	= $(notdir $(firstword $(subst @, ,${@})))
BIN_PATH	= ${PWD}/${@D}
BIN_NAME	= $(if $(filter ${LAST_PART},$(MAJOR_VER)),$(notdir ${BIN_PATH}),${LAST_PART})

${GO}: # Install required Go version
	@GOBIN=${PWD}/$(dir ${GO}) go install -mod=readonly golang.org/dl/go${GO_VERSION}@latest
	${GO} download

${GOLANGCI_LINT}: ${GO} # Install tools
	@mkdir -p ${BIN_PATH}
	@cd $(shell mktemp -d) && GOFLAGS='' GOBIN='${BIN_PATH}' ${PWD}/${GO} install ${@:${TOOLS_DIR}/%=%}
	@mv ${BIN_PATH}/${BIN_NAME} ${@}

.PHONY: help test benchmark fuzz lint build image go-version clean clean-tools

help: ## Show help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<recipe>\033[0m\n\nRecipes:\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-22s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

test: ${GO} ## Run tests
	${GO} test -v ./...

benchmark: ${GO} ## Run benchmarks
	${GO} test -run='^$$' -bench=.

fuzz: ${GO} ## Run fuzzer
	${GO} test -run=.

lint: ${GOLANGCI_LINT} ## Lint code
	@echo TODO

build: ${GO} ## Build binary
	@echo TODO

image: ## Build image
	@echo TODO

go-version: ## Print Go version
	@echo ${GO_VERSION}

clean: ## Remove all build and test artifacts
	rm -r target

clean-tools: ## Remove all tools
	rm -r .tools
