package pipeline

import (
	"github.com/grassrootseconomics/cic-chain-events/internal/fetch"
	"github.com/grassrootseconomics/cic-chain-events/internal/filter"
	"github.com/grassrootseconomics/cic-chain-events/internal/store"
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

// Run is the task executor which fetches and processes a block and its transactions through the pipeline filters
func (md *Pipeline) Run(blockNumber uint64) error {
	fetchResp, err := md.fetch.Block(blockNumber)
	if err != nil {
		md.logg.Error("pipeline block fetch error", "error", err)
		return err
	}

	for _, tx := range fetchResp.Data.Block.Transactions {
		for _, filter := range md.filters {
			next, err := filter.Execute(tx)
			if err != nil {
				md.logg.Error("pipeline run error", "error", err)
				return err
			}
			if !next {
				break
			}
		}
	}

	if err := md.store.CommitBlock(blockNumber); err != nil {
		return err
	}
	md.logg.Debug("successfully commited block", "block", blockNumber)

	return nil
}
