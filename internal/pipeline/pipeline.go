package pipeline

import (
	"context"

	"github.com/grassrootseconomics/cic-chain-events/internal/store"
	"github.com/grassrootseconomics/cic-chain-events/pkg/fetch"
	"github.com/grassrootseconomics/cic-chain-events/pkg/filter"
	"github.com/jackc/pgx/v5"
	"github.com/zerodha/logf"
)

type PipelineOpts struct {
	BlockFetcher fetch.Fetch
	Filters      []filter.Filter
	Logg         logf.Logger
	Store        store.Store[pgx.Rows]
}

type Pipeline struct {
	fetch   fetch.Fetch
	filters []filter.Filter
	logg    logf.Logger
	store   store.Store[pgx.Rows]
}

func NewPipeline(o PipelineOpts) *Pipeline {
	return &Pipeline{
		fetch:   o.BlockFetcher,
		filters: o.Filters,
		logg:    o.Logg,
		store:   o.Store,
	}
}

// Run is the task executor which runs in its own goroutine and does the following:
// 1. Fetches the block and all transactional data
// 2. Passes the block through all filters
// 3. Commits the block to store as successfully processed
//
// Note:
// - Blocks are processed atomically, a failure in-between will process the block from the start
// - Therefore, any side effect/event sink in the filter should support dedup
func (md *Pipeline) Run(ctx context.Context, blockNumber uint64) error {
	md.logg.Debug("pipeline: processing block", "block", blockNumber)
	fetchResp, err := md.fetch.Block(ctx, blockNumber)
	if err != nil {
		return err
	}

	for _, tx := range fetchResp.Data.Block.Transactions {
		for _, filter := range md.filters {
			next, err := filter.Execute(ctx, tx)
			if err != nil {
				return err
			}
			if !next {
				break
			}
		}
	}

	if err := md.store.CommitBlock(ctx, blockNumber); err != nil {
		return err
	}
	md.logg.Debug("pipeline: committed block", "block", blockNumber)

	return nil
}
