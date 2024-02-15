package main

import (
	"bytes"
	"context"
	"log"
	"time"

	cejs "github.com/cloudevents/sdk-go/protocol/nats_jetstream/v2"
	ce "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/binding"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"

	"github.com/patrickathpe/natstest/internal/adapters/dbtesting"
)

const (
	stream   = "rtm"
	consumer = "rtm_consumer"
	url      = "nats://nats-server:4222"
	userPath = "nats/keys/tester.creds"
	subject  = "rtm.events.foo"
)

var natsOpts = append([]nats.Option{},
	nats.MaxReconnects(3),
	nats.ReconnectWait(time.Second),
	nats.Timeout(5*time.Second),
	nats.Name("RTM client"),
	nats.UserCredentials(userPath),
)

func produce(ctx context.Context, conn *nats.Conn) error {
	js, err := jetstream.New(conn)
	if err != nil {
		log.Println("error creating JetStream")
		return err
	}

	// Create a random CloudEvent
	generator := dbtesting.CloudEventGenerator{}
	event := generator.GetCloudEvent()
	msg, err := cloudEventToNatsMessage(ctx, subject, event)
	if err != nil {
		return err
	}

	_, err = js.PublishMsg(ctx, msg)
	if err != nil {
		log.Println("error publishing message")
		return err
	}
	return nil
}

func cloudEventToNatsMessage(ctx context.Context, sub string, event *ce.Event) (*nats.Msg, error) {
	ceMsg := binding.ToMessage(event)
	writer := new(bytes.Buffer)
	header, err := cejs.WriteMsg(ctx, ceMsg, writer)
	if err != nil {
		log.Println("error writing event message")
		return nil, err
	}

	return &nats.Msg{
		Subject: sub,
		Data:    writer.Bytes(),
		Header:  header,
	}, nil
}

func main() {
	ctx := context.Background()
	conn, err := nats.Connect(url, natsOpts...)
	if err != nil {
		log.Fatalln(err)
	}

	err = produce(ctx, conn)
	if err != nil {
		log.Fatalln(err)
	}
}
