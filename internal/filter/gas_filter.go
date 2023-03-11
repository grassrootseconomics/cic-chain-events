package filter

import (
	"context"

	"github.com/celo-org/celo-blockchain/common/hexutil"
	"github.com/grassrootseconomics/celoutils"
	"github.com/grassrootseconomics/cic-chain-events/internal/pub"
	"github.com/grassrootseconomics/cic-chain-events/pkg/fetch"
	"github.com/zerodha/logf"
)

const (
	gasFilterEventSubject = "CHAIN.gas"
)

type (
	GasFilterOpts struct {
		Logg          logf.Logger
		Pub           *pub.Pub
		SystemAddress string
	}

	GasFilter struct {
		logg          logf.Logger
		pub           *pub.Pub
		systemAddress string
	}
)

func NewGasFilter(o GasFilterOpts) Filter {
	return &GasFilter{
		logg:          o.Logg,
		pub:           o.Pub,
		systemAddress: o.SystemAddress,
	}
}

func (f *GasFilter) Execute(_ context.Context, transaction *fetch.Transaction) (bool, error) {
	transferValue, err := hexutil.DecodeUint64(transaction.Value)
	if err != nil {
		return false, err
	}

	// TODO: This is a temporary shortcut to gift gas. Switch to gas faucet contract.
	if transaction.From.Address == f.systemAddress && transferValue > 0 {
		transferEvent := &pub.MinimalTxInfo{
			Block:   transaction.Block.Number,
			To:      celoutils.ChecksumAddress(transaction.To.Address),
			TxHash:  transaction.Hash,
			TxIndex: transaction.Index,
			Value:   transferValue,
		}

		if transaction.Status == 1 {
			transferEvent.Success = true
		}

		if err := f.pub.Publish(
			gasFilterEventSubject,
			transaction.Hash,
			transferEvent,
		); err != nil {
			return false, err
		}

		return true, nil
	}

	return true, nil
}
