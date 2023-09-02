package main

import (
	"fmt"
	"net/http"

	"github.com/mattbaron/telemetry-listener/listener"
)

func test() {
	h := func(w http.ResponseWriter, request *http.Request) {
		fmt.Println("Hello")
		fmt.Println(request.ContentLength)
	}

	http.HandleFunc("/test", h)
	http.ListenAndServe(":8080", nil)
}

func main() {
	l := listener.MakeListener()

	eventChannel := make(chan listener.Event)
	l.AddEventListener(eventChannel)
	go func() {
		for event := range eventChannel {
			fmt.Println(event.Message)
		}
	}()

	l.Start()

	<-make(chan int)
}
