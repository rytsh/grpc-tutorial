BIN_DIR = bin
PROTO_DIR = proto
SERVER_DIR = server
CLIENT_DIR = client

SHELL := bash
SHELL_VERSION = $(shell echo $$BASH_VERSION)
UNAME := $(shell uname -s)
VERSION_AND_ARCH = $(shell uname -rm)
ifeq ($(UNAME),Darwin)
	OS = macos ${VERSION_AND_ARCH}
else ifeq ($(UNAME),Linux)
	OS = linux ${VERSION_AND_ARCH}
else
$(error OS not supported by this Makefile)
endif
PACKAGE = $(shell head -1 go.mod | awk '{print $$2}')
HELP_CMD = grep -E '^[a-zA-Z_-]+:.*?\#\# .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?\#\# "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
RM_F_CMD = rm -f
RM_RF_CMD = ${RM_F_CMD} -r
SERVER_BIN = ${SERVER_DIR}
CLIENT_BIN = ${CLIENT_DIR}

.DEFAULT_GOAL := help
.PHONY: greet greet-build greet-proto calculator calculator-build calculator-proto help about
project := calculator greet

all: $(project) ## Generate Pbs and build

calculator: calculator-proto calculator-build ## Generate Pbs and build for calculator
greet: greet-proto greet-build ## Generate Pbs and build for greet

$(foreach item,$(project),$(item)-build):
	go build -o ${BIN_DIR}/$(subst -build,,$@)/${SERVER_BIN} ./$(subst -build,,$@)/${SERVER_DIR}
	go build -o ${BIN_DIR}/$(subst -build,,$@)/${CLIENT_BIN} ./$(subst -build,,$@)/${CLIENT_DIR}

$(foreach item,$(project),$(item)-proto):
	@test -d $(subst -proto,,$@) || (echo "\033[31mDirectory $(subst -proto,,$@) doesn't exist\033[0m" && false)
	protoc -I$(subst -proto,,$@)/${PROTO_DIR} --go_opt=module=${PACKAGE} --go_out=. --go-grpc_opt=module=${PACKAGE} --go-grpc_out=. $(subst -proto,,$@)/${PROTO_DIR}/*.proto

test: all ## Launch tests
	go test ./...

clean: clean_calculator ## Clean generated files
	${RM_F_CMD} ssl/*.crt
	${RM_F_CMD} ssl/*.csr
	${RM_F_CMD} ssl/*.key
	${RM_F_CMD} ssl/*.pem
	${RM_RF_CMD} ${BIN_DIR}

clean_calculator: ## Clean generated files for calculator
	${RM_F_CMD} calculator/${PROTO_DIR}/*.pb.go

rebuild: clean all ## Rebuild the whole project

about: ## Display info related to the build
	@echo "OS: ${OS}"
	@echo "Shell: ${SHELL} ${SHELL_VERSION}"
	@echo "Protoc version: $(shell protoc --version)"
	@echo "Go version: $(shell go version)"
	@echo "Go package: ${PACKAGE}"
	@echo "Openssl version: $(shell openssl version)"

help: ## Show this help
	@${HELP_CMD}
