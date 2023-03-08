package store

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/tern/v2/migrate"
	"github.com/knadh/goyesql/v2"
	"github.com/zerodha/logf"
)

const (
	schemaTable = "schema_version"
)

type (
	queries struct {
		CommitBlock         string `query:"commit-block"`
		GetMissingBlocks    string `query:"get-missing-blocks"`
		GetSearchBounds     string `query:"get-search-bounds"`
		InitSyncerMeta      string `query:"init-syncer-meta"`
		SetSearchLowerBound string `query:"set-search-lower-bound"`
	}

	PostgresStoreOpts struct {
		DSN                  string
		MigrationsFolderPath string
		InitialLowerBound    uint64
		Logg                 logf.Logger
		Queries              goyesql.Queries
	}
)

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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbPool, err := pgxpool.NewWithConfig(ctx, parsedConfig)
	if err != nil {
		return nil, err
	}

	conn, err := dbPool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	migrator, err := migrate.NewMigrator(ctx, conn.Conn(), schemaTable)
	if err != nil {
		return nil, err
	}

	if err := migrator.LoadMigrations(os.DirFS(o.MigrationsFolderPath)); err != nil {
		return nil, err
	}

	if err := migrator.Migrate(ctx); err != nil {
		return nil, err
	}

	_, err = dbPool.Exec(ctx, postgresStore.queries.InitSyncerMeta, o.InitialLowerBound)
	if err != nil {
		return nil, err
	}

	postgresStore.pool = dbPool

	return postgresStore, nil
}

func (s *PostgresStore) GetSearchBounds(ctx context.Context, batchSize uint64, headCursor uint64, headBlockLag uint64) (uint64, uint64, error) {
	var (
		lowerBound uint64
		upperBound uint64
	)

	if err := s.pool.QueryRow(
		ctx,
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

func (s *PostgresStore) GetMissingBlocks(ctx context.Context, lowerBound uint64, upperBound uint64) (pgx.Rows, error) {
	rows, err := s.pool.Query(ctx, s.queries.GetMissingBlocks, lowerBound, upperBound)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (s *PostgresStore) SetSearchLowerBound(ctx context.Context, newLowerBound uint64) error {
	_, err := s.pool.Exec(ctx, s.queries.SetSearchLowerBound, newLowerBound)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) CommitBlock(ctx context.Context, block uint64) error {
	_, err := s.pool.Exec(ctx, s.queries.CommitBlock, block)
	if err != nil {
		return err
	}

	return nil
}
