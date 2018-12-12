package main

import (
	"fmt"
	"github.com/go-stomp/stomp"
	"models/shared"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 4 {
		return
	}

	mqBroker := os.Args[1]
	user := os.Args[2]
	pass := os.Args[3]

	topic := "some-topic"

	cli := shared.NewStompClient(mqBroker, user, pass)
	cli.Send(topic, "andrew test message")

	handler := func(destination, message string) {
		fmt.Printf("received: %s on %s\n", message, destination)
	}

	cli.Subscribe(topic, handler)
	time.Sleep(5 * time.Second)
}
