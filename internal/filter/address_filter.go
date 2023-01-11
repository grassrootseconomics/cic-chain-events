package filter

import (
	"context"

	"github.com/grassrootseconomics/cic-chain-events/internal/fetch"
	"github.com/zerodha/logf"
)

const (
	cUSD = "0x765de816845861e75a25fca122bb6898b8b1282a"
)

type AddressFilterOpts struct {
	Logg logf.Logger
}

type AddressFilter struct {
	logg logf.Logger
}

func NewAddressFilter(o AddressFilterOpts) Filter {
	return &AddressFilter{
		logg: o.Logg,
	}
}

func (f *AddressFilter) Execute(ctx context.Context, transaction fetch.Transaction) (bool, error) {
	if transaction.To.Address == cUSD {
		return true, nil
	}

	return false, nil
}
