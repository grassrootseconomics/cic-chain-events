## Functionality

## Filters

Filters are initialized in `cmd/filters.go` and implemented in `internal/filters/*.go` folder. You will need to modify these files to suite your indexing needs.

The existing implementation demo's tracking Celo stables transfer events and gives a rough idea on how to write filters. The final filter should always emit an event to NATS JetStream.

## Syncers

### Head syncer

The head syncer processes newely produced blocks independently by connection to the geth websocket endpoint.

### Janitor

The janitor syncer checks for missing (blocks) gaps in the commited block sequence and queues them for processing. It can also function as a historical syncer too process older blocks.

With the default `config.toml`, The janitor can process around 950-1000 blocks/min.

_Ordering_

Missed/historical blocks are not guaranteed to be processed in order, however a low concurrency setting would somewhat give an "in-order" behaviour (not to be relied upon in any case).

## Block fetchers

The default GraphQL block fetcher is the recommended fetcher. An experimental RPC fetcher implementation is also provided as an example.

## Pipeline

The pipeline fetches a whole block with its full transaction and receipt objects, executes all loaded filters serially and finally commits the block value to the db. Blocks are processed atomically by the pipeline; a failure in one of the filters will trigger the janitor to re-queue the block and process the block again.

## Store

The postgres store keeps track of commited blocks and syncer curosors. Schema:

- The `blocks` table keeps track of processed blocks.
- The `syncer_meta` table keeps track of the lower_bound cursor. Below the lower_bound cursor, all blocks are guarnteed to have been processsed hence it is safe to trim the `blocks` table below that pointer.
