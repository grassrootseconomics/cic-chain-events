package filter

import (
	"context"
	"sync"

	"github.com/grassrootseconomics/cic-chain-events/pkg/fetch"
	"github.com/zerodha/logf"
)

type AddressFilterOpts struct {
	Cache *sync.Map
	Logg  logf.Logger
}

type AddressFilter struct {
	cache *sync.Map
	logg  logf.Logger
}

func NewAddressFilter(o AddressFilterOpts) Filter {
	return &AddressFilter{
		cache: o.Cache,
		logg:  o.Logg,
	}
}

func (f *AddressFilter) Execute(_ context.Context, transaction *fetch.Transaction) (bool, error) {
	if _, found := f.cache.Load(transaction.To.Address); found {
		return true, nil
	}

	return false, nil
}
