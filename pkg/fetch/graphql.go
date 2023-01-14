package fetch

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/goccy/go-json"
)

const (
	graphqlQuery = `{"query":"{block(number:%d){transactions{block{number,timestamp},hash,index,from{address},to{address},value,inputData,status,gasUsed}}}"}`
)

type GraphqlOpts struct {
	GraphqlEndpoint string
}

type Graphql struct {
	graphqlEndpoint string
	httpClient      *http.Client
}

func NewGraphqlFetcher(o GraphqlOpts) Fetch {
	return &Graphql{
		httpClient: &http.Client{
			Timeout: time.Second * 2,
		},
		graphqlEndpoint: o.GraphqlEndpoint,
	}
}

func (f *Graphql) Block(ctx context.Context, blockNumber uint64) (FetchResponse, error) {
	var (
		fetchResponse FetchResponse
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, f.graphqlEndpoint, bytes.NewBufferString(fmt.Sprintf(graphqlQuery, blockNumber)))
	if err != nil {
		return FetchResponse{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := f.httpClient.Do(req)
	if err != nil {
		return FetchResponse{}, err
	}

	if resp.StatusCode >= http.StatusBadRequest {
		return FetchResponse{}, fmt.Errorf("error fetching block %s", resp.Status)
	}

	out, err := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		return FetchResponse{}, nil
	}

	if err := json.Unmarshal(out, &fetchResponse); err != nil {
		return FetchResponse{}, err
	}

	return fetchResponse, nil
}
