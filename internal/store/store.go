package store

// Store defines indexer get and set queries.
// GetMissingBlocks returns a generic iterable.
type Store[T any] interface {
	GetSearchBounds(batchSize uint64, headCursor uint64, headBlockLag uint64) (lowerBound uint64, upperBound uint64, err error)
	GetMissingBlocks(lowerBound uint64, upperBound uint64) (missingBlocksIterable T, err error)
	SetSearchLowerBound(newLowerBound uint64) (err error)
	CommitBlock(block uint64) (err error)
}
