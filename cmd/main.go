package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/VictoriaMetrics/metrics"
	"github.com/grassrootseconomics/cic-chain-events/internal/exporter"
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
			initNoopFilter(),
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

	if ko.Bool("metrics.expose") {
		metricsServer := &http.Server{
			Addr: ":9090",
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			exporter.Register(syncerStats)

			http.HandleFunc("/metrics", func(w http.ResponseWriter, _ *http.Request) {
				metrics.WritePrometheus(w, true)
			})

			if err := metricsServer.ListenAndServe(); err != nil {
				lo.Fatal("metrics server error", "error", err)
			}
		}()
	}

	<-ctx.Done()
	lo.Info("graceful shutdown triggered")

	workerPool.Stop()
}
