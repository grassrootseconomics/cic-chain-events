package main

import (
	"sync"

	"github.com/grassrootseconomics/celoutils"
	"github.com/inethi/inethi-cic-chain-events/internal/filter"
	"github.com/inethi/inethi-cic-chain-events/internal/pub"
)

func initAddressFilter(celoProvider *celoutils.Provider, cache *sync.Map) filter.Filter {

	return filter.NewAddressFilter(filter.AddressFilterOpts{
		Logg: lo,
	})
}

func initTransferFilter(pub *pub.Pub) filter.Filter {
	return filter.NewTransferFilter(filter.TransferFilterOpts{
		Pub:  pub,
		Logg: lo,
	})

}
