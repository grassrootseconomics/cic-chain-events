package pool

import (
	"context"
	"time"

	"github.com/alitto/pond"
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
		pond.IdleTimeout(time.Second*1),
		pond.Context(ctx),
	)
}
