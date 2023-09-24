package splunk

import (
	"sync"
)

type EventBuffer struct {
	BatchSize     int32
	FlushInterval int32
	events        []*Event
	started       bool
	mutex         sync.Mutex
}

func NewEventBuffer() *EventBuffer {
	return &EventBuffer{
		BatchSize:     25,
		FlushInterval: 5,
		events:        make([]*Event, 0),
		started:       false,
	}
}

func (b *EventBuffer) BufferEvent(event *Event) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.events = append(b.events, event)
}

func (b *EventBuffer) FlushEvents() {
	b.mutex.Lock()
	batch := b.events
	b.events = make([]*Event, 0)
	b.mutex.Unlock()

	buff := make([]byte, 0)
	for _, event := range batch {
		eventBytes, _ := event.Serialize()
		buff = append(buff, eventBytes...)
	}

}
