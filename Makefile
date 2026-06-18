BINARY              := k8shell
BIN_DIR             := bin

GOTOOLCHAIN_VERSION := go1.25.0
GOPATH              := $(HOME)/go
TOOLCHAIN_DIR       := $(GOPATH)/pkg/mod/golang.org/toolchain@v0.0.1-$(GOTOOLCHAIN_VERSION).linux-amd64
GO                  := $(TOOLCHAIN_DIR)/bin/go

export GOROOT      := $(TOOLCHAIN_DIR)
export GOPATH
export GOTOOLCHAIN := $(GOTOOLCHAIN_VERSION)

.PHONY: build build-darwin build-darwin-amd64 build-darwin-arm64 clean

build:
	@mkdir -p $(BIN_DIR)
	@$(GO) build -o $(BIN_DIR)/$(BINARY) .

build-darwin-amd64:
	@mkdir -p $(BIN_DIR)
	@GOOS=darwin GOARCH=amd64 $(GO) build -o $(BIN_DIR)/$(BINARY)-darwin-amd64 .

build-darwin-arm64:
	@mkdir -p $(BIN_DIR)
	@GOOS=darwin GOARCH=arm64 $(GO) build -o $(BIN_DIR)/$(BINARY)-darwin-arm64 .

build-darwin: build-darwin-amd64 build-darwin-arm64

build-all: build build-darwin

clean:
	@rm -rf $(BIN_DIR)
