package main

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/celo-org/celo-blockchain/common"
	"github.com/grassrootseconomics/celoutils"
	"github.com/grassrootseconomics/cic-chain-events/internal/filter"
	"github.com/grassrootseconomics/cic-chain-events/internal/pub"
	"github.com/grassrootseconomics/w3-celo-patch"
	"github.com/grassrootseconomics/w3-celo-patch/module/eth"
	"github.com/grassrootseconomics/w3-celo-patch/w3types"
)

func initAddressFilter(celoProvider *celoutils.Provider, cache *sync.Map) filter.Filter {
	var (
		tokenIndexEntryCount big.Int
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	registryMap, err := celoProvider.RegistryMap(ctx, celoutils.HexToAddress(ko.MustString("chain.registry_address")))
	if err != nil {
		lo.Fatal("init: critical error creating address filter", "error", err)
	}

	for k, v := range registryMap {
		cache.Store(strings.ToLower(v.Hex()), k)
	}

	if err := celoProvider.Client.CallCtx(
		ctx,
		eth.CallFunc(w3.MustNewFunc("entryCount()", "uint256"), registryMap[celoutils.TokenIndex]).Returns(&tokenIndexEntryCount),
	); err != nil {
		lo.Fatal("init: critical error creating address filter", "error", err)
	}

	calls := make([]w3types.Caller, tokenIndexEntryCount.Int64())
	tokenAddresses := make([]common.Address, tokenIndexEntryCount.Int64())

	entrySig := w3.MustNewFunc("entry(uint256 _idx)", "address")

	// TODO: There is a 5MB limit to a RPC batch call size.
	// Test if 10k entries will raise an error (future proofed for a lot of years)
	for i := 0; i < int(tokenIndexEntryCount.Int64()); i++ {
		calls[i] = eth.CallFunc(entrySig, registryMap[celoutils.TokenIndex], new(big.Int).SetInt64(int64(i))).Returns(&tokenAddresses[i])
	}

	if err := celoProvider.Client.CallCtx(
		ctx,
		calls...,
	); err != nil {
		lo.Fatal("init: critical error creating address filter", "error", err)
	}

	for i, v := range tokenAddresses {
		cache.Store(strings.ToLower(v.Hex()), fmt.Sprintf("TOKEN_%d", i))
	}

	return filter.NewAddressFilter(filter.AddressFilterOpts{
		Cache: cache,
		Logg:  lo,
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
		Pub:  pub,
		Logg: lo,
	})
}

func initRegisterFilter(pub *pub.Pub) filter.Filter {
	return filter.NewRegisterFilter(filter.RegisterFilterOpts{
		Pub:  pub,
		Logg: lo,
	})
}

func initApproveFilter(pub *pub.Pub) filter.Filter {
	return filter.NewApproveFilter(filter.ApproveFilterOpts{
		Pub:  pub,
		Logg: lo,
	})
}

func initTokenIndexFilter(cache *sync.Map, pub *pub.Pub) filter.Filter {
	return filter.NewTokenIndexFilter(filter.TokenIndexFilterOpts{
		Cache: cache,
		Pub:   pub,
		Logg:  lo,
	})
}
