## NATS JetStream

### Server setup

- Enable `-js -sd`

## Stream setup

```go
_, err = js.AddStream(&nats.StreamConfig{
	Name:       streamName,
    // Remove from JS once Acked (Should not be used with 1 consumer which acts as a  relayer e.g. Benthos).
	// Retention:  nats.WorkQueuePolicy,
    // MaxAge allows us to replay it within 48 hrs
    MaxAge:     time.Hour * 48,
	Storage:    nats.FileStorage,
	Subjects:   []string{streamSubjects},
    // Sliding window dedup period.
	Duplicates: time.Minute * 5,
})
```

## Producer

```go
// nats.MsgId is the unique identifier for dedup
ctx.Publish("*", []byte("*"), nats.MsgId("*"))
```

## Consumer setup

- Explicit ACK
- Durable
- Deliver: all

### Benthos example

```toml
input:
  label: jetstream
  nats_jetstream:
    urls:
      - nats://127.0.0.1:4222
    subject: "*"
    durable: benthos
    deliver: all
  output:
  stdout:
    codec: lines
```

### Replay example

```bash
nats sub "*" --all --start-sequence=$N
```