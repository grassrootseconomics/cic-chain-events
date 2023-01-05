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

func (hs *HeadSyncer) Start(ctx context.Context) error {
	headerReceiver := make(chan *types.Header, 10)

	sub, err := hs.ethClient.SubscribeNewHead(ctx, headerReceiver)
	if err != nil {
		return err
	}

	for {
		select {
		case header := <-headerReceiver:
			block := header.Number.Uint64()
			hs.logg.Debug("head syncer received new block", "block", block)

			hs.stats.UpdateHeadCursor(block)
			hs.pool.Submit(func() {
				if err := hs.pipeline.Run(block); err != nil {
					hs.logg.Error("pipeline run error", "error", err)
				}
			})
		case err := <-sub.Err():
			hs.logg.Error("head syncer error", "error", err)
			return err
		case <-ctx.Done():
			hs.logg.Debug("head syncer shutdown signnal received")
			sub.Unsubscribe()
			return nil
		}
	}
}
