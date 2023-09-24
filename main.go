package main

import (
	"time"

	"github.com/mattbaron/telemetry-listener/splunk"
)

// func test() {
// 	h := func(w http.ResponseWriter, request *http.Request) {
// 		fmt.Println("Hello")
// 		fmt.Println(request.ContentLength)
// 	}

// 	http.HandleFunc("/test", h)
// 	http.ListenAndServe(":8080", nil)
// }

// func main() {
// 	metricListner := listener.MakeListener()

// 	eventChannel := metricListner.NewEventChannel()
// 	go func() {
// 		for event := range eventChannel {
// 			if event.Type >= listener.ERROR {
// 				fmt.Printf("Fatal error: %s\n", event.Message)
// 				os.Exit(3)
// 			} else {
// 				fmt.Printf("Event: %d %s\n", event.Type, event.Message)
// 			}
// 		}
// 	}()

// 	metricListner.Start()
// 	metricListner.WaitUntilDone()
// }

// func main() {
// 	amqp := publisher.NewAMQP()

// 	amqp.Brokers = []string{"amqp://localhost:5672", "amqp://localhost:5672"}

// 	err := amqp.Connect()
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	<-make(chan int)
// }

func main() {
	// splunkConfig := splunk.NewConfig()
	// splunkConfig.BatchSize = 30
	// splunkConfig.FlushInterval = 3
	// splunkConfig.Source = "Foo"
	// splunkConfig.SourceType = "Bar"
	// splunkConfig.Defaults["env"] = "production"
	// splunk.Configure(splunkConfig)
	// fmt.Println(splunkConfig)

	splunk.Logger().Info("This is a test")
	splunk.Logger().Error("An exception happened")

	splunk.Logger().FlushEvents()

	time.Sleep(10 * time.Second)

	//fmt.Println(data)

	// metricListner := listener.MakeListener()
	// metricListner.Start()
	// metricListner.WaitUntilDone()
}
