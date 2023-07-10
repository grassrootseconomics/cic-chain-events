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

type (
	ApproveFilterOpts struct {
		Logg logf.Logger
		Pub  *pub.Pub
	}

	ApproveFilter struct {
		logg logf.Logger
		pub  *pub.Pub
	}
)

const (
	approveEventSubject = "CHAIN.approve"
)

var (
	approveSig = w3.MustNewFunc("approve(address, uint256)", "bool")
)

func NewApproveFilter(o ApproveFilterOpts) Filter {
	return &ApproveFilter{
		logg: o.Logg,
		pub:  o.Pub,
	}
}

func (f *ApproveFilter) Execute(_ context.Context, transaction *fetch.Transaction) (bool, error) {
	if len(transaction.InputData) < 10 {
		return true, nil
	}

	if transaction.InputData[:10] == "0x095ea7b3" {
		var (
			address common.Address
			value   big.Int
		)

		if err := approveSig.DecodeArgs(w3.B(transaction.InputData), &address, &value); err != nil {
			return false, err
		}

		approveEvent := &pub.MinimalTxInfo{
			Block:           transaction.Block.Number,
			ContractAddress: celoutils.ChecksumAddress(transaction.To.Address),
			Timestamp:       hexutil.MustDecodeUint64(transaction.Block.Timestamp),
			From:            celoutils.ChecksumAddress(transaction.From.Address),
			To:              address.Hex(),
			TxHash:          transaction.Hash,
			TxIndex:         transaction.Index,
			TXType:          "approve",
		}

		if transaction.Status == 1 {
			approveEvent.Success = true
		}

		if err := f.pub.Publish(
			approveEventSubject,
			transaction.Hash,
			approveEvent,
		); err != nil {
			return false, err
		}

		return true, nil
	}

	return true, nil
}
