package filter

import (
	"context"
	"sync"

	celo "github.com/grassrootseconomics/cic-celo-sdk"
	"github.com/grassrootseconomics/cic-chain-events/internal/fetch"
	"github.com/zerodha/logf"
)

type AddressFilterOpts struct {
	Cache        *sync.Map
	CeloProvider *celo.Provider
	Logg         logf.Logger
}

type AddressFilter struct {
	cache *sync.Map
	logg  logf.Logger
}

func NewAddressFilter(o AddressFilterOpts) Filter {
	// TODO: Bootstrap addresses from registry smart contract
	return &AddressFilter{
		cache: o.Cache,
		logg:  o.Logg,
	}
}

func (f *AddressFilter) Execute(_ context.Context, transaction fetch.Transaction) (bool, error) {
	if _, found := f.cache.Load(transaction.To.Address); found {
		return true, nil
	}

	return false, nil
}
