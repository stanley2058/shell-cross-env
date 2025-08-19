# Usage:
#   make           # build
#   make install   # build + install to ~/.local/bin (default)
#   make clean     # remove build cache (optional)

BINARY := shell-cross-env
PREFIX ?= $(HOME)/.local
BIN_DIR := $(PREFIX)/bin

.PHONY: all build install clean

all: build

build:
	go build

install: build
	mkdir -p "$(BIN_DIR)"
	cp "./$(BINARY)" "$(BIN_DIR)/"

build-linux_x64:
	GOOS=linux GOARCH=amd64 go build -o $(BINARY)-linux_x64

build-macos_arm64:
	GOOS=darwin GOARCH=arm64 go build -o $(BINARY)-macos_arm64

build-all: build-linux_x64 build-macos_arm64

clean:
	go clean
