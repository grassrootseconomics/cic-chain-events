package filter

import (
	"context"

	"github.com/celo-org/celo-blockchain/common"
	"github.com/grassrootseconomics/cic-chain-events/internal/events"
	"github.com/grassrootseconomics/cic-chain-events/pkg/fetch"
	"github.com/grassrootseconomics/w3-celo-patch"
	"github.com/zerodha/logf"
)

const (
	registerEventSubject = "CHAIN.register"
)

var (
	addSig = w3.MustNewFunc("add(address)", "bool")
)

type RegisterFilterOpts struct {
	EventEmitter events.EventEmitter
	Logg         logf.Logger
}

type RegisterFilter struct {
	eventEmitter events.EventEmitter
	logg         logf.Logger
}

func NewRegisterFilter(o RegisterFilterOpts) Filter {
	return &RegisterFilter{
		eventEmitter: o.EventEmitter,
		logg:         o.Logg,
	}
}

func (f *RegisterFilter) Execute(_ context.Context, transaction fetch.Transaction) (bool, error) {
	if len(transaction.InputData) < 10 {
		return true, nil
	}

	if transaction.InputData[:10] == "0x0a3b0a4f" {
		var address common.Address

		if err := addSig.DecodeArgs(w3.B(transaction.InputData), &address); err != nil {
			return false, err
		}

		addEvent := &events.MinimalTxInfo{
			Block:           transaction.Block.Number,
			ContractAddress: transaction.To.Address,
			To:              transaction.To.Address,
			TxHash:          transaction.Hash,
			TxIndex:         transaction.Index,
		}

		if transaction.Status == 1 {
			addEvent.Success = true
		}

		if err := f.eventEmitter.Publish(
			registerEventSubject,
			transaction.Hash,
			addEvent,
		); err != nil {
			return false, err
		}

		return true, nil
	}
	
	return true, nil
}
