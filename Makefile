###############################################################################
###                                Variables                                  ###
###############################################################################

DOCKER := $(shell which docker)
BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
COMMIT := $(shell git log -1 --format='%H')

# Get version from git tag or branch+commit
ifeq (,$(VERSION))
  VERSION := $(shell git describe --exact-match 2>/dev/null)
  ifeq (,$(VERSION))
    VERSION := $(BRANCH)-$(COMMIT)
  endif
endif

# Build flags
ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=mini \
	-X github.com/cosmos/cosmos-sdk/version.AppName=minid \
	-X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
	-X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT)

BUILD_FLAGS := -ldflags '$(ldflags)'

# Protobuf
protoVer=0.14.0
protoImageName=ghcr.io/cosmos/proto-builder:$(protoVer)
protoImage=$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace $(protoImageName)

# Linting
golangci_lint_cmd=golangci-lint
golangci_version=v1.51.2

###############################################################################
###                                Targets                                    ###
###############################################################################

.PHONY: all install init clean build test lint proto-all

all: install

install: go.sum
	@echo "--> ensure dependencies have not been modified"
	@go mod verify
	@echo "--> installing minid"
	@go install $(BUILD_FLAGS) -mod=readonly ./cmd/minid

init:
	@echo "--> initializing chain"
	./scripts/init.sh

build:
	@echo "--> building minid binary"
	@go build $(BUILD_FLAGS) -mod=readonly -o bin/minid ./cmd/minid

clean:
	@echo "--> cleaning up build artifacts"
	rm -rf bin/
	rm -rf ~/.minid

###############################################################################
###                                Testing                                    ###
###############################################################################

test:
	@echo "--> running tests"
	@go test -v ./...

test-integration:
	@echo "--> running integration tests" 
	cd integration; go test -v ./...

###############################################################################
###                                Protobuf                                   ###
###############################################################################

proto-all: proto-format proto-lint proto-gen

proto-gen:
	@echo "Generating protobuf files..."
	@$(protoImage) sh ./scripts/protocgen.sh
	@go mod tidy

proto-format:
	@$(protoImage) find ./ -name "*.proto" -exec clang-format -i {} \;

proto-lint:
	@$(protoImage) buf lint proto/ --error-format=json

###############################################################################
###                                Linting                                    ###
###############################################################################

lint:
	@echo "--> running linter"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(golangci_version)
	@$(golangci_lint_cmd) run ./... --timeout 15m

lint-fix:
	@echo "--> running linter and fixing issues"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(golangci_version)
	@$(golangci_lint_cmd) run ./... --fix --timeout 15m

# Ensure go.sum is up to date
go.sum: go.mod
	@echo "--> ensure dependencies have not been modified"
	@go mod verify
	@go mod tidy

###############################################################################
###                                Helpers                                    ###
###############################################################################

help:
	@echo "Available targets:"
	@echo "  make install  - Install minid binary"
	@echo "  make init     - Initialize the chain"
	@echo "  make build    - Build minid binary"
	@echo "  make clean    - Clean build artifacts"
	@echo "  make test     - Run tests"
	@echo "  make lint     - Run linter"
	@echo "  make proto-all- Generate and lint protobuf"
	@echo "  make help     - Show this help message"

.DEFAULT_GOAL := help