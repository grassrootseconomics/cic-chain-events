package fetch

import (
	"context"
	"math/big"
	"strings"

	"github.com/celo-org/celo-blockchain/common/hexutil"
	"github.com/celo-org/celo-blockchain/core/types"
	celo "github.com/grassrootseconomics/cic-celo-sdk"
	"github.com/grassrootseconomics/w3-celo-patch/module/eth"
	"github.com/grassrootseconomics/w3-celo-patch/w3types"
)

// RPCOpts reprsents the required paramters for an RPC fetcher.
type RPCOpts struct {
	RPCProvider *celo.Provider
}

// RPC is a RPC based block and transaction fetcher.
type RPC struct {
	provider *celo.Provider
}

// NewRPCFetcher returns a new RPC fetcher which implemnts Fetch.
// Note: No rate limiting feeature.
func NewRPCFetcher(o RPCOpts) Fetch {
	return &RPC{
		provider: o.RPCProvider,
	}
}

// Block fetches via RPC and transforms the response to adapt to the GraphQL JSON response struct.
func (f *RPC) Block(ctx context.Context, blockNumber uint64) (FetchResponse, error) {
	var (
		block         types.Block
		fetchResponse FetchResponse
	)

	if err := f.provider.Client.CallCtx(
		ctx,
		eth.BlockByNumber(big.NewInt(int64(blockNumber))).Returns(&block),
	); err != nil {
		return fetchResponse, err
	}

	txCount := len(block.Transactions())
	batchCalls := make([]w3types.Caller, txCount*2)

	txs := make([]types.Transaction, txCount)
	txsReceipt := make([]types.Receipt, txCount)

	// Prepare batch calls.
	for i, tx := range block.Transactions() {
		batchCalls[i] = eth.Tx(tx.Hash()).Returns(&txs[i])
		batchCalls[txCount+i] = eth.TxReceipt(tx.Hash()).Returns(&txsReceipt[i])
	}

	if err := f.provider.Client.CallCtx(
		ctx,
		batchCalls...,
	); err != nil {
		return fetchResponse, err
	}

	// Transform response and adapt to FetchResponse.
	for i := 0; i < txCount; i++ {
		var txObject Transaction

		txObject.Block.Number = block.NumberU64()
		txObject.Block.Timestamp = hexutil.EncodeUint64(block.Time())

		from, err := types.Sender(types.LatestSignerForChainID(txs[i].ChainId()), &txs[i])
		if err != nil {
			return fetchResponse, err
		}
		txObject.From.Address = strings.ToLower(from.Hex())
		// This check ignores contract deployment transactions.
		if txs[i].To() != nil {
			txObject.To.Address = strings.ToLower(txs[i].To().Hex())
		}
		txObject.Value = hexutil.EncodeBig(txs[i].Value())
		txObject.InputData = hexutil.Encode(txs[i].Data())

		txObject.Hash = txsReceipt[i].TxHash.Hex()
		txObject.Index = txsReceipt[i].TransactionIndex
		txObject.Status = txsReceipt[i].Status
		txObject.GasUsed = txsReceipt[i].GasUsed

		fetchResponse.Data.Block.Transactions = append(
			fetchResponse.Data.Block.Transactions,
			txObject,
		)
	}

	return fetchResponse, nil
}
