package splunk

import (
	"bytes"
	"net/http"
)

type Publisher struct {
	Config Config
}

func NewPublisher(config Config) *Publisher {
	return &Publisher{Config: config}
}

func (p *Publisher) Publish(payload []byte) error {
	_, err := http.NewRequest("POST", p.Config.Endpoint, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	return nil
}
