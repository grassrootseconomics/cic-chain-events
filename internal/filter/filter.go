package filter

import "github.com/grassrootseconomics/cic-chain-events/internal/fetch"

// Filter defines a read only filter which must return next as true/false or an error
type Filter interface {
	Execute(inputTransaction fetch.Transaction) (next bool, err error)
}
