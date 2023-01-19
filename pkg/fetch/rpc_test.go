package fetch

import (
	"context"
	"os"
	"testing"

	celo "github.com/grassrootseconomics/cic-celo-sdk"
	"github.com/stretchr/testify/suite"
)

var (
	rpcEndpoint = os.Getenv("TEST_RPC_ENDPOINT")
)

type RPCTestSuite struct {
	suite.Suite
	fetch Fetch
}

func (s *RPCTestSuite) SetupSuite() {
	celoProvider, err := celo.NewProvider(celo.ProviderOpts{
		ChainId:     celo.MainnetChainId,
		RpcEndpoint: rpcEndpoint,
	})

	if err != nil {
		return
	}

	s.fetch = NewRPCFetcher(RPCOpts{
		RPCProvider: celoProvider,
	})
}

func (s *RPCTestSuite) Test_E2E_Fetch_Existing_Block() {
	resp, err := s.fetch.Block(context.Background(), 14974600)
	s.NoError(err)
	s.Len(resp.Data.Block.Transactions, 3)
}

func (s *RPCTestSuite) Test_E2E_Fetch_Non_Existing_Block() {
	_, err := s.fetch.Block(context.Background(), 14974600000)
	s.Error(err)
}

func TestRPCSuite(t *testing.T) {
	suite.Run(t, new(RPCTestSuite))
}
