package pool

import (
	"context"
	"time"

	"github.com/alitto/pond"
)

const (
	idleTimeout = 1 * time.Second
)

type Opts struct {
	Concurrency int
	QueueSize   int
}

// NewPool creates a fixed size (and buffered) go routine worker pool.
func NewPool(ctx context.Context, o Opts) *pond.WorkerPool {
	return pond.New(
		o.Concurrency,
		o.QueueSize,
		pond.MinWorkers(o.Concurrency),
		pond.IdleTimeout(idleTimeout),
		pond.Context(ctx),
	)
}
