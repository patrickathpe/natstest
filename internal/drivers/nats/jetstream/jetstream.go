package jetstream

import (
	"context"

	natsio "github.com/nats-io/nats.go"
	natsjs "github.com/nats-io/nats.go/jetstream"
)

// interface so that it can be mocked
type JetStream interface {
	// jetstream.PubAck is a struct with fields an no methods. No mock needed for it.
	PublishMsg(ctx context.Context, msg *natsio.Msg, opts ...natsjs.PublishOpt) (*natsjs.PubAck, error)

	Consumer(ctx context.Context, stream, consumerName string) (Consumer, error)

	CreateOrUpdateConsumer(ctx context.Context, stream string, cfg natsjs.ConsumerConfig) (Consumer, error)

	CreateOrUpdateStream(ctx context.Context, cfg natsjs.StreamConfig) (Stream, error)
}

type Consumer interface {
	Info(context.Context) (*natsjs.ConsumerInfo, error)
	Consume(handler natsjs.MessageHandler, opts ...natsjs.PullConsumeOpt) (ConsumeContext, error)
}

type ConsumeContext interface {
	Stop()
}

type Stream interface {
	Info(ctx context.Context, opts ...natsjs.StreamInfoOpt) (*natsjs.StreamInfo, error)
}

// jetstream wraps a NATS JetStream so it can be mocked
type jetstream struct {
	js natsjs.JetStream
}

// NewJetStream creates and returns a wrapped NATS JetStream
func NewJetStream(conn *natsio.Conn) (JetStream, error) {
	js, err := natsjs.New(conn)
	if err != nil {
		return nil, err
	}
	return &jetstream{js: js}, nil
}

func (w *jetstream) PublishMsg(ctx context.Context, msg *natsio.Msg, opts ...natsjs.PublishOpt) (*natsjs.PubAck, error) {
	return w.js.PublishMsg(ctx, msg, opts...)
}

func (w *jetstream) Consumer(ctx context.Context, stream, consumerName string) (Consumer, error) {
	jsConsumer, err := w.js.Consumer(ctx, stream, consumerName)
	if err != nil {
		return nil, err
	}
	return &consumer{consumer: jsConsumer}, nil
}

func (w *jetstream) CreateOrUpdateConsumer(ctx context.Context, stream string, cfg natsjs.ConsumerConfig) (Consumer, error) {
	jsConsumer, err := w.js.CreateOrUpdateConsumer(ctx, stream, cfg)
	if err != nil {
		return nil, err
	}
	return &consumer{consumer: jsConsumer}, nil
}

func (w *jetstream) CreateOrUpdateStream(ctx context.Context, cfg natsjs.StreamConfig) (Stream, error) {
	jsStream, err := w.js.CreateOrUpdateStream(ctx, cfg)
	if err != nil {
		return nil, err
	}
	return &stream{stream: jsStream}, nil
}

// consumer wraps a NATS JetStream Consumer so it can be mocked
type consumer struct {
	consumer natsjs.Consumer
}

func (c *consumer) Info(ctx context.Context) (*natsjs.ConsumerInfo, error) {
	return c.consumer.Info(ctx)
}

func (c *consumer) Consume(handler natsjs.MessageHandler, opts ...natsjs.PullConsumeOpt) (ConsumeContext, error) {
	jsContext, err := c.consumer.Consume(handler, opts...)
	if err != nil {
		return nil, err
	}
	return &consumeContext{consumeContext: jsContext}, nil
}

// consumeContext wraps a NATS JetStream ConsumeContext so it can be mocked
type consumeContext struct {
	consumeContext natsjs.ConsumeContext
}

func (c *consumeContext) Stop() {
	c.consumeContext.Stop()
}

// stream wraps a NATS JetStream Stream so it can be mocked
type stream struct {
	stream natsjs.Stream
}

func (s *stream) Info(ctx context.Context, opts ...natsjs.StreamInfoOpt) (*natsjs.StreamInfo, error) {
	return s.stream.Info(ctx, opts...)
}
