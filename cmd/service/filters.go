package main

import (
	"strings"
	"sync"

	"github.com/grassrootseconomics/cic-chain-events/internal/filter"
	"github.com/grassrootseconomics/cic-chain-events/internal/pub"
)

var (
	systemAddress = strings.ToLower("0x3D85285e39f05773aC92EAD27CB50a4385A529E4")
)

func initAddressFilter() filter.Filter {
	// TODO: Bootstrap addresses from smart contract
	// TODO: Add route to update cache
	cache := &sync.Map{}

	// Example bootstrap addresses
	cache.Store(strings.ToLower("0xB92463E2262E700e29c16416270c9Fdfa17934D7"), "TRNVoucher")
	cache.Store(strings.ToLower("0xf2a1fc19Ad275A0EAe3445798761FeD1Eea725d5"), "GasFaucet")
	cache.Store(strings.ToLower("0x1e041282695C66944BfC53cabce947cf35CEaf87"), "AddressIndex")

	return filter.NewAddressFilter(filter.AddressFilterOpts{
		Cache:         cache,
		Logg:          lo,
		SystemAddress: systemAddress,
	})
}

func initTransferFilter(pub *pub.Pub) filter.Filter {
	return filter.NewTransferFilter(filter.TransferFilterOpts{
		Pub:  pub,
		Logg: lo,
	})

}

func initGasGiftFilter(pub *pub.Pub) filter.Filter {
	return filter.NewGasFilter(filter.GasFilterOpts{
		Pub:           pub,
		Logg:          lo,
		SystemAddress: systemAddress,
	})
}

func initRegisterFilter(pub *pub.Pub) filter.Filter {
	return filter.NewRegisterFilter(filter.RegisterFilterOpts{
		Pub:  pub,
		Logg: lo,
	})
}