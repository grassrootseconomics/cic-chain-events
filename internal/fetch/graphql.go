package fetch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	contentType  = "application/json"
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
			Timeout: time.Second * 5,
		},
		graphqlEndpoint: o.GraphqlEndpoint,
	}
}

func (f *Graphql) Block(blockNumber uint64) (FetchResponse, error) {
	var (
		fetchResponse FetchResponse
	)

	resp, err := f.httpClient.Post(
		f.graphqlEndpoint,
		contentType,
		bytes.NewBufferString(fmt.Sprintf(graphqlQuery, blockNumber)),
	)
	if err != nil {
		return FetchResponse{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return FetchResponse{}, fmt.Errorf("error fetching block %s", resp.Status)
	}
	defer resp.Body.Close()

	out, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return FetchResponse{}, nil
	}

	if err := json.Unmarshal(out, &fetchResponse); err != nil {
		return FetchResponse{}, err
	}

	return fetchResponse, nil
}
