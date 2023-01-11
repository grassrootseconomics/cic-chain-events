package filter

import (
	"context"

	"github.com/grassrootseconomics/cic-chain-events/internal/fetch"
	"github.com/zerodha/logf"
)

type NoopFilterOpts struct {
	Logg logf.Logger
}

type NoopFilter struct {
	logg logf.Logger
}

func NewNoopFilter(o NoopFilterOpts) Filter {
	return &NoopFilter{
		logg: o.Logg,
	}
}

func (f *NoopFilter) Execute(ctx context.Context, transaction fetch.Transaction) (bool, error) {
	f.logg.Debug("noop filter", "block", transaction.Block.Number, "index", transaction.Index)
	return true, nil
}
