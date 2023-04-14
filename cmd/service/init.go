package main

import (
	"context"
	"strings"
	"time"

	"github.com/alitto/pond"
	"github.com/grassrootseconomics/celoutils"
	"github.com/grassrootseconomics/cic-chain-events/internal/pool"
	"github.com/grassrootseconomics/cic-chain-events/internal/pub"
	"github.com/grassrootseconomics/cic-chain-events/internal/store"
	"github.com/grassrootseconomics/cic-chain-events/pkg/fetch"
	"github.com/jackc/pgx/v5"
	"github.com/knadh/goyesql/v2"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/nats-io/nats.go"
	"github.com/zerodha/logf"
)

func initLogger() logf.Logger {
	loggOpts := logf.Opts{}

	if debugFlag {
		loggOpts.EnableColor = true
		loggOpts.EnableColor = true
		loggOpts.Level = logf.DebugLevel
	}

	return logf.New(loggOpts)
}

func initConfig() *koanf.Koanf {
	var (
		ko = koanf.New(".")
	)

	confFile := file.Provider(confFlag)
	if err := ko.Load(confFile, toml.Parser()); err != nil {
		lo.Fatal("init: could not load config file", "error", err)
	}

	if err := ko.Load(env.Provider("EVENTS_", ".", func(s string) string {
		return strings.ReplaceAll(strings.ToLower(
			strings.TrimPrefix(s, "EVENTS_")), "__", ".")
	}), nil); err != nil {
		lo.Fatal("init: could not override config from env vars", "error", err)
	}

	if debugFlag {
		ko.Print()
	}

	return ko
}

func initQueries(queriesPath string) goyesql.Queries {
	queries, err := goyesql.ParseFile(queriesPath)
	if err != nil {
		lo.Fatal("init: could not load queries file", "error", err)
	}

	return queries
}

func initPgStore(migrationsPath string, queries goyesql.Queries) store.Store[pgx.Rows] {
	pgStore, err := store.NewPostgresStore(store.PostgresStoreOpts{
		MigrationsFolderPath: migrationsPath,
		DSN:                  ko.MustString("postgres.dsn"),
		InitialLowerBound:    uint64(ko.MustInt64("syncer.initial_lower_bound")),
		Logg:                 lo,
		Queries:              queries,
	})
	if err != nil {
		lo.Fatal("init: critical error loading chain provider", "error", err)
	}

	return pgStore
}

func initFetcher() fetch.Fetch {
	return fetch.NewGraphqlFetcher(fetch.GraphqlOpts{
		GraphqlEndpoint: ko.MustString("chain.graphql_endpoint"),
	})
}

func initJanitorWorkerPool(ctx context.Context) *pond.WorkerPool {
	return pool.NewPool(ctx, pool.Opts{
		Concurrency: ko.MustInt("syncer.janitor_concurrency"),
		QueueSize:   ko.MustInt("syncer.janitor_queue_size"),
	})
}

func initHeadSyncerWorkerPool(ctx context.Context) *pond.WorkerPool {
	return pool.NewPool(ctx, pool.Opts{
		Concurrency: 1,
		QueueSize:   1,
	})
}

func initJetStream() (*nats.Conn, nats.JetStreamContext) {
	natsConn, err := nats.Connect(ko.MustString("jetstream.endpoint"))
	if err != nil {
		lo.Fatal("init: critical error connecting to NATS", "error", err)
	}

	js, err := natsConn.JetStream()
	if err != nil {
		lo.Fatal("init: bad JetStream opts", "error", err)

	}

	return natsConn, js
}

func initPub(natsConn *nats.Conn, jsCtx nats.JetStreamContext) *pub.Pub {
	pub, err := pub.NewPub(pub.PubOpts{
		DedupDuration:   time.Duration(ko.MustInt("jetstream.dedup_duration_hrs")) * time.Hour,
		JsCtx:           jsCtx,
		NatsConn:        natsConn,
		PersistDuration: time.Duration(ko.MustInt("jetstream.persist_duration_hrs")) * time.Hour,
	})
	if err != nil {
		lo.Fatal("init: critical error bootstrapping pub", "error", err)
	}

	return pub
}

func initCeloProvider() *celoutils.Provider {
	providerOpts := celoutils.ProviderOpts{
		RpcEndpoint: ko.MustString("chain.rpc_endpoint"),
	}

	if ko.Bool("chain.testnet") {
		providerOpts.ChainId = celoutils.TestnetChainId
	} else {
		providerOpts.ChainId = celoutils.MainnetChainId
	}

	provider, err := celoutils.NewProvider(providerOpts)
	if err != nil {
		lo.Fatal("init: critical error loading chain provider", "error", err)
	}

	return provider
}
