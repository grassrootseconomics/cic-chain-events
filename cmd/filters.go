package main

import (
	"strings"
	"sync"

	"github.com/grassrootseconomics/cic-chain-events/pkg/filter"
)

func initAddressFilter() filter.Filter {
	// TODO: Bootstrap addresses from smart contract
	// TODO: Add route to update cache
	cache := &sync.Map{}

	// Example bootstrap addresses
	cache.Store(strings.ToLower("0x617f3112bf5397D0467D315cC709EF968D9ba546"), "USDT")
	cache.Store(strings.ToLower("0x765DE816845861e75A25fCA122bb6898B8B1282a"), "cUSD")
	cache.Store(strings.ToLower("0xD8763CBa276a3738E6DE85b4b3bF5FDed6D6cA73"), "cEUR")

	return filter.NewAddressFilter(filter.AddressFilterOpts{
		Cache: cache,
		Logg:  lo,
	})
}

func initDecodeFilter() filter.Filter {
	js, err := initJetStream()
	if err != nil {
		lo.Fatal("filters: critical error loading jetstream", "error", err)
	}

	return filter.NewDecodeFilter(filter.DecodeFilterOpts{
		Logg:  lo,
		JSCtx: js,
	})
}
