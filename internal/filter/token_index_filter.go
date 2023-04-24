package filter

import (
	"context"
	"strings"
	"sync"

	"github.com/celo-org/celo-blockchain/common"
	"github.com/celo-org/celo-blockchain/common/hexutil"
	"github.com/grassrootseconomics/celoutils"
	"github.com/grassrootseconomics/cic-chain-events/internal/pub"
	"github.com/grassrootseconomics/cic-chain-events/pkg/fetch"
	"github.com/grassrootseconomics/w3-celo-patch"
	"github.com/zerodha/logf"
)

const (
	tokenIndexFilterEventSubject = "CHAIN.tokenAdded"
)

var (
	addSig = w3.MustNewFunc("add(address)", "bool")
)

type (
	TokenIndexFilterOpts struct {
		Cache *sync.Map
		Logg  logf.Logger
		Pub   *pub.Pub
	}

	TokenIndexFilter struct {
		pub   *pub.Pub
		cache *sync.Map
		logg  logf.Logger
	}
)

func NewTokenIndexFilter(o TokenIndexFilterOpts) Filter {
	return &TokenIndexFilter{
		cache: o.Cache,
		logg:  o.Logg,
		pub:   o.Pub,
	}
}

func (f *TokenIndexFilter) Execute(_ context.Context, transaction *fetch.Transaction) (bool, error) {
	if len(transaction.InputData) < 10 {
		return true, nil
	}

	if transaction.InputData[:10] == "0x0a3b0a4f" {
		var address common.Address

		if err := addSig.DecodeArgs(w3.B(transaction.InputData), &address); err != nil {
			return false, err
		}

		f.cache.Store(strings.ToLower(address.Hex()), transaction.Hash)

		addEvent := &pub.MinimalTxInfo{
			Block:           transaction.Block.Number,
			ContractAddress: celoutils.ChecksumAddress(transaction.To.Address),
			Timestamp:       hexutil.MustDecodeUint64(transaction.Block.Timestamp),
			To:              address.Hex(),
			TxHash:          transaction.Hash,
			TxIndex:         transaction.Index,
			TXType:          "tokenAdd",
		}

		if transaction.Status == 1 {
			addEvent.Success = true
		}

		if err := f.pub.Publish(
			tokenIndexFilterEventSubject,
			transaction.Hash,
			addEvent,
		); err != nil {
			return false, err
		}
	}
	return true, nil
}
