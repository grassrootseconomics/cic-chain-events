package syncer

import (
	"context"
	"time"

	"github.com/alitto/pond"
	"github.com/grassrootseconomics/cic-chain-events/internal/pipeline"
	"github.com/grassrootseconomics/cic-chain-events/internal/store"
	"github.com/jackc/pgx/v5"
	"github.com/zerodha/logf"
)

type JanitorOpts struct {
	BatchSize     uint64
	HeadBlockLag  uint64
	Logg          logf.Logger
	Pipeline      *pipeline.Pipeline
	Pool          *pond.WorkerPool
	Stats         *Stats
	Store         store.Store[pgx.Rows]
	SweepInterval time.Duration
}

type Janitor struct {
	batchSize     uint64
	headBlockLag  uint64
	pipeline      *pipeline.Pipeline
	logg          logf.Logger
	pool          *pond.WorkerPool
	stats         *Stats
	store         store.Store[pgx.Rows]
	sweepInterval time.Duration
}

func NewJanitor(o JanitorOpts) *Janitor {
	return &Janitor{
		batchSize:     o.BatchSize,
		headBlockLag:  o.HeadBlockLag,
		logg:          o.Logg,
		pipeline:      o.Pipeline,
		pool:          o.Pool,
		stats:         o.Stats,
		store:         o.Store,
		sweepInterval: o.SweepInterval,
	}
}

func (j *Janitor) Start(ctx context.Context) error {
	timer := time.NewTimer(j.sweepInterval)

	for {
		select {
		case <-ctx.Done():
			j.logg.Info("janitor: shutdown signal received")
			return nil
		case <-timer.C:
			j.logg.Debug("janitor: starting sweep")
			if err := j.QueueMissingBlocks(context.Background()); err != nil {
				j.logg.Error("janitor: queue missing blocks error", "error", err)
			}

			timer.Reset(j.sweepInterval)
		}
	}
}

// QueueMissingBlocks searches for missing block and queues the block for processing.
// It will run twice for a given search range and only after, raise the lower bound.
func (j *Janitor) QueueMissingBlocks(ctx context.Context) error {
	if j.stats.GetHeadCursor() == 0 {
		j.logg.Warn("janitor: (skipping) awaiting head synchronization")
		return nil
	}

	if j.pool.WaitingTasks() >= j.batchSize {
		j.logg.Warn("janitor: (skipping) avoiding queue pressure")
		return nil
	}

	lowerBound, upperBound, err := j.store.GetSearchBounds(
		ctx,
		j.batchSize,
		j.stats.GetHeadCursor(),
		j.headBlockLag,
	)
	if err != nil {
		return err
	}

	rows, err := j.store.GetMissingBlocks(ctx, lowerBound, upperBound)
	if err != nil {
		return err
	}

	rowsProcessed := 0
	for rows.Next() {
		var n uint64
		if err := rows.Scan(&n); err != nil {
			return err
		}

		j.pool.Submit(func() {
			if err := j.pipeline.Run(ctx, n); err != nil {
				j.logg.Error("janitor: pipeline run error", "error", err)
			}
		})
		rowsProcessed++
	}

	j.logg.Debug("janitor: missing blocks count", "count", rowsProcessed)
	if rowsProcessed == 0 {
		j.logg.Debug("janitor: rasing lower bound")
		j.stats.UpdateLowerBound(upperBound)
		j.store.SetSearchLowerBound(ctx, upperBound)
	}

	if rows.Err() != nil {
		return err
	}

	return nil
}
