.PHONY: install test binary race cover 

CURRENT_DIR := $(shell pwd)
VERSION := $(shell git describe --always)
GOLIST := $(shell go list -m)
BINARY := service-discovery 

binary:
	@go build -o ${BINARY} ${CURRENT_DIR}/cmd/srv

test:
	@go test ./... -cover -count=1 

test-unit:
	@go test ./... -run Unit -v -race -count=1

test-integration:
	@go test ./... -run Integration -v -race -count=1
