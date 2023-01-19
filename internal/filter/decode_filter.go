package filter

import (
	"context"
	"encoding/json"
	"math/big"

	"github.com/celo-org/celo-blockchain/common"
	"github.com/grassrootseconomics/cic-chain-events/pkg/fetch"
	"github.com/grassrootseconomics/w3-celo-patch"
	"github.com/nats-io/nats.go"
	"github.com/zerodha/logf"
)

var (
	transferSig     = w3.MustNewFunc("transfer(address, uint256)", "bool")
	transferFromSig = w3.MustNewFunc("transferFrom(address, address, uint256)", "bool")
	mintToSig       = w3.MustNewFunc("mintTo(address, uint256)", "bool")
)

type DecodeFilterOpts struct {
	Logg  logf.Logger
	JSCtx nats.JetStreamContext
}

type DecodeFilter struct {
	logg logf.Logger
	js   nats.JetStreamContext
}

type minimalTxInfo struct {
	Block        uint64 `json:"block"`
	From         string `json:"from"`
	Success      bool   `json:"success"`
	To           string `json:"to"`
	TokenAddress string `json:"tokenAddress"`
	TxHash       string `json:"transactionHash"`
	TxIndex      uint   `json:"transactionIndex"`
	Value        uint64 `json:"value"`
}

func NewDecodeFilter(o DecodeFilterOpts) Filter {
	return &DecodeFilter{
		logg: o.Logg,
		js:   o.JSCtx,
	}
}

func (f *DecodeFilter) Execute(_ context.Context, transaction *fetch.Transaction) (bool, error) {
	switch transaction.InputData[:10] {
	case "0xa9059cbb":
		var (
			to    common.Address
			value big.Int
		)

		if err := transferSig.DecodeArgs(w3.B(transaction.InputData), &to, &value); err != nil {
			return false, err
		}

		transferEvent := &minimalTxInfo{
			Block:        transaction.Block.Number,
			From:         transaction.From.Address,
			To:           to.Hex(),
			TokenAddress: transaction.To.Address,
			TxHash:       transaction.Hash,
			TxIndex:      transaction.Index,
			Value:        value.Uint64(),
		}

		if transaction.Status == 1 {
			transferEvent.Success = true
		}

		json, err := json.Marshal(transferEvent)
		if err != nil {
			return false, err
		}

		_, err = f.js.Publish("CHAIN.transfer", json, nats.MsgId(transaction.Hash))
		if err != nil {
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

		transferFromEvent := &minimalTxInfo{
			Block:        transaction.Block.Number,
			From:         from.Hex(),
			To:           to.Hex(),
			TokenAddress: transaction.To.Address,
			TxHash:       transaction.Hash,
			TxIndex:      transaction.Index,
			Value:        value.Uint64(),
		}

		if transaction.Status == 1 {
			transferFromEvent.Success = true
		}

		json, err := json.Marshal(transferFromEvent)
		if err != nil {
			return false, err
		}

		_, err = f.js.Publish("CHAIN.transferFrom", json, nats.MsgId(transaction.Hash))
		if err != nil {
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

		mintToEvent := &minimalTxInfo{
			Block:        transaction.Block.Number,
			From:         transaction.From.Address,
			To:           to.Hex(),
			TokenAddress: transaction.To.Address,
			TxHash:       transaction.Hash,
			TxIndex:      transaction.Index,
			Value:        value.Uint64(),
		}

		if transaction.Status == 1 {
			mintToEvent.Success = true
		}

		json, err := json.Marshal(mintToEvent)
		if err != nil {
			return false, err
		}

		_, err = f.js.Publish("CHAIN.mintTo", json, nats.MsgId(transaction.Hash))
		if err != nil {
			return false, err
		}

		return true, nil
	default:
		f.logg.Debug("unknownSignature", "inpuData", transaction.InputData)
		return false, nil
	}
}
