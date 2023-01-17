package fetch

import (
	"context"
	"testing"

	celo "github.com/grassrootseconomics/cic-celo-sdk"
)

func Benchmark_RPC(b *testing.B) {
	celoProvider, err := celo.NewProvider(celo.ProviderOpts{
		ChainId:     celo.MainnetChainId,
		RpcEndpoint: rpcEndpoint,
	})

	rpc := NewRPCFetcher(RPCOpts{
		RPCProvider: celoProvider,
	})

	if err != nil {
		return
	}

	b.Run("RPC_Block_Fetcher_Benchmark", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			_, err := rpc.Block(context.Background(), 14974600)
			if err != nil {
				b.Fatal(err)
			}
		}
		b.ReportAllocs()
	})
}

func Benchmark_GraphQL(b *testing.B) {
	graphql := NewGraphqlFetcher(GraphqlOpts{
		GraphqlEndpoint: graphqlEndpoint,
	})

	b.Run("GraphQL_Block_Fetcher_Benchmark", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			_, err := graphql.Block(context.Background(), 14974600)
			if err != nil {
				b.Fatal(err)
			}
		}
		b.ReportAllocs()
	})

}
