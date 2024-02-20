package natspubsub

import (
	"bytes"
	"context"

	cejs "github.com/cloudevents/sdk-go/protocol/nats_jetstream/v2"
	ce "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/binding"
	natsio "github.com/nats-io/nats.go"
	natsjs "github.com/nats-io/nats.go/jetstream"
)

func CloudEventToMessage(ctx context.Context, sub string, event *ce.Event) (*natsio.Msg, error) {
	ceMsg := binding.ToMessage(event)
	writer := new(bytes.Buffer)
	header, err := cejs.WriteMsg(ctx, ceMsg, writer)
	if err != nil {
		return nil, err
	}

	return &natsio.Msg{
		Subject: sub,
		Data:    writer.Bytes(),
		Header:  header,
	}, nil
}

func MessageToCloudEvent(ctx context.Context, msg *natsio.Msg) (*ce.Event, error) {
	ceMsg := cejs.NewMessage(msg)
	return binding.ToEvent(ctx, ceMsg)
}

func JetStreamMessageToCloudEvent(ctx context.Context, msg natsjs.Msg) (*ce.Event, error) {
	// TODO: find a better way to convert jetstream.Msg to nats.Msg
	natsMsg := &natsio.Msg{
		Subject: msg.Subject(),
		Data:    msg.Data(),
		Header:  msg.Headers(),
	}
	return MessageToCloudEvent(ctx, natsMsg)
}
