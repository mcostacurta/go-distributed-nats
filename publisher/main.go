package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/mcostacurta/go-distributed-nats/core"
)

func main() {

	message := flag.String("msg", "", "Message to publish")
	flag.Parse()
	fmt.Println(">>>> MESSAGE: ", *message)

	// Connect to NATS
	fmt.Println(">>>> Connecting to NATS...")
	nc, err := core.ConnectToNATS()
	if err != nil {
		fmt.Printf("connect to nats: %v\n", err)
		return
	}
	defer nc.Close()

	fmt.Println(">>>> Connected to NATS")

	// Get JetStream context
	fmt.Println(">>>> Getting JetStream context...")
	jsCtx, err := core.GetNatsJetStream(nc)
	if err != nil {
		fmt.Printf("get jetstream: %v\n", err)
		return
	}
	fmt.Println(">>>> Got JetStream context")

	// Create a stream
	fmt.Println(">>>> Creating a stream...")
	_, err = core.CreateStream(context.Background(), jsCtx)
	if err != nil {
		fmt.Printf("create stream: %v\n", err)
		return
	}
	fmt.Println(">>>> Stream created")

	// Create a payload
	payload := []byte(*message)

	// Publish a message
	fmt.Println(">>>> Publishing a message...")
	err = core.PublishMsg(nc, "topic1", payload)
	if err != nil {
		fmt.Printf("publish msg: %v\n", err)
		return
	}
	fmt.Println(">>>> Message published")

}
