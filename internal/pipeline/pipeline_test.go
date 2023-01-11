package pipeline

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/grassrootseconomics/cic-chain-events/internal/fetch"
	"github.com/grassrootseconomics/cic-chain-events/internal/filter"
	"github.com/stretchr/testify/suite"
	"github.com/zerodha/logf"
)

var (
	graphqlEndpoint = os.Getenv("TEST_GRAPHQL_ENDPOINT")
)

type itPipelineTest struct {
	suite.Suite
	errorPipeline     *Pipeline
	normalPipeline    *Pipeline
	earlyExitPipeline *Pipeline
}

func TestPipelineSuite(t *testing.T) {
	suite.Run(t, new(itPipelineTest))
}

type errorFilter struct{}

func newErrorFilter() filter.Filter {
	return &errorFilter{}
}

func (f *errorFilter) Execute(ctx context.Context, transaction fetch.Transaction) (bool, error) {
	return false, errors.New("crash")
}

type earlyExitFilter struct{}

func newEarlyExitFilter() filter.Filter {
	return &earlyExitFilter{}
}

func (f *earlyExitFilter) Execute(ctx context.Context, transaction fetch.Transaction) (bool, error) {
	return false, nil
}

func (s *itPipelineTest) SetupSuite() {
	logger := logf.New(
		logf.Opts{
			Level: logf.DebugLevel,
		},
	)

	fetcher := fetch.NewGraphqlFetcher(fetch.GraphqlOpts{
		GraphqlEndpoint: graphqlEndpoint,
	})

	noopFilter := filter.NewNoopFilter(filter.NoopFilterOpts{
		Logg: logger,
	})
	errFilter := newErrorFilter()
	earlyFilter := newEarlyExitFilter()

	s.errorPipeline = NewPipeline(PipelineOpts{
		Filters: []filter.Filter{
			noopFilter,
			errFilter,
		},
		BlockFetcher: fetcher,
		Logg:         logger,
	})

	s.normalPipeline = NewPipeline(PipelineOpts{
		Filters: []filter.Filter{
			noopFilter,
		},
		BlockFetcher: fetcher,
		Logg:         logger,
	})

	s.earlyExitPipeline = NewPipeline(PipelineOpts{
		Filters: []filter.Filter{
			noopFilter,
			earlyFilter,
			errFilter,
		},
		BlockFetcher: fetcher,
		Logg:         logger,
	})
}

func (s *itPipelineTest) Test_E2E_Pipeline_Run_On_Existing_Block_No_Err() {
	err := s.normalPipeline.Run(context.Background(), 14974600)
	s.NoError(err)
}

func (s *itPipelineTest) Test_E2E_Pipeline_Run_On_Non_Existing_Block_No_Err() {
	err := s.normalPipeline.Run(context.Background(), 14974600000)
	s.Error(err)
}

func (s *itPipelineTest) Test_E2E_Pipeline_Run_On_Existing_Block_Early() {
	err := s.earlyExitPipeline.Run(context.Background(), 14974600)
	s.NoError(err)
}

func (s *itPipelineTest) Test_E2E_Pipeline_Run_On_Existing_Block_With_Err() {
	err := s.errorPipeline.Run(context.Background(), 14974600)
	s.Error(err)
}

func (s *itPipelineTest) Test_E2E_Pipeline_Run_On_Non_Existing_Block_With_Err() {
	err := s.errorPipeline.Run(context.Background(), 14974600000)
	s.Error(err)
}

func (s *itPipelineTest) Test_E2E_Pipeline_Run_On_Non_Existing_Block_Early() {
	err := s.earlyExitPipeline.Run(context.Background(), 14974600000)
	s.Error(err)
}

func (s *itPipelineTest) Test_E2E_Pipeline_Run_On_Empty_Block_With_No_Err() {
	err := s.normalPipeline.Run(context.Background(), 15370320)
	s.NoError(err)
}

func (s *itPipelineTest) Test_E2E_Pipeline_Run_On_Empty_Block_With_Err() {
	err := s.errorPipeline.Run(context.Background(), 15370320)
	s.NoError(err)
}

func (s *itPipelineTest) Test_E2E_Pipeline_Run_On_Empty_Block_Early() {
	err := s.earlyExitPipeline.Run(context.Background(), 15370320)
	s.NoError(err)
}
