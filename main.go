package main

import (
	"fmt"
	"os"

	"github.com/mattbaron/telemetry-listener/listener"
)

// func test() {
// 	h := func(w http.ResponseWriter, request *http.Request) {
// 		fmt.Println("Hello")
// 		fmt.Println(request.ContentLength)
// 	}

// 	http.HandleFunc("/test", h)
// 	http.ListenAndServe(":8080", nil)
// }

func main() {
	metricListner := listener.MakeListener()

	eventChannel := metricListner.NewEventChannel()
	go func() {
		for event := range eventChannel {
			if event.Type >= listener.ERROR {
				fmt.Printf("Fatal error: %s\n", event.Message)
				os.Exit(3)
			} else {
				fmt.Printf("Event: %d %s\n", event.Type, event.Message)
			}
		}
	}()

	metricListner.Start()
	metricListner.WaitUntilDone()
}
