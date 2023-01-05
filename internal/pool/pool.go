package pool

import (
	"time"

	"github.com/alitto/pond"
)

type Opts struct {
	ConcurrencyFactor int
	PoolQueueSize     int
}

func NewPool(o Opts) *pond.WorkerPool {
	return pond.New(
		o.ConcurrencyFactor,
		o.PoolQueueSize,
		pond.MinWorkers(o.ConcurrencyFactor),
		pond.IdleTimeout(time.Second*1),
	)
}
