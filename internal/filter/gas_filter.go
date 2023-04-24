package filter

import (
	"context"

	"github.com/celo-org/celo-blockchain/common"
	"github.com/celo-org/celo-blockchain/common/hexutil"
	"github.com/grassrootseconomics/celoutils"
	"github.com/grassrootseconomics/cic-chain-events/internal/pub"
	"github.com/grassrootseconomics/cic-chain-events/pkg/fetch"
	"github.com/grassrootseconomics/w3-celo-patch"
	"github.com/zerodha/logf"
)

const (
	gasEventSubject = "CHAIN.gas"
)

var (
	giveToSig = w3.MustNewFunc("giveTo(address)", "uint256")
)

type (
	GasFilterOpts struct {
		Logg logf.Logger
		Pub  *pub.Pub
	}

	GasFilter struct {
		logg logf.Logger
		pub  *pub.Pub
	}
)

func NewGasFilter(o GasFilterOpts) Filter {
	return &GasFilter{
		logg: o.Logg,
		pub:  o.Pub,
	}
}

func (f *GasFilter) Execute(_ context.Context, transaction *fetch.Transaction) (bool, error) {
	if len(transaction.InputData) < 10 {
		return true, nil
	}

	if transaction.InputData[:10] == "0x63e4bff4" {
		var address common.Address

		if err := giveToSig.DecodeArgs(w3.B(transaction.InputData), &address); err != nil {
			return false, err
		}

		giveToEvent := &pub.MinimalTxInfo{
			Block:           transaction.Block.Number,
			ContractAddress: celoutils.ChecksumAddress(transaction.To.Address),
			Timestamp:       hexutil.MustDecodeUint64(transaction.Block.Timestamp),
			To:              address.Hex(),
			TxHash:          transaction.Hash,
			TxIndex:         transaction.Index,
			TXType:          "gas",
		}

		if transaction.Status == 1 {
			giveToEvent.Success = true
		}

		if err := f.pub.Publish(
			gasEventSubject,
			transaction.Hash,
			giveToEvent,
		); err != nil {
			return false, err
		}

		return true, nil
	}

	return true, nil
}
