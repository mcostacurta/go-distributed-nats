package main

import (
	"context"
	"fmt"

	"github.com/mcostacurta/go-distributed-nats/core"
)

func main() {

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

	// Create a consumer
	consumerGroupName := "test_consumer"
	streamName := "test_stream"
	fmt.Println(">>>> Creating a consumer...")
	_, err = core.CreateConsumer(context.Background(), jsCtx, consumerGroupName, streamName)
	if err != nil {
		fmt.Printf("create consumer: %v\n", err)
		return
	}
	fmt.Println(">>>> Consumer created")

	// Subscribe to a topic
	subject := "topic1"
	fmt.Println(">>>> Subscribing to a topic...")
	pullSub, err := core.Subscribe(context.Background(), jsCtx, subject, consumerGroupName, streamName)
	if err != nil {
		fmt.Printf("subscribe: %v\n", err)
		return
	}
	fmt.Println(">>>> Subscribed to a topic")

	// Fetch a message
	fmt.Println(">>>> Fetching a message...")

	// Get total messages and bytes
	streamInfo, err := jsCtx.StreamInfo("test_stream")
	if err != nil {
		fmt.Printf("stream info: %v\n", err)
	}

	// Print total messages and bytes
	fmt.Printf("Total Messages: %d\n", streamInfo.State.Msgs)
	fmt.Printf("Total Bytes: %d\n", streamInfo.State.Bytes)

	// msgs, err := core.FetchOne(context.Background(), pullSub)
	msgs, err := core.FetchAll(context.Background(), pullSub, int(streamInfo.State.Msgs)-1)
	if err != nil {
		fmt.Printf("fetch one: %v\n", err)
		return
	}
	fmt.Println(">>>> Message fetched")

	//	print all messages
	for _, msg := range msgs {
		fmt.Printf("**** Message received: %s ****\n", msg.Data)
	}
}
