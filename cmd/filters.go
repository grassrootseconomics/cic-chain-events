package main

import "github.com/grassrootseconomics/cic-chain-events/internal/filter"

func initNoopFilter() filter.Filter {
	return filter.NewNoopFilter(filter.NoopFilterOpts{
		Logg: lo,
	})
}
