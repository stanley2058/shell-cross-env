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

clean:
	go clean
