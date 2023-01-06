FROM debian:11-slim

RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /cic-chain-events

COPY cic-chain-events .
COPY config.toml .
COPY queries.sql .
COPY LICENSE .
COPY migrations migrations/

CMD ["./cic-chain-events"]
