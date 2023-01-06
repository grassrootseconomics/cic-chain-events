package filter

import (
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

func (f *NoopFilter) Execute(transaction fetch.Transaction) (bool, error) {
	f.logg.Debug("noop filter", "block", transaction.Block.Number, "tx", transaction.Hash)
	return true, nil
}