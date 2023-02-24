package main

import (
	"strings"
	"sync"

	"github.com/grassrootseconomics/cic-chain-events/internal/filter"
	"github.com/nats-io/nats.go"
)

func initAddressFilter() filter.Filter {
	// TODO: Bootstrap addresses from smart contract
	// TODO: Add route to update cache
	cache := &sync.Map{}

	// Example bootstrap addresses
	cache.Store(strings.ToLower("0x54c8D8718Ea9E7b2b4542e630fd36Ccab32cE74E"), "BABVoucher")
	cache.Store(strings.ToLower("0xdD4F5ea484F6b16f031eF7B98F3810365493BC20"), "GasFaucet")

	return filter.NewAddressFilter(filter.AddressFilterOpts{
		Cache: cache,
		Logg:  lo,
	})
}

func initTransferFilter(jsCtx nats.JetStreamContext) filter.Filter {
	return filter.NewTransferFilter(filter.TransferFilterOpts{
		Logg:  lo,
		JSCtx: jsCtx,
	})

}

func initGasGiftFilter(jsCtx nats.JetStreamContext) filter.Filter {
	return filter.NewGasFilter(filter.GasFilterOpts{
		Logg:  lo,
		JSCtx: jsCtx,
	})
}
