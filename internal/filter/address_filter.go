package filter

import (
	"context"

	"github.com/inethi/inethi-cic-chain-events/pkg/fetch"
	"github.com/zerodha/logf"
)

const (
	KRNVoucherAddress = "0x8bab657c88eb3c724486d113e650d2c659aa23d2"
)

type (
	AddressFilterOpts struct {
		Logg logf.Logger
	}

	AddressFilter struct {
		logg logf.Logger
	}
)

func NewAddressFilter(o AddressFilterOpts) Filter {
	return &AddressFilter{
		logg: o.Logg,
	}
}

func (f *AddressFilter) Execute(_ context.Context, transaction *fetch.Transaction) (bool, error) {
	if transaction.To.Address == KRNVoucherAddress {
		return true, nil
	}

	return false, nil
}
