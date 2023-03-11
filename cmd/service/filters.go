package main

import (
	"strings"
	"sync"

	"github.com/grassrootseconomics/cic-chain-events/internal/filter"
	"github.com/grassrootseconomics/cic-chain-events/internal/pub"
)

var (
	systemAddress string
)

func initAddressFilter() filter.Filter {
	// TODO: Temporary shortcut
	systemAddress = ko.MustString("chain.system_address")

	// TODO: Bootstrap addresses from smart contract
	// TODO: Add route to update cache
	cache := &sync.Map{}

	cache.Store(strings.ToLower(ko.MustString("chain.token_index_address")), "TokenIndex")
	cache.Store(strings.ToLower(ko.MustString("chain.gas_faucet_address")), "GasFaucet")
	cache.Store(strings.ToLower(ko.MustString("chain.user_index_address")), "UserIndex")

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
