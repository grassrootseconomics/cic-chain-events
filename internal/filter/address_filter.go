package filter

import (
	"context"
	"sync"

	"github.com/grassrootseconomics/cic-chain-events/pkg/fetch"
	"github.com/zerodha/logf"
)

type (
	AddressFilterOpts struct {
		Cache         *sync.Map
		Logg          logf.Logger
		SystemAddress string
	}

	AddressFilter struct {
		cache         *sync.Map
		logg          logf.Logger
		systemAddress string
	}
)

func NewAddressFilter(o AddressFilterOpts) Filter {
	return &AddressFilter{
		cache:         o.Cache,
		logg:          o.Logg,
		systemAddress: o.SystemAddress,
	}
}

func (f *AddressFilter) Execute(_ context.Context, transaction fetch.Transaction) (bool, error) {
	if transaction.From.Address == f.systemAddress {
		return true, nil
	}

	if _, found := f.cache.Load(transaction.To.Address); found {
		return true, nil
	}

	return false, nil
}
