## Functionality

### Head syncer

Opens a websocket connection and processes live transactions.

### Janitor

Periodically checks for missed (and historical) blocks missed by the head syncer and queues them for processing. A gap range is processed twice to guarantee there is no missing block.

### Pipeline

Fetches a block and executes the filters in serial order for every transaction in the block before finally committing the block to the store.

### Filter

Processes a transaction and passes it on to the next filter or terminates the pipeline for that transaction if it is irrelevant.

### Store schema

- The `blocks` table keeps track of processed blocks.
- The `syncer_meta` table keeps track of the lower_bound cursor. Below the lower_bound cursor, all blocks are guarnteed to have been processsed hence it is safe to trim the `blocks` table below that pointer.

### GraphQL

- Fetches a block (and some of its header details), transactions and transaction receipts embedded within the transaction object in a single call.

### NATS JetStream

- The final filter will emit an event to JetStream.

To view/debug the JetStream messages, you can use [Benthos](https://benthos.dev)

With a config like:

```yaml
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

## Caveats

- Blocks are not guaranteed to be processed in order, however a low concurrency setting would somewhat give an "in-order" behaviour (not to be relied upon in any case).
