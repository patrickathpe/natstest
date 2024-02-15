package main

import (
	"context"
	"fmt"
	"log"
	"time"

	cejs "github.com/cloudevents/sdk-go/protocol/nats_jetstream/v2"
	ce "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/binding"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

const (
	stream   = "rtm"
	consumer = "rtm_consumer"
	url      = "nats://nats-server:4222"
	userPath = "nats/keys/tester.creds"
)

var natsOpts = append([]nats.Option{},
	nats.MaxReconnects(3),
	nats.ReconnectWait(time.Second),
	nats.Timeout(5*time.Second),
	nats.Name("RTM client"),
	nats.UserCredentials(userPath),
)

func consume(ctx context.Context, conn *nats.Conn) (jetstream.ConsumeContext, error) {
	js, err := jetstream.New(conn)
	if err != nil {
		log.Println("error creating JetStream")
		return nil, err
	}

	consumer, err := js.Consumer(ctx, stream, consumer)
	if err != nil {
		log.Println("error getting consumer")
		return nil, err
	}

	return consumer.Consume(func(msg jetstream.Msg) {
		log.Printf("received message from %s\n", msg.Subject())
		if len(msg.Headers()) > 0 {
			log.Println("MESSAGE HEADERS")
			for key, value := range msg.Headers() {
				log.Printf("%s: %s\n", key, value)
			}
		}
		if len(msg.Data()) > 0 {
			log.Println("MESSAGE DATA")
			log.Println(string(msg.Data()))
		}

		// TODO: find a better way to convert jetstream.Msg to nats.Msg
		natsMsg := &nats.Msg{
			Subject: msg.Subject(),
			Data:    msg.Data(),
			Header:  msg.Headers(),
		}

		event, err := natsMessageToCloudEvent(ctx, natsMsg)
		if err != nil {
			log.Println("error converting NATS message to CloudEvent")
			nackMessage(msg)
			return
		}
		log.Println("CLOUD EVENT")
		log.Printf("%+v\n", event)
		ackMessage(msg)
	}, jetstream.PullMaxMessages(1))
}

func ackMessage(msg jetstream.Msg) {
	err := msg.Ack()
	if err != nil {
		log.Println("error acking message")
		log.Println(err)
	}
}

func nackMessage(msg jetstream.Msg) {
	err := msg.Nak()
	if err != nil {
		log.Println("error nacking message")
		log.Println(err)
	}
}

func natsMessageToCloudEvent(ctx context.Context, msg *nats.Msg) (*ce.Event, error) {
	ceMsg := cejs.NewMessage(msg)
	return binding.ToEvent(ctx, ceMsg)
}

func main() {
	ctx := context.Background()
	conn, err := nats.Connect(url, natsOpts...)
	if err != nil {
		log.Fatalln(err)
	}

	cc, err := consume(ctx, conn)
	if err != nil {
		log.Fatalln(err)
	}

	_, _ = fmt.Scanln()
	cc.Stop()
}
