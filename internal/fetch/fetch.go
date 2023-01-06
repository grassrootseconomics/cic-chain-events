package fetch

// Fetch defines a block fetcher that must return a full JSON response
type Fetch interface {
	Block(block uint64) (fetchResponse FetchResponse, err error)
}

// Transaction reprsents a JSON object of all important mined transaction information
type Transaction struct {
	Block struct {
		Number    uint   `json:"number"`
		Timestamp string `json:"timestamp"`
	} `json:"block"`
	Hash  string `json:"hash"`
	Index uint   `json:"index"`
	From  struct {
		Address string `json:"address"`
	} `json:"from"`
	To struct {
		Address string `json:"address"`
	} `json:"to"`
	Value     string `json:"value"`
	InputData string `json:"inputData"`
	Status    uint   `json:"status"`
	GasUsed   uint   `json:"gasUsed"`
}

// BlockResponse represents a full fetch JSON response
type FetchResponse struct {
	Data struct {
		Block struct {
			Transactions []Transaction `json:"transactions"`
		} `json:"block"`
	} `json:"data"`
}