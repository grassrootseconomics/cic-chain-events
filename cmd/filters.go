package main

import (
	"github.com/grassrootseconomics/cic-chain-events/internal/filter"
)

func initAddressFilter() filter.Filter {
	return filter.NewAddressFilter(filter.AddressFilterOpts{
		Logg: lo,
	})
}

func initDecodeFilter() filter.Filter {
	return filter.NewDecodeFilter(filter.DecodeFilterOpts{
		Logg: lo,
	})
}
