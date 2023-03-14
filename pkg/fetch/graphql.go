package fetch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
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
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 100,
				IdleConnTimeout:     60 * time.Second,
			},
			Timeout: time.Second * 5,
		},
		graphqlEndpoint: o.GraphqlEndpoint,
	}
}

func (f *Graphql) Block(ctx context.Context, blockNumber uint64) (FetchResponse, error) {
	fetchResponse := FetchResponse{}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, f.graphqlEndpoint, bytes.NewBufferString(fmt.Sprintf(graphqlQuery, blockNumber)))
	if err != nil {
		return fetchResponse, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := f.httpClient.Do(req)
	if err != nil {
		return fetchResponse, err
	}

	if resp.StatusCode >= http.StatusBadRequest {
		return fetchResponse, fmt.Errorf("error fetching block %s", resp.Status)
	}

	out, err := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		return fetchResponse, nil
	}

	if err := json.Unmarshal(out, &fetchResponse); err != nil {
		return fetchResponse, err
	}

	return fetchResponse, nil
}
