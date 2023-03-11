package pub

import (
	"encoding/json"
	"time"

	"github.com/nats-io/nats.go"
)

const (
	streamName     string = "CHAIN"
	streamSubjects string = "CHAIN.*"
)

type (
	PubOpts struct {
		DedupDuration   time.Duration
		JsCtx           nats.JetStreamContext
		NatsConn        *nats.Conn
		PersistDuration time.Duration
	}

	Pub struct {
		natsConn *nats.Conn
		jsCtx    nats.JetStreamContext
	}

	MinimalTxInfo struct {
		Block           uint64 `json:"block"`
		From            string `json:"from"`
		To              string `json:"to"`
		ContractAddress string `json:"contractAddress"`
		Success         bool   `json:"success"`
		TxHash          string `json:"transactionHash"`
		TxIndex         uint   `json:"transactionIndex"`
		Value           uint64 `json:"value"`
	}
)

func NewPub(o PubOpts) (*Pub, error) {
	stream, _ := o.JsCtx.StreamInfo(streamName)
	if stream == nil {
		_, err := o.JsCtx.AddStream(&nats.StreamConfig{
			Name:       streamName,
			MaxAge:     o.PersistDuration,
			Storage:    nats.FileStorage,
			Subjects:   []string{streamSubjects},
			Duplicates: o.DedupDuration,
		})
		if err != nil {
			return nil, err
		}
	}

	return &Pub{
		jsCtx:    o.JsCtx,
		natsConn: o.NatsConn,
	}, nil
}

// Close gracefully shutdowns the JetStream connection.
func (p *Pub) Close() {
	if p.natsConn != nil {
		p.natsConn.Close()
	}
}

// Publish publishes the JSON data to the NATS stream.
func (p *Pub) Publish(subject string, dedupId string, eventPayload interface{}) error {
	jsonData, err := json.Marshal(eventPayload)
	if err != nil {
		return err
	}

	_, err = p.jsCtx.Publish(subject, jsonData, nats.MsgId(dedupId))
	if err != nil {
		return err
	}

	return nil
}
