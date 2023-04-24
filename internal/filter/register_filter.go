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
	registerEventSubject = "CHAIN.register"
)

var (
	registerSig = w3.MustNewFunc("register(address)", "")
)

type (
	RegisterFilterOpts struct {
		Logg logf.Logger
		Pub  *pub.Pub
	}

	RegisterFilter struct {
		logg logf.Logger
		pub  *pub.Pub
	}
)

func NewRegisterFilter(o RegisterFilterOpts) Filter {
	return &RegisterFilter{
		logg: o.Logg,
		pub:  o.Pub,
	}
}

func (f *RegisterFilter) Execute(_ context.Context, transaction *fetch.Transaction) (bool, error) {
	if len(transaction.InputData) < 10 {
		return true, nil
	}

	if transaction.InputData[:10] == "0x4420e486" {
		var address common.Address

		if err := registerSig.DecodeArgs(w3.B(transaction.InputData), &address); err != nil {
			return false, err
		}

		addEvent := &pub.MinimalTxInfo{
			Block:           transaction.Block.Number,
			ContractAddress: celoutils.ChecksumAddress(transaction.To.Address),
			Timestamp:       hexutil.MustDecodeUint64(transaction.Block.Timestamp),
			To:              address.Hex(),
			TxHash:          transaction.Hash,
			TxIndex:         transaction.Index,
			TXType:          "register",
		}

		if transaction.Status == 1 {
			addEvent.Success = true
		}

		if err := f.pub.Publish(
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
