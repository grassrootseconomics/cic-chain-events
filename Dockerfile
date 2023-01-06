FROM golang:1.19-bullseye as build

WORKDIR /build
COPY . .
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o cic-chain-events -ldflags="-s -w" cmd/*.go

FROM debian:bullseye-slim

ENV DEBIAN_FRONTEND=noninteractive
RUN set -x && \
    apt-get update && \
    apt-get install -y ca-certificates && \
    rm -rf /var/cache/apt/archives /var/lib/apt/lists/*

WORKDIR /service

COPY --from=build /build/cic-chain-events .
COPY config.toml .
COPY queries.sql .
COPY LICENSE .
COPY migrations migrations/

CMD ["./cic-chain-events"]

EXPOSE 8080