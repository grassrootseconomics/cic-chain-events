package filter

import (
	"context"

	"github.com/inethi/inethi-cic-chain-events/pkg/fetch"
)

// Filter defines a read only filter which must return next as true/false or an error
type Filter interface {
	Execute(ctx context.Context, inputTransaction *fetch.Transaction) (next bool, err error)
}
