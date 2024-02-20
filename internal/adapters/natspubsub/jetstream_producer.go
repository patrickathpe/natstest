package natspubsub

import (
	"context"

	ce "github.com/cloudevents/sdk-go/v2"

	"github.hpe.com/cloud/go-gadgets/x/logging"

	"github.com/patrickathpe/natstest/internal/drivers/nats/jetstream"
)

type JetStreamProducer struct {
	js jetstream.JetStream
}

func NewJetStreamProducer(conn JetStreamConnection) (*JetStreamProducer, error) {
	js, err := conn.JetStream()
	if err != nil {
		return nil, err
	}
	return &JetStreamProducer{js: js}, nil
}

func (producer *JetStreamProducer) Produce(ctx context.Context, logger logging.Logger, subject string, event *ce.Event) error {
	msg, err := CloudEventToMessage(ctx, subject, event)
	if err != nil {
		// TODO: log and wrap in terror InternalError
		return err
	}

	_, err = producer.js.PublishMsg(ctx, msg)
	if err != nil {
		// TODO: log and wrap in terror InternalError
		return err
	}
	return nil
}
