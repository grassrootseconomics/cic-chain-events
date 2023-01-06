package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/grassrootseconomics/cic-chain-events/internal/api"
	"github.com/grassrootseconomics/cic-chain-events/internal/filter"
	"github.com/grassrootseconomics/cic-chain-events/internal/pipeline"
	"github.com/grassrootseconomics/cic-chain-events/internal/syncer"
	"github.com/knadh/goyesql/v2"
	"github.com/knadh/koanf"
	"github.com/zerodha/logf"
)

var (
	confFlag    string
	debugFlag   bool
	queriesFlag string

	ko *koanf.Koanf
	lo logf.Logger
	q  goyesql.Queries
)

func init() {
	flag.StringVar(&confFlag, "config", "config.toml", "Config file location")
	flag.BoolVar(&debugFlag, "log", true, "Enable debug logging")
	flag.StringVar(&queriesFlag, "queries", "queries.sql", "Queries file location")
	flag.Parse()

	lo = initLogger(debugFlag)
	ko = initConfig(confFlag)
	q = initQueries(queriesFlag)
}

func main() {
	syncerStats := &syncer.Stats{}
	wg := &sync.WaitGroup{}
	apiServer := initApiServer()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pgStore, err := initPgStore()
	if err != nil {
		lo.Fatal("error loading pg store", "error", err)
	}
	workerPool := initWorkerPool()
	graphqlFetcher := initFetcher()

	pipeline := pipeline.NewPipeline(pipeline.PipelineOpts{
		BlockFetcher: graphqlFetcher,
		Filters: []filter.Filter{
			initAddressFilter(),
			initTransferFilter(),
			// initNoopFilter(),
		},
		Logg:  lo,
		Store: pgStore,
	})

	headSyncer, err := syncer.NewHeadSyncer(syncer.HeadSyncerOpts{
		Logg:       lo,
		Pipeline:   pipeline,
		Pool:       workerPool,
		Stats:      syncerStats,
		WsEndpoint: ko.MustString("chain.ws_endpoint"),
	})
	if err != nil {
		lo.Fatal("error loading head syncer", "error", err)
	}

	janitor := syncer.NewJanitor(syncer.JanitorOpts{
		BatchSize:     uint64(ko.MustInt64("indexer.batch_size")),
		HeadBlockLag:  uint64(ko.MustInt64("indexer.head_block_lag")),
		Logg:          lo,
		Pipeline:      pipeline,
		Pool:          workerPool,
		Stats:         syncerStats,
		Store:         pgStore,
		SweepInterval: time.Second * time.Duration(ko.MustInt64("indexer.sweep_interval")),
	})

	apiServer.GET("/stats", api.StatsHandler(syncerStats, workerPool, lo))

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := headSyncer.Start(ctx); err != nil {
			lo.Fatal("head syncer error", "error", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := janitor.Start(ctx); err != nil {
			lo.Fatal("janitor error", "error", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		lo.Info("starting API server")
		if err := apiServer.Start(ko.MustString("api.address")); err != nil {
			if strings.Contains(err.Error(), "Server closed") {
				lo.Info("shutting down server")
			} else {
				lo.Fatal("could not start api server", "err", err)
			}
		}
	}()

	<-ctx.Done()
	lo.Info("graceful shutdown triggered")

	workerPool.Stop()
	if err := apiServer.Shutdown(ctx); err != nil {
		lo.Error("could not gracefully shutdown api server", "err", err)
	}

	wg.Wait()
}
