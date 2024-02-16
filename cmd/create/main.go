package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

const (
	stream   = "rtm"
	consumer = "rtm_consumer"
	url      = "nats://nats-server:4222"
	userPath = "nats/keys/tester.creds"
	maxMsgs  = int64(250000)
	maxAge   = (24 * 7) * time.Hour
	maxBytes = (1 << 30) * 2 // 2 GiB

	sourceSubject = "rtm.events.*"
	numPartitions = 3
)

var destSubject = fmt.Sprintf("rtm.events.{{partition(%d, 1)}}", numPartitions)

var streamCfg = jetstream.StreamConfig{
	Name:        stream,
	Description: "RTM Events",
	Subjects:    []string{sourceSubject},

	// Message retention limits
	MaxMsgs:  maxMsgs,
	MaxAge:   maxAge,
	MaxBytes: maxBytes,

	// Retain until one of limits above is reached
	// Retention: jetstream.LimitsPolicy,

	// Retain until acknowledged by all active consumers
	// Above limits still apply
	Retention: jetstream.InterestPolicy,

	// Retain until acknowledged by one consumer
	// Above limits still apply
	// Retention: jetstream.WorkQueuePolicy,

	// MaxConsumers         int
	// Discard              DiscardPolicy
	// DiscardNewPerSubject bool
	// MaxMsgsPerSubject    int64
	// MaxMsgSize           int32
	// Storage              StorageType
	Replicas: 1,
	// NoAck                bool
	// Template             string
	// Duplicates           time.Duration
	// Placement            *Placement
	// Mirror               *StreamSource
	// Sources              []*StreamSource
	// Sealed               bool
	// DenyDelete           bool
	// DenyPurge            bool
	// AllowRollup          bool
	// Compression          StoreCompression
	// FirstSeq             uint64

	// Allow applying a subject transform to incoming messages before doing anything else
	SubjectTransform: &jetstream.SubjectTransformConfig{
		Source:      sourceSubject,
		Destination: destSubject,
	},

	// Allow republish of the message after being sequenced and stored.
	// RePublish *RePublish

	// Allow higher performance, direct access to get individual messages. E.g. KeyValue
	// AllowDirect bool

	// Allow higher performance and unified direct access for mirrors as well.
	// MirrorDirect bool

	// Limits for consumers on this stream.
	// ConsumerLimits StreamConsumerLimits

	// Metadata is additional metadata for the Stream.
	// Keys starting with `_nats` are reserved.
	// NOTE: Metadata requires nats-server v2.10.0+
	// Metadata map[string]string
}

var consumerCfg = jetstream.ConsumerConfig{
	Name:    consumer,
	Durable: consumer,

	// Description        string
	// DeliverPolicy      DeliverPolicy
	// OptStartSeq        uint64
	// OptStartTime       *time.Time

	// Set the following for ordered message processing
	// AckPolicy:     jetstream.AckExplicitPolicy,
	AckWait:       5 * time.Second,
	MaxDeliver:    3,
	MaxAckPending: 1,

	// BackOff            []time.Duration
	// FilterSubject      string
	// ReplayPolicy       ReplayPolicy
	// RateLimit          uint64
	// SampleFrequency    string
	// MaxWaiting         int
	// HeadersOnly        bool
	// MaxRequestBatch    int
	// MaxRequestExpires  time.Duration
	// MaxRequestMaxBytes int

	// Inactivity threshold.
	// InactiveThreshold time.Duration

	// Generally inherited by parent stream and other markers, now can be configured directly.
	// Replicas int `json:"num_replicas"`
	// Force memory storage.
	// MemoryStorage bool

	// NOTE: FilterSubjects requires nats-server v2.10.0+
	FilterSubjects: []string{"rtm.events.*"},

	// Metadata is additional metadata for the Consumer.
	// Keys starting with `_nats` are reserved.
	// NOTE: Metadata requires nats-server v2.10.0+
	// Metadata map[string]string `json:"metadata,omitempty"`
}

var natsOpts = append([]nats.Option{},
	nats.MaxReconnects(3),
	nats.ReconnectWait(time.Second),
	nats.Timeout(5*time.Second),
	nats.Name("RTM client"),
	nats.UserCredentials(userPath),
)

func create(ctx context.Context, conn *nats.Conn) error {
	js, err := jetstream.New(conn)
	if err != nil {
		log.Println("error creating JetStream")
		return err
	}

	log.Println("upserting stream")
	_, err = js.CreateOrUpdateStream(ctx, streamCfg)
	if err != nil {
		log.Println("error upserting stream")
		return err
	}

	// TODO: create one consumer for each partition
	log.Println("upserting consumer")
	_, err = js.CreateOrUpdateConsumer(ctx, stream, consumerCfg)
	if err != nil {
		log.Println("error creating or updating consumer")
		return err
	}
	log.Println("upsert consumer success")
	return nil
}

func main() {
	ctx := context.TODO()
	conn, err := nats.Connect(url, natsOpts...)
	if err != nil {
		log.Fatalln(err)
	}

	err = create(ctx, conn)
	if err != nil {
		log.Fatalln(err)
	}
}
