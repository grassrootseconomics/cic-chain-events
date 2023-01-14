package fetch

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

var (
	graphqlEndpoint = os.Getenv("TEST_GRAPHQL_ENDPOINT")
)

type GraphQlTestSuite struct {
	suite.Suite
	fetch Fetch
}

func (s *GraphQlTestSuite) SetupSuite() {
	s.fetch = NewGraphqlFetcher(GraphqlOpts{
		GraphqlEndpoint: graphqlEndpoint,
	})
}

func (s *GraphQlTestSuite) Test_E2E_Fetch_Existing_Block() {
	resp, err := s.fetch.Block(context.Background(), 14974600)
	s.NoError(err)
	s.Len(resp.Data.Block.Transactions, 3)
}

func (s *GraphQlTestSuite) Test_E2E_Fetch_Non_Existing_Block() {
	_, err := s.fetch.Block(context.Background(), 14974600000)
	s.Error(err)
}

func TestGraphQlSuite(t *testing.T) {
	suite.Run(t, new(GraphQlTestSuite))
}
