package splunk

import (
	"fmt"
	"os"
	"sync"
	"time"
)

var defaultLogger *logger
var defaultConfig *Config

type logger struct {
	config  *Config
	events  []*Event
	started bool
	mutex   sync.Mutex
}

type Config struct {
	Endpoint      string
	Token         string
	Source        string
	SourceType    string
	Host          string
	FlushInterval int
	BatchSize     int
	Defaults      map[string]any
}

func NewConfig() *Config {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	return &Config{
		Defaults:      make(map[string]any, 0),
		Host:          hostname,
		FlushInterval: 5,
		BatchSize:     25,
	}
}

func NewLogger(config *Config) *logger {
	if config == nil {
		config = NewConfig()
	}

	logger := &logger{
		started: false,
		config:  config,
	}

	logger.Start()
	return logger
}

func Configure(config *Config) {
	defaultConfig = config
}

func Logger() *logger {
	if defaultLogger == nil {
		defaultLogger = NewLogger(defaultConfig)
	}

	return defaultLogger
}

func (l *logger) Start() {

	fmt.Println(l.config)
	if l.started {
		return
	}

	go func() {
		for {
			time.Sleep(time.Duration(l.config.FlushInterval) * time.Second)
			l.FlushEvents()
		}
	}()

	l.started = true
}

func (l *logger) Event(event *Event) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.events = append(l.events, event)
}

func (l *logger) Info(message string) {
	l.Event(NewEvent().Set("message", message).Set("severity", "info"))
}

func (l *logger) Warn(message string) {
	l.Event(NewEvent().Set("message", message).Set("severity", "warn"))
}

func (l *logger) Error(message string) {
	l.Event(NewEvent().Set("message", message).Set("severity", "error"))
}

func (l *logger) FlushEvents() {
	fmt.Println("FlushEvents()")
	l.mutex.Lock()
	batch := l.events
	l.events = make([]*Event, 0)
	l.mutex.Unlock()

	buff := make([]byte, 0)
	for _, event := range batch {
		eventBytes, _ := event.Serialize()
		buff = append(buff, eventBytes...)
	}

}

func (l *logger) publish(payload []byte) error {

	return nil
}
