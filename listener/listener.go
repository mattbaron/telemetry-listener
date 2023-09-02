package listener

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Listener struct {
	server http.Server
	mux    http.ServeMux

	ServiceAddress string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration

	eventChannel chan (Event)

	writeCount int64
	pingCount  int64
}

type Event struct {
	Fatal   bool
	Message string
}

func MakeListener() *Listener {
	l := &Listener{
		ReadTimeout:    time.Second * 5,
		WriteTimeout:   time.Second * 5,
		ServiceAddress: ":8080",
		eventChannel:   nil,
	}

	return l
}

func (l *Listener) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	l.mux.ServeHTTP(res, req)
}

func (l *Listener) Start() {
	l.prepare()
	go func() {
		err := l.server.ListenAndServe()
		if err != nil {
			fmt.Println("Failed to start Listener")
		} else {
			fmt.Printf("Started server on %s\n", l.ServiceAddress)
			l.reportEvent(false, "Server Started")
		}
	}()
}

func (l *Listener) Stop() {
	fmt.Println("Stop()")
	l.server.Shutdown(context.Background())
}

func (l *Listener) AddEventListener(channel chan Event) {
	l.eventChannel = channel
}

func (l *Listener) reportEvent(fatal bool, message string) {
	if l.eventChannel != nil {
		l.eventChannel <- Event{Fatal: fatal, Message: message}
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
		handlePing(response, request)
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

func handlePing(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusOK)
	response.Write([]byte("Foobar"))
}
