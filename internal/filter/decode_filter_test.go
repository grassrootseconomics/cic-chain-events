package filter

import (
	"context"
	"testing"

	"github.com/grassrootseconomics/cic-chain-events/pkg/fetch"
	"github.com/stretchr/testify/suite"
	"github.com/zerodha/logf"
)

type DecodeFilterSuite struct {
	suite.Suite
	filter Filter
}

func (s *DecodeFilterSuite) SetupSuite() {
	logg := logf.New(
		logf.Opts{
			Level: logf.DebugLevel,
		},
	)

	s.filter = NewDecodeFilter(DecodeFilterOpts{
		Logg: logg,
	})
}

func (s *DecodeFilterSuite) TestTranfserInputs() {
	type testCase struct {
		transactionData fetch.Transaction
		want            bool
		wantErr         bool
	}

	// Generated with eth-encode
	tests := []testCase{
		{
			transactionData: fetch.Transaction{
				InputData: "0xa9059cbb000000000000000000000000000000000000000000000000000000000000dEaD00000000000000000000000000000000000000000000000000000000000003e8",
			},
			want:    true,
			wantErr: false,
		},
		{
			transactionData: fetch.Transaction{
				InputData: "0x23b872dd000000000000000000000000000000000000000000000000000000000000dEaD000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000003e8",
			},
			want:    true,
			wantErr: false,
		},
		{
			transactionData: fetch.Transaction{
				InputData: "0x449a52f8000000000000000000000000000000000000000000000000000000000000dEaD00000000000000000000000000000000000000000000000000000000000003e8",
			},
			want:    true,
			wantErr: false,
		},
		{
			transactionData: fetch.Transaction{
				InputData: "0x8d72ec9d000000000000000000000000000000000000000000000000000000000000dEaD00000000000000000000000000000000000000000000000000000000000003e8",
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

func TestDecodeFilterSuite(t *testing.T) {
	suite.Run(t, new(DecodeFilterSuite))
}
