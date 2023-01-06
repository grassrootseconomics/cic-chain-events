package filter

import (
	"github.com/grassrootseconomics/cic-chain-events/internal/fetch"
	"github.com/zerodha/logf"
)

type TransferFilterOpts struct {
	Logg logf.Logger
}

type TransferFilter struct {
	logg logf.Logger
}

func NewTransferFilter(o TransferFilterOpts) Filter {
	return &TransferFilter{
		logg: o.Logg,
	}
}

func (f *TransferFilter) Execute(transaction fetch.Transaction) (bool, error) {
	switch transaction.InputData[:10] {
	case "0xa9059cbb":
		f.logg.Info("cUSD transfer", "block", transaction.Block.Number, "index", transaction.Index)
	case "0x23b872dd":
		f.logg.Info("cUSD transferFrom", "block", transaction.Block.Number, "index", transaction.Index)
	default:
		f.logg.Info("cUSD otherMethod", "block", transaction.Block.Number, "index", transaction.Index)
	}

	return true, nil
}
