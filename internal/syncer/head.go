package syncer

import (
	"context"

	"github.com/alitto/pond"
	"github.com/celo-org/celo-blockchain/core/types"
	"github.com/celo-org/celo-blockchain/ethclient"
	"github.com/grassrootseconomics/cic-chain-events/internal/pipeline"
	"github.com/zerodha/logf"
)

type HeadSyncerOpts struct {
	Stats      *Stats
	Pipeline   *pipeline.Pipeline
	Logg       logf.Logger
	Pool       *pond.WorkerPool
	WsEndpoint string
}

type HeadSyncer struct {
	stats     *Stats
	pipeline  *pipeline.Pipeline
	logg      logf.Logger
	ethClient *ethclient.Client
	pool      *pond.WorkerPool
}

func NewHeadSyncer(o HeadSyncerOpts) (*HeadSyncer, error) {
	ethClient, err := ethclient.Dial(o.WsEndpoint)
	if err != nil {
		return nil, err
	}

	return &HeadSyncer{
		stats:     o.Stats,
		pipeline:  o.Pipeline,
		logg:      o.Logg,
		ethClient: ethClient,
		pool:      o.Pool,
	}, nil
}

// Start creates a websocket subscription and actively receives new blocks untill stopped
// or a critical error occurs.
func (hs *HeadSyncer) Start(ctx context.Context) error {
	headerReceiver := make(chan *types.Header, 1)

	sub, err := hs.ethClient.SubscribeNewHead(ctx, headerReceiver)
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
				if err := hs.pipeline.Run(context.Background(), blockNumber); err != nil {
					hs.logg.Error("head syncer: piepline run error", "error", err)
				}
			})
		}
	}
}
