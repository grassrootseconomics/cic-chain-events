package store

import "context"

// Store defines all relevant get/set queries against the implemented storage backend.
type Store[T any] interface {
	GetSearchBounds(ctx context.Context, batchSize uint64, headCursor uint64, headBlockLag uint64) (lowerBound uint64, upperBound uint64, err error)
	GetMissingBlocks(ctx context.Context, lowerBound uint64, upperBound uint64) (missingBlocksIterable T, err error)
	SetSearchLowerBound(ctx context.Context, newLowerBound uint64) (err error)
	CommitBlock(ctx context.Context, block uint64) (err error)
}
