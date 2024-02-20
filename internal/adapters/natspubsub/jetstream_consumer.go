package natspubsub

import (
	"context"

	ce "github.com/cloudevents/sdk-go/v2"
	natsjs "github.com/nats-io/nats.go/jetstream"

	"github.hpe.com/cloud/go-gadgets/x/logging"

	"github.com/patrickathpe/natstest/internal/drivers/nats/jetstream"
)

type ConsumeCallback func(ctx context.Context, logger logging.Logger, event *ce.Event) error

// interface that can be mocked in tests
type JetStreamConnection interface {
	JetStream() (jetstream.JetStream, error)
}

type JetStreamConsumer struct {
	js jetstream.JetStream
}

func NewJetStreamConsumer(conn JetStreamConnection) (*JetStreamConsumer, error) {
	js, err := conn.JetStream()
	if err != nil {
		return nil, err
	}
	return &JetStreamConsumer{js: js}, nil
}

func (c *JetStreamConsumer) Consume(
	ctx context.Context,
	logger logging.Logger,
	stream, consumerName string,
	callback ConsumeCallback,
) (jetstream.ConsumeContext, error) {
	consumer, err := c.js.Consumer(ctx, stream, consumerName)
	if err != nil {
		// TODO: log and wrap in terrors InternalError
		return nil, err
	}

	// Consumer will continue to run on a background goroutine
	// until consumerContext.Stop() is called
	return consumer.Consume(func(msg natsjs.Msg) {
		event, err := JetStreamMessageToCloudEvent(ctx, msg)
		if err != nil {
			// TODO log err
			// TODO handle and log nak error
			// TODO determine if message should be naked. Message will be re-delivered
			//	and will probably have the same error. This is a deserialization error.
			_ = msg.Nak()
			return
		}

		// TODO: add fields to logger, e.g. tracing from the event
		// eventLogger := logger.WithField...
		err = callback(ctx, logger, event)
		if err != nil {
			// TODO log err
			// TODO handle and log nak error
			// TODO determine if message should be naked.
			_ = msg.Nak()
			return
		}

		// TODO handle and log ack error
		_ = msg.Ack()
	}, natsjs.PullMaxMessages(1))
}
