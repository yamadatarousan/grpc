GOBIN := $(HOME)/go/bin
BUF := $(GOBIN)/buf
GOCACHE := $(CURDIR)/.cache/go-build
GOMODCACHE := $(CURDIR)/.cache/go-mod
BUF_CACHE_DIR := $(CURDIR)/.cache/buf
PATH := $(GOBIN):$(CURDIR)/web/node_modules/.bin:$(PATH)

export GOCACHE
export GOMODCACHE
export BUF_CACHE_DIR
export PATH

.PHONY: setup tidy deps web-build go-build build run test

setup:
	go install github.com/bufbuild/buf/cmd/buf@v1.58.0
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.11
	go install connectrpc.com/connect/cmd/protoc-gen-connect-go@v1.19.1

generate:
	$(BUF) generate

tidy:
	go mod tidy

deps:
	cd web && npm install

web-build:
	cd web && npm run build

go-build:
	go build ./...

build: generate deps web-build go-build

run:
	go run ./cmd/server

test:
	go test ./...
