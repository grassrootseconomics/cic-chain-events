package syncer

import (
	"context"
	"time"

	"github.com/alitto/pond"
	"github.com/celo-org/celo-blockchain/core/types"
	"github.com/celo-org/celo-blockchain/ethclient"
	"github.com/grassrootseconomics/cic-chain-events/internal/pipeline"
	"github.com/zerodha/logf"
)

const (
	jobTimeout = 5 * time.Second
)

type (
	HeadSyncerOpts struct {
		Logg       logf.Logger
		Pipeline   *pipeline.Pipeline
		Pool       *pond.WorkerPool
		Stats      *Stats
		WsEndpoint string
	}

	HeadSyncer struct {
		ethClient *ethclient.Client
		logg      logf.Logger
		pipeline  *pipeline.Pipeline
		pool      *pond.WorkerPool
		stats     *Stats
	}
)

func NewHeadSyncer(o HeadSyncerOpts) (*HeadSyncer, error) {
	ethClient, err := ethclient.Dial(o.WsEndpoint)
	if err != nil {
		return nil, err
	}

	return &HeadSyncer{
		ethClient: ethClient,
		logg:      o.Logg,
		pipeline:  o.Pipeline,
		pool:      o.Pool,
		stats:     o.Stats,
	}, nil
}

// Start creates a websocket subscription and actively receives new blocks until stopped
// or a critical error occurs.
func (hs *HeadSyncer) Start(ctx context.Context) error {
	headerReceiver := make(chan *types.Header, 1)

	sub, err := hs.ethClient.SubscribeNewHead(context.Background(), headerReceiver)
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()

	for {
		select {
		case <-ctx.Done():
			hs.logg.Info("head syncer: shutdown signal received")
			return nil
		case err := <-sub.Err():
			return err
		case header := <-headerReceiver:
			blockNumber := header.Number.Uint64()
			hs.logg.Debug("head syncer: received new block", "block", blockNumber)
			hs.stats.UpdateHeadCursor(blockNumber)
			hs.pool.Submit(func() {
				ctx, cancel := context.WithTimeout(context.Background(), jobTimeout)
				defer cancel()

				if err := hs.pipeline.Run(ctx, blockNumber); err != nil {
					hs.logg.Error("head syncer: piepline run error", "error", err)
				}
			})
		}
	}
}
