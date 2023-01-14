package filter

import (
	"context"
	"math/big"

	"github.com/celo-org/celo-blockchain/common"
	"github.com/grassrootseconomics/cic-chain-events/pkg/fetch"
	"github.com/grassrootseconomics/w3-celo-patch"
	"github.com/zerodha/logf"
)

var (
	transferSig     = w3.MustNewFunc("transfer(address, uint256)", "bool")
	transferFromSig = w3.MustNewFunc("transferFrom(address, address, uint256)", "bool")
	mintToSig       = w3.MustNewFunc("mintTo(address, uint256)", "bool")
)

type DecodeFilterOpts struct {
	Logg logf.Logger
}

type DecodeFilter struct {
	logg logf.Logger
}

func NewDecodeFilter(o DecodeFilterOpts) Filter {
	return &DecodeFilter{
		logg: o.Logg,
	}
}

func (f *DecodeFilter) Execute(_ context.Context, transaction fetch.Transaction) (bool, error) {
	switch transaction.InputData[:10] {
	case "0xa9059cbb":
		var (
			to    common.Address
			value big.Int
		)

		if err := transferSig.DecodeArgs(w3.B(transaction.InputData), &to, &value); err != nil {
			return false, err
		}

		f.logg.Debug("transfer", "to", to.Hex(), "value", value.String())
		return true, nil
	case "0x23b872dd":
		var (
			from  common.Address
			to    common.Address
			value big.Int
		)

		if err := transferFromSig.DecodeArgs(w3.B(transaction.InputData), &from, &to, &value); err != nil {
			return false, err
		}

		f.logg.Debug("transferFrom", "from", from.Hex(), "to", to.Hex(), "value", value.String())
		return true, nil
	case "0x449a52f8":
		var (
			to    common.Address
			value big.Int
		)

		if err := mintToSig.DecodeArgs(w3.B(transaction.InputData), &to, &value); err != nil {
			return false, err
		}

		f.logg.Debug("mintTo", "to", to.Hex(), "value", value.String())
		return true, nil
	default:
		f.logg.Debug("unknownSignature", "inpuData", transaction.InputData)
		return false, nil
	}
}
