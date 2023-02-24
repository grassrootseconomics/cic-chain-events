package events

type EventEmitter interface {
	Close()
	Publish(subject string, dedupId string, eventPayload interface{}) error
}

type MinimalTxInfo struct {
	Block           uint64 `json:"block"`
	From            string `json:"from"`
	To              string `json:"to"`
	ContractAddress string `json:"contractAddress"`
	Success         bool   `json:"success"`
	TxHash          string `json:"transactionHash"`
	TxIndex         uint   `json:"transactionIndex"`
	Value           uint64 `json:"value"`
}
