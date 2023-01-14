SHELL := /bin/bash
BIN := cic-chain-events

.PHONY: build

build:
	CGO_ENABLED=1 go build -v -ldflags="-s -w" -o ${BIN} cmd/*.go

run:
	CGO_ENABLED=1 go run -ldflags="-s -w" cmd/*.go

mod:
	go mod tidy
	go mod verify

test:
	source .env.test
	go test -v -covermode atomic -coverprofile=covprofile ./internal/...