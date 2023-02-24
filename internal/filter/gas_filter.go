package filter

import (
	"context"
	"encoding/json"

	"github.com/celo-org/celo-blockchain/common"
	"github.com/grassrootseconomics/cic-chain-events/pkg/fetch"
	"github.com/grassrootseconomics/w3-celo-patch"
	"github.com/nats-io/nats.go"
	"github.com/zerodha/logf"
)

var (
	giveToSig = w3.MustNewFunc("giveTo(address)", "uint256")
)

type GasFilterOpts struct {
	Logg  logf.Logger
	JSCtx nats.JetStreamContext
}

type GasFilter struct {
	logg logf.Logger
	js   nats.JetStreamContext
}

type minimalGasGiftTxInfo struct {
	Block   uint64 `json:"block"`
	Success bool   `json:"success"`
	To      string `json:"to"`
	TxHash  string `json:"transactionHash"`
	TxIndex uint   `json:"transactionIndex"`
}

func NewGasFilter(o GasFilterOpts) Filter {
	return &GasFilter{
		logg: o.Logg,
		js:   o.JSCtx,
	}
}

func (f *GasFilter) Execute(_ context.Context, transaction *fetch.Transaction) (bool, error) {
	switch transaction.InputData[:10] {
	case "0x63e4bff4":
		var (
			to common.Address
		)

		if err := giveToSig.DecodeArgs(w3.B(transaction.InputData), &to); err != nil {
			return false, err
		}

		transferEvent := &minimalGasGiftTxInfo{
			Block:   transaction.Block.Number,
			To:      to.Hex(),
			TxHash:  transaction.Hash,
			TxIndex: transaction.Index,
		}

		if transaction.Status == 1 {
			transferEvent.Success = true
		}

		json, err := json.Marshal(transferEvent)
		if err != nil {
			return false, err
		}

		_, err = f.js.Publish("CHAIN.gasGiveTo", json, nats.MsgId(transaction.Hash))
		if err != nil {
			return false, err
		}

		return true, nil
	default:
		return false, nil
	}
}
