package pool

import (
	"context"
	"time"

	"github.com/alitto/pond"
)

type Opts struct {
	ConcurrencyFactor int
	PoolQueueSize     int
}

// NewPool creates a fixed size (and buffered) go routine worker pool.
func NewPool(ctx context.Context, o Opts) *pond.WorkerPool {
	return pond.New(
		o.ConcurrencyFactor,
		o.PoolQueueSize,
		pond.MinWorkers(o.ConcurrencyFactor),
		pond.IdleTimeout(time.Second*1),
		pond.Context(ctx),
	)
}
