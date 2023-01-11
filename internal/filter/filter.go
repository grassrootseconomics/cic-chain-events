package filter

import (
	"context"

	"github.com/grassrootseconomics/cic-chain-events/internal/fetch"
)

// Filter defines a read only filter which must return next as true/false or an error
type Filter interface {
	Execute(ctx context.Context, inputTransaction fetch.Transaction) (next bool, err error)
}
