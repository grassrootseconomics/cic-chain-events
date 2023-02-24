package filter

import (
	"context"

	"github.com/celo-org/celo-blockchain/common/hexutil"
	"github.com/grassrootseconomics/cic-chain-events/internal/events"
	"github.com/grassrootseconomics/cic-chain-events/pkg/fetch"
	"github.com/zerodha/logf"
)

const (
	gasFilterEventSubject = "CHAIN.gas"
)

type GasFilterOpts struct {
	EventEmitter  events.EventEmitter
	Logg          logf.Logger
	SystemAddress string
}

type GasFilter struct {
	eventEmitter  events.EventEmitter
	logg          logf.Logger
	systemAddress string
}

func NewGasFilter(o GasFilterOpts) Filter {
	return &GasFilter{
		eventEmitter:  o.EventEmitter,
		logg:          o.Logg,
		systemAddress: o.SystemAddress,
	}
}

func (f *GasFilter) Execute(_ context.Context, transaction fetch.Transaction) (bool, error) {
	transferValue, err := hexutil.DecodeUint64(transaction.Value)
	if err != nil {
		return false, err
	}

	// TODO: This is a temporary shortcut to gift gas. Switch to gas faucet contract.
	if transaction.From.Address == f.systemAddress && transferValue > 0 {
		transferEvent := &events.MinimalTxInfo{
			Block:   transaction.Block.Number,
			To:      transaction.To.Address,
			TxHash:  transaction.Hash,
			TxIndex: transaction.Index,
			Value:   transferValue,
		}

		if transaction.Status == 1 {
			transferEvent.Success = true
		}

		if err := f.eventEmitter.Publish(
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
