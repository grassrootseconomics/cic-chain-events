package filter

import (
	"context"
	"sync"
	"testing"

	"github.com/grassrootseconomics/cic-chain-events/pkg/fetch"
	"github.com/stretchr/testify/suite"
	"github.com/zerodha/logf"
)

type AddressFilterSuite struct {
	suite.Suite
	filter Filter
}

func (s *AddressFilterSuite) SetupSuite() {
	addressCache := &sync.Map{}

	addressCache.Store("0x6914ba1c49d3c3f32a9e65a0661d7656cb292e9f", "")

	logg := logf.New(
		logf.Opts{
			Level: logf.DebugLevel,
		},
	)

	s.filter = NewAddressFilter(AddressFilterOpts{
		Cache: addressCache,
		Logg:  logg,
	})
}

func (s *AddressFilterSuite) TestAddresses() {
	type testCase struct {
		transactionData fetch.Transaction
		want            bool
		wantErr         bool
	}

	// Generated with eth-encode
	tests := []testCase{
		{
			transactionData: fetch.Transaction{
				To: struct {
					Address string "json:\"address\""
				}{
					Address: "0x6914ba1c49d3c3f32a9e65a0661d7656cb292e9f",
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			transactionData: fetch.Transaction{
				To: struct {
					Address string "json:\"address\""
				}{
					Address: "0x6914ba1c49d3c3f32a9e65a0661d7656cb292e9x",
				},
			},
			want:    false,
			wantErr: false,
		},
	}

	for _, test := range tests {
		next, err := s.filter.Execute(context.Background(), test.transactionData)
		s.NoError(err)
		s.Equal(test.want, next)
	}
}

func TestAddressFilterSuite(t *testing.T) {
	suite.Run(t, new(AddressFilterSuite))
}
