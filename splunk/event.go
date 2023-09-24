package splunk

import (
	"encoding/json"
	"fmt"
	"time"
)

type Event struct {
	Source     string         `json:"source,omitempty"`
	SourceType string         `json:"sourcetype,omitempty"`
	Host       string         `json:"host,omitempty"`
	Time       int64          `json:"time"`
	Values     map[string]any `json:"event"`
}

func NewEvent() *Event {
	event := &Event{
		Values: make(map[string]any),
		Time:   time.Now().Unix(),
	}

	if defaultConfig != nil {
		event.Source = defaultConfig.Source
		event.SourceType = defaultConfig.SourceType
		event.Host = defaultConfig.Host
		for k, v := range defaultConfig.Defaults {
			event.Set(k, v)
		}
	}

	return event
}

func (e *Event) Set(key string, value any) *Event {
	e.Values[key] = value
	return e
}

func (e *Event) Dump() {
	fmt.Println(e.Values)
}

func (e *Event) Serialize() ([]byte, error) {
	value, err := json.Marshal(e)
	if err != nil {
		return make([]byte, 0), err
	}

	return value, err
}
