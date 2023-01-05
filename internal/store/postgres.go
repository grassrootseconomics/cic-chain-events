package store

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/knadh/goyesql/v2"
	"github.com/zerodha/logf"
)

type queries struct {
	CommitBlock         string `query:"commit-block"`
	GetMissingBlocks    string `query:"get-missing-blocks"`
	GetSearchBounds     string `query:"get-search-bounds"`
	InitSyncerMeta      string `query:"init-syncer-meta"`
	SetSearchLowerBound string `query:"set-search-lower-bound"`
}

type PostgresStoreOpts struct {
	DSN               string
	InitialLowerBound uint64
	Logg              logf.Logger
	Queries           goyesql.Queries
}

type PostgresStore struct {
	logg    logf.Logger
	pool    *pgxpool.Pool
	queries queries
}

func NewPostgresStore(o PostgresStoreOpts) (Store[pgx.Rows], error) {
	postgresStore := &PostgresStore{
		logg: o.Logg,
	}

	if err := goyesql.ScanToStruct(&postgresStore.queries, o.Queries, nil); err != nil {
		return nil, fmt.Errorf("failed to scan queries %v", err)
	}

	parsedConfig, err := pgxpool.ParseConfig(o.DSN)
	if err != nil {
		return nil, err
	}

	dbPool, err := pgxpool.NewWithConfig(context.Background(), parsedConfig)
	if err != nil {
		return nil, err
	}

	_, err = dbPool.Exec(context.Background(), postgresStore.queries.InitSyncerMeta, o.InitialLowerBound)
	if err != nil {
		return nil, err
	}

	postgresStore.pool = dbPool

	return postgresStore, nil
}

func (s *PostgresStore) GetSearchBounds(batchSize uint64, headCursor uint64, headBlockLag uint64) (uint64, uint64, error) {
	var (
		lowerBound uint64
		upperBound uint64
	)

	if err := s.pool.QueryRow(
		context.Background(),
		s.queries.GetSearchBounds,
		batchSize,
		headCursor,
		headBlockLag,
	).Scan(&lowerBound, &upperBound); err != nil {
		s.logg.Error("pgx error", "error", err)
		return 0, 0, err
	}

	return lowerBound, upperBound, nil
}

func (s *PostgresStore) GetMissingBlocks(lowerBound uint64, upperBound uint64) (pgx.Rows, error) {
	rows, err := s.pool.Query(context.Background(), s.queries.GetMissingBlocks, lowerBound, upperBound)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (s *PostgresStore) SetSearchLowerBound(newLowerBound uint64) error {
	_, err := s.pool.Exec(context.Background(), s.queries.SetSearchLowerBound, newLowerBound)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) CommitBlock(block uint64) error {
	_, err := s.pool.Exec(context.Background(), s.queries.CommitBlock, block)
	if err != nil {
		return err
	}

	return nil
}
