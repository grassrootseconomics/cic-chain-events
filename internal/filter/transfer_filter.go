package filter

import (
	"context"
	"math/big"

	"github.com/celo-org/celo-blockchain/common"
	"github.com/celo-org/celo-blockchain/common/hexutil"
	"github.com/grassrootseconomics/celoutils"
	"github.com/grassrootseconomics/cic-chain-events/internal/pub"
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

type (
	TransferFilterOpts struct {
		Logg logf.Logger
		Pub  *pub.Pub
	}

	TransferFilter struct {
		logg logf.Logger
		pub  *pub.Pub
	}
)

func NewTransferFilter(o TransferFilterOpts) Filter {
	return &TransferFilter{
		logg: o.Logg,
		pub:  o.Pub,
	}
}

func (f *TransferFilter) Execute(_ context.Context, transaction *fetch.Transaction) (bool, error) {
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

		transferEvent := &pub.MinimalTxInfo{
			Block:           transaction.Block.Number,
			From:            celoutils.ChecksumAddress(transaction.From.Address),
			Timestamp:       hexutil.MustDecodeUint64(transaction.Block.Timestamp),
			To:              to.Hex(),
			ContractAddress: celoutils.ChecksumAddress(transaction.To.Address),
			TxHash:          transaction.Hash,
			TxIndex:         transaction.Index,
			TXType:          "transfer",
			Value:           value.Uint64(),
		}

		if transaction.Status == 1 {
			transferEvent.Success = true
		}

		if err := f.pub.Publish(
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

		transferEvent := &pub.MinimalTxInfo{
			Block:           transaction.Block.Number,
			From:            from.Hex(),
			Timestamp:       hexutil.MustDecodeUint64(transaction.Block.Timestamp),
			To:              to.Hex(),
			ContractAddress: celoutils.ChecksumAddress(transaction.To.Address),
			TxHash:          transaction.Hash,
			TxIndex:         transaction.Index,
			TXType:          "transferFrom",
			Value:           value.Uint64(),
		}

		if transaction.Status == 1 {
			transferEvent.Success = true
		}

		if err := f.pub.Publish(
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

		transferEvent := &pub.MinimalTxInfo{
			Block:           transaction.Block.Number,
			From:            celoutils.ChecksumAddress(transaction.From.Address),
			Timestamp:       hexutil.MustDecodeUint64(transaction.Block.Timestamp),
			To:              to.Hex(),
			ContractAddress: celoutils.ChecksumAddress(transaction.To.Address),
			TxHash:          transaction.Hash,
			TxIndex:         transaction.Index,
			TXType:          "mintTo",
			Value:           value.Uint64(),
		}

		if transaction.Status == 1 {
			transferEvent.Success = true
		}

		if err := f.pub.Publish(
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
