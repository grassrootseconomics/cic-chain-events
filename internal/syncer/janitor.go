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
		case <-timer.C:
			j.logg.Debug("janitor starting sweep")
			if err := j.QueueMissingBlocks(); err != nil {
				j.logg.Error("janitor error", "error", err)
			}

			timer.Reset(j.sweepInterval)
		case <-ctx.Done():
			j.logg.Debug("janitor shutdown signal received")
			return nil
		}
	}
}

func (j *Janitor) QueueMissingBlocks() error {
	if j.stats.GetHeadCursor() == 0 {
		j.logg.Debug("janitor waiting for head synchronization")
		return nil
	}

	if j.pool.WaitingTasks() >= j.batchSize {
		j.logg.Debug("janitor skipping potential queue pressure")
		return nil
	}

	lowerBound, upperBound, err := j.store.GetSearchBounds(
		j.batchSize,
		j.stats.GetHeadCursor(),
		j.headBlockLag,
	)
	if err != nil {
		return err
	}
	j.logg.Debug("janitor search bounds", "lower_bound", lowerBound, "upper_bound", upperBound)

	rows, err := j.store.GetMissingBlocks(lowerBound, upperBound)
	if err != nil {
		return err
	}

	rowsProcessed := 0
	for rows.Next() {
		var n uint64
		if err := rows.Scan(&n); err != nil {
			return err
		}

		j.logg.Debug("submitting missing block for processing", "block", n)
		j.pool.Submit(func() {
			if err := j.pipeline.Run(n); err != nil {
				j.logg.Error("pipeline run error", "error", err)
			}
		})

		rowsProcessed++
	}

	j.logg.Debug("janitor missing block count", "count", rowsProcessed)
	if rowsProcessed == 0 {
		j.logg.Debug("no missing blocks, rasing lower bound")
		j.stats.UpdateLowerBound(upperBound)
		j.store.SetSearchLowerBound(upperBound)
	}

	if rows.Err() != nil {
		return err
	}

	return nil
}
