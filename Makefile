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

test-pkg:
	TEST_GRAPHQL_ENDPOINT=https://rpc.alfajores.celo.grassecon.net/graphql go test -v -covermode atomic -coverprofile=covprofile ./pkg/...

migrate:
	tern migrate -c migrations/tern.conf

docker-clean:
	docker-compose down
	docker volume rm cic-chain-events_cic-indexer-pg cic-chain-events_cic-indexer-nats
