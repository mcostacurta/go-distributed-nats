package core

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
)

func ConnectToNATS() (*nats.Conn, error) {
	nc, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		return nil, fmt.Errorf("nats connect: %w", err)
	}

	return nc, nil
}

func GetNatsJetStream(nc *nats.Conn) (nats.JetStreamContext, error) {

	// Create a context with a deadline - 2 seconds timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel() // Ensure the context is canceled to avoid leaks

	jsCtx, err := nc.JetStream(nats.Context(ctx))
	if err != nil {
		return nil, fmt.Errorf("jetstream: %w", err)
	}

	return jsCtx, nil
}

// AddStream adds a new stream to the JetStream server
func CreateStream(ctx context.Context, jsCtx nats.JetStreamContext) (*nats.StreamInfo, error) {
	stream, err := jsCtx.AddStream(&nats.StreamConfig{
		Name:              "test_stream",
		Subjects:          []string{"topic1", "topic2", "topic3"},
		Retention:         nats.InterestPolicy, // remove confirmed messages
		Discard:           nats.DiscardOld,     // when the stream is full, discard old messages
		MaxAge:            7 * 24 * time.Hour,  // max age of stored messages is 7 days
		Storage:           nats.FileStorage,    // type of message storage
		MaxMsgsPerSubject: 100_000_000,         // max stored messages per subject
		MaxMsgSize:        4 << 20,             // max single message size is 4 MB
		NoAck:             false,               // we need the "ack" system for the message queue system
	}, nats.Context(ctx))
	if err != nil {
		return nil, fmt.Errorf("add stream: %w", err)
	}

	fmt.Printf("Stream created: %+v\n", stream)
	return stream, nil
}

func PublishMsg(nc *nats.Conn, subject string, payload []byte) error {
	err := nc.Publish(subject, payload)
	if err != nil {
		return fmt.Errorf("publish: %w", err)
	}

	fmt.Printf("Message published: %s\n", payload)

	return nil
}

func CreateConsumer(ctx context.Context, jsCtx nats.JetStreamContext, consumerGroupName, streamName string) (*nats.ConsumerInfo, error) {
	consumer, err := jsCtx.AddConsumer(streamName, &nats.ConsumerConfig{
		Durable:       consumerGroupName,      // durable name is the same as consumer group name
		DeliverPolicy: nats.DeliverAllPolicy,  // deliver all messages, even if they were sent before the consumer was created
		AckPolicy:     nats.AckExplicitPolicy, // ack messages manually
		AckWait:       5 * time.Second,        // wait for ack for 5 seconds
		MaxAckPending: -1,                     // unlimited number of pending acks
	}, nats.Context(ctx))
	if err != nil {
		return nil, fmt.Errorf("add consumer: %w", err)
	}

	return consumer, nil
}

func Subscribe(ctx context.Context, js nats.JetStreamContext, subject, consumerGroupName, streamName string) (*nats.Subscription, error) {
	pullSub, err := js.PullSubscribe(
		subject,
		consumerGroupName,
		nats.ManualAck(),                         // ack messages manually
		nats.Bind(streamName, consumerGroupName), // bind consumer to the stream
		nats.Context(ctx),                        // use context to cancel the subscription
	)
	if err != nil {
		return nil, fmt.Errorf("pull subscribe: %w", err)
	}

	return pullSub, nil
}

// fetch first message from the subscription
func FetchOne(ctx context.Context, pullSub *nats.Subscription) (*nats.Msg, error) {

	// Create a context with a deadline - 2 seconds timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // Ensure the context is canceled to avoid leaks

	msgs, err := pullSub.Fetch(1, nats.Context(ctx))
	if err != nil {
		return nil, fmt.Errorf("fetch: %w", err)
	}
	if len(msgs) == 0 {
		return nil, errors.New("no messages")
	}

	return msgs[0], nil
}

// fetch all messages from the subscription
func FetchAll(ctx context.Context, pullSub *nats.Subscription, numMessages int) ([]*nats.Msg, error) {

	// Create a context with a deadline - 2 seconds timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // Ensure the context is canceled to avoid leaks

	msgs, err := pullSub.Fetch(numMessages, nats.Context(ctx))
	if err != nil {
		return nil, fmt.Errorf("fetch: %w", err)
	}

	return msgs, nil
}
