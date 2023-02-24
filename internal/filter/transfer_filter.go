package filter

import (
	"context"
	"math/big"

	"github.com/celo-org/celo-blockchain/common"
	"github.com/grassrootseconomics/cic-chain-events/internal/events"
	"github.com/grassrootseconomics/cic-chain-events/pkg/fetch"
	"github.com/grassrootseconomics/w3-celo-patch"
	"github.com/zerodha/logf"
)

const (
	transferFilterEventSubject = "CHAIN.transfer"
)

var (
	transferSig     = w3.MustNewFunc("transfer(address, uint256)", "bool")
	transferFromSig = w3.MustNewFunc("transferFrom(address, address, uint256)", "bool")
	mintToSig       = w3.MustNewFunc("mintTo(address, uint256)", "bool")
)

type TransferFilterOpts struct {
	EventEmitter events.EventEmitter
	Logg         logf.Logger
}

type TransferFilter struct {
	eventEmitter events.EventEmitter
	logg         logf.Logger
}

func NewTransferFilter(o TransferFilterOpts) Filter {
	return &TransferFilter{
		eventEmitter: o.EventEmitter,
		logg:         o.Logg,
	}
}

func (f *TransferFilter) Execute(_ context.Context, transaction fetch.Transaction) (bool, error) {
	if len(transaction.InputData) < 10 {
		return true, nil
	}

	switch transaction.InputData[:10] {
	case "0xa9059cbb":
		var (
			to    common.Address
			value big.Int
		)

		if err := transferSig.DecodeArgs(w3.B(transaction.InputData), &to, &value); err != nil {
			return false, err
		}

		f.logg.Debug("transfer_filter: new reg", "transfer", to)

		transferEvent := &events.MinimalTxInfo{
			Block:           transaction.Block.Number,
			From:            transaction.From.Address,
			To:              to.Hex(),
			ContractAddress: transaction.To.Address,
			TxHash:          transaction.Hash,
			TxIndex:         transaction.Index,
			Value:           value.Uint64(),
		}

		if transaction.Status == 1 {
			transferEvent.Success = true
		}

		if err := f.eventEmitter.Publish(
			transferFilterEventSubject,
			transaction.Hash,
			transferEvent,
		); err != nil {
			return false, err
		}

		return true, nil
	case "0x23b872dd":
		var (
			from  common.Address
			to    common.Address
			value big.Int
		)

		if err := transferFromSig.DecodeArgs(w3.B(transaction.InputData), &from, &to, &value); err != nil {
			return false, err
		}

		f.logg.Debug("transfer_filter: new reg", "transferFrom", to)

		transferEvent := &events.MinimalTxInfo{
			Block:           transaction.Block.Number,
			From:            from.Hex(),
			To:              to.Hex(),
			ContractAddress: transaction.To.Address,
			TxHash:          transaction.Hash,
			TxIndex:         transaction.Index,
			Value:           value.Uint64(),
		}

		if transaction.Status == 1 {
			transferEvent.Success = true
		}

		if err := f.eventEmitter.Publish(
			transferFilterEventSubject,
			transaction.Hash,
			transferEvent,
		); err != nil {
			return false, err
		}

		return true, nil
	case "0x449a52f8":
		var (
			to    common.Address
			value big.Int
		)

		if err := mintToSig.DecodeArgs(w3.B(transaction.InputData), &to, &value); err != nil {
			return false, err
		}

		f.logg.Debug("transfer_filter: new reg", "mintTo", to)

		transferEvent := &events.MinimalTxInfo{
			Block:           transaction.Block.Number,
			From:            transaction.From.Address,
			To:              to.Hex(),
			ContractAddress: transaction.To.Address,
			TxHash:          transaction.Hash,
			TxIndex:         transaction.Index,
			Value:           value.Uint64(),
		}

		if transaction.Status == 1 {
			transferEvent.Success = true
		}

		if err := f.eventEmitter.Publish(
			transferFilterEventSubject,
			transaction.Hash,
			transferEvent,
		); err != nil {
			return false, err
		}

		return true, nil
	default:
		return true, nil
	}
}
