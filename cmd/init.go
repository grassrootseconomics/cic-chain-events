package main

import (
	"strings"

	"github.com/alitto/pond"
	"github.com/grassrootseconomics/cic-chain-events/internal/fetch"
	"github.com/grassrootseconomics/cic-chain-events/internal/pool"
	"github.com/grassrootseconomics/cic-chain-events/internal/store"
	"github.com/jackc/pgx/v5"
	"github.com/knadh/goyesql/v2"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/zerodha/logf"
)

func initLogger(debug bool) logf.Logger {
	loggOpts := logf.Opts{
		EnableColor: true,
	}

	if debug {
		loggOpts.Level = logf.DebugLevel
	}

	return logf.New(loggOpts)
}

func initConfig(configFilePath string) *koanf.Koanf {
	var (
		ko = koanf.New(".")
	)

	confFile := file.Provider(configFilePath)
	if err := ko.Load(confFile, toml.Parser()); err != nil {
		lo.Fatal("could not load config file", "error", err)
	}

	if err := ko.Load(env.Provider("", ".", func(s string) string {
		return strings.ReplaceAll(strings.ToLower(
			strings.TrimPrefix(s, "")), "_", ".")
	}), nil); err != nil {
		lo.Fatal("could not override config from env vars", "error", err)
	}

	return ko
}

func initQueries(queriesPath string) goyesql.Queries {
	queries, err := goyesql.ParseFile(queriesPath)
	if err != nil {
		lo.Fatal("could not load queries file", "error", err)
	}

	return queries
}

func initPgStore() (store.Store[pgx.Rows], error) {
	pgStore, err := store.NewPostgresStore(store.PostgresStoreOpts{
		DSN:               ko.MustString("postgres.dsn"),
		InitialLowerBound: uint64(ko.MustInt64("indexer.initial_lower_bound")),
		Logg:              lo,
		Queries:           q,
	})
	if err != nil {
		return nil, err
	}

	return pgStore, nil
}

func initWorkerPool() *pond.WorkerPool {
	return pool.NewPool(pool.Opts{
		ConcurrencyFactor: ko.MustInt("indexer.concurrency"),
		PoolQueueSize:     ko.MustInt("indexer.queue_size"),
	})
}

func initFetcher() fetch.Fetch {
	return fetch.NewGraphqlFetcher(fetch.GraphqlOpts{
		GraphqlEndpoint: ko.MustString("chain.graphql_endpoint"),
	})
}
