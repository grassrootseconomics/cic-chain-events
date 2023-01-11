## Requirements

- Celo (geth) node with GraphQL enabled
- Postgres 14+

## Running

### 1. Run migrations

Run the migrations inside the `migrations` folder.

### 2. Update the config

The base config is described in `config.toml`. Values can be overriden with env variables e.g. to disable metrics, set `METRICS_GO_PROCESS=false`.

### 3. Start the service:

**Compiling**:

- Requires CGO_ENABLED=1
- Prebuilt binaries (for amd64 only) available on the releases page

**Docker**:

- `docker pull ghcr.io/grassrootseconomics/cic-chain-events/cic-chain-events:latest`

After compiling or within a Docker container:

`$ ./cic-chain-events`

Optional flags:

- `-config` - `config.toml` file path
- `-debug` - Enable/disable debug level logs
- `-queries` - `queries.sql` file path

