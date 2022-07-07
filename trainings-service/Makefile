.PHONY: install test binary race cover run

CURRENT_DIR := $(shell pwd)
VERSION := $(shell git describe --always) 
GOLIST := $(shell go list -m)
BINARY := trainings-service 

run:
	./${BINARY}

binary:
	@go build -o ${BINARY} ${CURRENT_DIR}/cmd/srv 

test:
	@go test ./... -cover -count=1 -timeout 10s

test-unit:
	@go test ./... -run Unit -race -count=1

test-integration:
	@go test ./... -run Integration -race -count=1 -cover
