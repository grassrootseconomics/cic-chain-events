
## Prerequisites

- Linux OS (amd64) or Docker
- Postgres >= 14
- Celo geth with GraphQL API enabled
- NATS server with JetStream enabled

## Usage

The provided `docker-compose.yaml` is the fastest way to get up and running. Bring up the Postgres and NATS conatiners with `docker-compose up -d`

### 1. Run migrations

Run the SQL migrations inside the `migrations` folder with `psql` or [`tern`](https://github.com/jackc/tern) (recommended).

### 2. Update the config

The base config is described in `config.toml`. Values can be overriden with env variables e.g. to disable metrics, set `METRICS_GO_PROCESS=false`.

### 3. Start the service

**Compiling the binary**

Run `make build` or download pre-compiled binaries from the [releases](https://github.com/grassrootseconomics/cic-chain-events/releases) page.

Then start the service with `./cic-chain-events`

Optional flags:

- `-config` - `config.toml` file path
- `-debug` - Enable/disable debug level logs
- `-queries` - `queries.sql` file path

**Docker**

To pull the pre-built docker image:

`docker pull ghcr.io/grassrootseconomics/cic-chain-events/cic-chain-events:latest`

Or to build it:

`DOCKER_BUILDKIT=1 COMPOSE_DOCKER_CLI_BUILD=1 docker-compose -f docker-compose.build.yaml build --progress plain`

### 4. NATS JetStream consumer

A consumer with the following NATS JetStream config is required:

- Durable
- Stream: `CHAIN.*` (See `config.toml` for stream subjects)

[Benthos](https://benthos.dev) (Benthos can act as a JetStream consumer) example.

```yaml
# config.yaml
input:
  label: jetstream
  nats_jetstream:
    urls:
      - nats://127.0.0.1:4222
    subject: "CHAIN.*"
    durable: benthos
    deliver: all
output:
  stdout:
    codec: lines
```

Then run:

`benthos -c config.yaml`
