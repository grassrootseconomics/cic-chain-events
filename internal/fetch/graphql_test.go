package fetch

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

var (
	graphqlEndpoint = "https://rpc.celo.grassecon.net/graphql"
)

type itGraphqlTest struct {
	suite.Suite
	graphqlFetcher Fetch
}

func TestPipelineSuite(t *testing.T) {
	suite.Run(t, new(itGraphqlTest))
}

func (s *itGraphqlTest) SetupSuite() {
	s.graphqlFetcher = NewGraphqlFetcher(GraphqlOpts{
		GraphqlEndpoint: graphqlEndpoint,
	})
}

func (s *itGraphqlTest) Test_E2E_Fetch_Existing_Block() {
	resp, err := s.graphqlFetcher.Block(14974600)
	s.NoError(err)
	s.Len(resp.Data.Block.Transactions, 3)
}

func (s *itGraphqlTest) Test_E2E_Fetch_Non_Existing_Block() {
	_, err := s.graphqlFetcher.Block(14974600000)
	s.Error(err)
}