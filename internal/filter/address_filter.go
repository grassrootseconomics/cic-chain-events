package filter

import (
	"context"

	"github.com/inethi/inethi-cic-chain-events/pkg/fetch"
	"github.com/zerodha/logf"
)

const (
	KRNVoucherAddress = "0x8bab657c88eb3c724486d113e650d2c659aa23d2"
	SRFVoucherAddress = "0x45d747172e77d55575c197cba9451bc2cd8f4958"
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
	if transaction.To.Address == KRNVoucherAddress || transaction.To.Address == SRFVoucherAddress {
		return true, nil
	}

	return false, nil
}
