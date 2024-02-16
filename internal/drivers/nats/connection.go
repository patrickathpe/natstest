package nats

import (
	"time"

	natsio "github.com/nats-io/nats.go"

	"github.com/patrickathpe/natstest/internal/adapters/config"
	"github.com/patrickathpe/natstest/internal/drivers/nats/jetstream"
)

type Connection struct {
	conn *natsio.Conn
}

func NewConnection(cfg config.NATSConnectionConfig) (*Connection, error) {
	opts := append([]natsio.Option{},
		natsio.MaxReconnects(cfg.MaxReconnects),
		natsio.ReconnectWait(time.Duration(cfg.ReconnectWaitMilliSecs)*time.Millisecond),
		natsio.Timeout(time.Duration(cfg.TimeoutMilliSecs)*time.Millisecond),
		natsio.UserCredentials(cfg.CredentialsPath),
	)

	conn, err := natsio.Connect(cfg.URL, opts...)
	if err != nil {
		// TODO: return terror InternalError
		return nil, err
	}
	return &Connection{conn: conn}, nil
}

func (c *Connection) JetStream() (jetstream.JetStream, error) {
	return jetstream.NewJetStream(c.conn)
}
