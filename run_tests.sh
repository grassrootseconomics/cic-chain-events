#! /bin/bash
set -e

source .env.test && go test -v -covermode atomic -coverprofile=covprofile ./internal/...
