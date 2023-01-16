package fetch

import (
	"context"
	"math/big"

	"github.com/celo-org/celo-blockchain/core/types"
	celo "github.com/grassrootseconomics/cic-celo-sdk"
	"github.com/grassrootseconomics/w3-celo-patch/module/eth"
	"github.com/grassrootseconomics/w3-celo-patch/w3types"
)

type RPCOpts struct {
	RPCProvider *celo.Provider
}

type RPC struct {
	provider *celo.Provider
}

func NewRPCFetcher(o RPCOpts) Fetch {
	return &RPC{
		provider: o.RPCProvider,
	}
}

// This method makes 2 calls. 1 for the block and 1 batched call for txs + receipts.
// Should work on free tier RPC services.
func (f *RPC) Block(ctx context.Context, blockNumber uint64) (FetchResponse, error) {
	var (
		block types.Block

		fetchResponse FetchResponse
	)

	if err := f.provider.Client.CallCtx(
		ctx,
		eth.BlockByNumber(big.NewInt(int64(blockNumber))).Returns(&block),
	); err != nil {
		return FetchResponse{}, nil
	}

	txCount := len(block.Transactions())
	batchCalls := make([]w3types.Caller, txCount*2)

	txs := make([]types.Transaction, txCount)
	txsReceipt := make([]types.Receipt, txCount)

	for i, tx := range block.Transactions() {
		batchCalls[i] = eth.Tx(tx.Hash()).Returns(&txs[i])
		batchCalls[txCount+i] = eth.TxReceipt(tx.Hash()).Returns(&txsReceipt[i])
	}

	if err := f.provider.Client.CallCtx(
		ctx,
		batchCalls...,
	); err != nil {
		return FetchResponse{}, nil
	}

	// TODO: Create FetchResponse

	return fetchResponse, nil
}
