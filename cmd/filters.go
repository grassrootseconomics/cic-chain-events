package main

import (
	"github.com/grassrootseconomics/cic-chain-events/internal/filter"
)

func initAddressFilter() filter.Filter {
	return filter.NewAddressFilter(filter.AddressFilterOpts{
		Logg: lo,
	})
}

func initTransferFilter() filter.Filter {
	return filter.NewTransferFilter(filter.TransferFilterOpts{
		Logg: lo,
	})
}

func initNoopFilter() filter.Filter {
	return filter.NewNoopFilter(filter.NoopFilterOpts{
		Logg: lo,
	})
}
