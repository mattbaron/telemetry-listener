package listener

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	TRACE = iota
	DEBUG = iota
	INFO  = iota
	WARN  = iota
	ERROR = iota
	FATAL = iota
)

type Listener struct {
	ServiceAddress string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	server         http.Server
	mux            http.ServeMux
	eventChannels  []chan Event
	writeCount     int64
	pingCount      int64
}

type Event struct {
	Message string
	Type    int
}

func MakeListener() *Listener {
	l := &Listener{
		ReadTimeout:    time.Second * 5,
		WriteTimeout:   time.Second * 5,
		ServiceAddress: ":8080",
		eventChannels:  make([]chan Event, 0),
	}

	return l
}

func (l *Listener) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	l.mux.ServeHTTP(res, req)
}

func (l *Listener) Start() {
	l.prepare()

	go func() {
		l.event("Staring server", INFO)
		err := l.server.ListenAndServe()
		if err != nil {
			l.event("Failed to start server", FATAL)
		}
	}()
}

func (l *Listener) WaitUntilDone() {
	<-make(chan int)
}

func (l *Listener) Stop() {
	l.server.Shutdown(context.Background())
}

func (l *Listener) NewEventChannel() chan Event {
	channel := make(chan Event)
	l.eventChannels = append(l.eventChannels, channel)
	return channel
}

func (l *Listener) event(message string, eventType int) {
	for _, channel := range l.eventChannels {
		channel <- Event{Message: message, Type: eventType}
	}
}

func (l *Listener) prepare() {
	l.server = http.Server{
		Handler:      l,
		Addr:         l.ServiceAddress,
		WriteTimeout: l.WriteTimeout,
		ReadTimeout:  l.ReadTimeout,
	}

	l.mux.HandleFunc("/ping", (func(response http.ResponseWriter, request *http.Request) {
		l.pingCount += 1
		fmt.Printf("%d /ping\n", l.pingCount)
		response.WriteHeader(http.StatusOK)
		response.Write([]byte("Foobar"))
	}))

	l.mux.HandleFunc("/write", (func(response http.ResponseWriter, request *http.Request) {
		l.writeCount += 1

		body, err := io.ReadAll(request.Body)
		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
		}

		fmt.Printf("%d /write\n%s\n", l.pingCount, body)
	}))
}
