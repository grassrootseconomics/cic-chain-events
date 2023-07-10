package main

import (
	"context"
	"flag"
	"strings"
	"sync"
	"time"

	"github.com/grassrootseconomics/cic-chain-events/internal/filter"
	"github.com/grassrootseconomics/cic-chain-events/internal/pipeline"
	"github.com/grassrootseconomics/cic-chain-events/internal/pub"
	"github.com/grassrootseconomics/cic-chain-events/internal/syncer"
	"github.com/knadh/koanf/v2"
	"github.com/labstack/echo/v4"
	"github.com/zerodha/logf"
)

type (
	internalServicesContainer struct {
		apiService *echo.Echo
		pub        *pub.Pub
	}
)

var (
	build string

	confFlag             string
	debugFlag            bool
	migrationsFolderFlag string
	queriesFlag          string

	ko *koanf.Koanf
	lo logf.Logger
)

func init() {
	flag.StringVar(&confFlag, "config", "config.toml", "Config file location")
	flag.BoolVar(&debugFlag, "debug", false, "Enable debug logging")
	flag.StringVar(&migrationsFolderFlag, "migrations", "migrations/", "Migrations folder location")
	flag.StringVar(&queriesFlag, "queries", "queries.sql", "Queries file location")
	flag.Parse()

	lo = initLogger()
	ko = initConfig()
}

func main() {
	lo.Info("main: starting cic-chain-events", "build", build)

	parsedQueries := initQueries(queriesFlag)
	graphqlFetcher := initFetcher()
	pgStore := initPgStore(migrationsFolderFlag, parsedQueries)
	natsConn, jsCtx := initJetStream()
	jsPub := initPub(natsConn, jsCtx)

	celoProvider := initCeloProvider()
	cache := &sync.Map{}

	pipeline := pipeline.NewPipeline(pipeline.PipelineOpts{
		BlockFetcher: graphqlFetcher,
		Filters: []filter.Filter{
			initAddressFilter(celoProvider, cache),
			initGasGiftFilter(jsPub),
			initTransferFilter(jsPub),
			initRegisterFilter(jsPub),
			initApproveFilter(jsPub),
			initTokenIndexFilter(cache, jsPub),
		},
		Logg:  lo,
		Store: pgStore,
	})

	internalServices := &internalServicesContainer{
		pub: jsPub,
	}
	syncerStats := &syncer.Stats{}
	wg := &sync.WaitGroup{}

	signalCh, closeCh := createSigChannel()
	defer closeCh()

	ctx, cancel := context.WithCancel(context.Background())

	headSyncer, err := syncer.NewHeadSyncer(syncer.HeadSyncerOpts{
		Logg:       lo,
		Pipeline:   pipeline,
		Pool:       initHeadSyncerWorkerPool(ctx),
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
		Pool:          initJanitorWorkerPool(ctx),
		Stats:         syncerStats,
		Store:         pgStore,
		SweepInterval: time.Second * time.Duration(ko.MustInt64("syncer.janitor_sweep_interval")),
	})

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := headSyncer.Start(ctx); err != nil {
			lo.Info("main: starting head syncer")
			lo.Fatal("main: critical error starting head syncer", "error", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		lo.Info("main: starting janitor")
		if err := janitor.Start(ctx); err != nil {
			lo.Fatal("main: critical error starting janitor", "error", err)
		}
	}()

	internalServices.apiService = initApiServer()
	wg.Add(1)
	go func() {
		defer wg.Done()
		host := ko.MustString("service.address")
		lo.Info("main: starting API server", "host", host)
		if err := internalServices.apiService.Start(host); err != nil {
			if strings.Contains(err.Error(), "Server closed") {
				lo.Info("main: shutting down server")
			} else {
				lo.Fatal("main: critical error shutting down server", "err", err)
			}
		}
	}()

	lo.Info("main: graceful shutdown triggered", "signal", <-signalCh)
	cancel()
	startGracefulShutdown(context.Background(), internalServices)

	wg.Wait()
}
