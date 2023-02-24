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

	"github.com/grassrootseconomics/cic-chain-events/internal/filter"
	"github.com/grassrootseconomics/cic-chain-events/internal/pipeline"
	"github.com/grassrootseconomics/cic-chain-events/internal/pool"
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

	janitorWorkerPool := pool.NewPool(ctx, pool.Opts{
		Concurrency: ko.MustInt("syncer.janitor_concurrency"),
		QueueSize:   ko.MustInt("syncer.janitor_queue_size"),
	})

	pgStore, err := initPgStore()
	if err != nil {
		lo.Fatal("main: critical error loading pg store", "error", err)
	}

	jsCtx, err := initJetStream()
	if err != nil {
		lo.Fatal("main: critical error loading jetstream context", "error", err)
	}

	graphqlFetcher := initFetcher()

	pipeline := pipeline.NewPipeline(pipeline.PipelineOpts{
		BlockFetcher: graphqlFetcher,
		Filters: []filter.Filter{
			initAddressFilter(),
			initTransferFilter(jsCtx),
			initGasGiftFilter(jsCtx),
		},
		Logg:  lo,
		Store: pgStore,
	})

	headSyncerWorker := pool.NewPool(ctx, pool.Opts{
		Concurrency: 1,
		QueueSize:   1,
	})

	headSyncer, err := syncer.NewHeadSyncer(syncer.HeadSyncerOpts{
		Logg:       lo,
		Pipeline:   pipeline,
		Pool:       headSyncerWorker,
		Stats:      syncerStats,
		WsEndpoint: ko.MustString("chain.ws_endpoint"),
	})
	if err != nil {
		lo.Fatal("main: crticial error loading head syncer", "error", err)
	}

	janitor := syncer.NewJanitor(syncer.JanitorOpts{
		BatchSize:     uint64(ko.MustInt64("syncer.janitor_queue_size")),
		Logg:          lo,
		Pipeline:      pipeline,
		Pool:          janitorWorkerPool,
		Stats:         syncerStats,
		Store:         pgStore,
		SweepInterval: time.Second * time.Duration(ko.MustInt64("syncer.janitor_sweep_interval")),
	})

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := headSyncer.Start(ctx); err != nil {
			lo.Fatal("main: critical error starting head syncer", "error", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := janitor.Start(ctx); err != nil {
			lo.Fatal("main: critical error starting janitor", "error", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		lo.Info("starting API server")
		if err := apiServer.Start(ko.MustString("api.address")); err != nil {
			if strings.Contains(err.Error(), "Server closed") {
				lo.Info("main: shutting down server")
			} else {
				lo.Fatal("main: critical error shutting down server", "err", err)
			}
		}
	}()

	<-ctx.Done()

	if err := apiServer.Shutdown(ctx); err != nil {
		lo.Error("main: could not gracefully shutdown api server", "err", err)
	}

	wg.Wait()
}
