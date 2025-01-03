package nats

import (
	"context"
	"encoding/json"
	"fmt"

	nats "github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	geEvent "github.com/grassrootseconomics/eth-tracker/pkg/event"
	"git.defalsify.org/vise.git/logging"
	"git.grassecon.net/urdt/ussd/common"
	"git.grassecon.net/term/event"
	"git.grassecon.net/term/config"
)

var (
	logg = logging.NewVanilla().WithDomain("term-nats")
)

// NatsSubscription encapsulates the jetstream session providing events.
//
// Extends Router.
type NatsSubscription struct {
	event.Router
	ctx context.Context
	conn *nats.Conn
	js jetstream.JetStream
	cs jetstream.Consumer
	cctx jetstream.ConsumeContext
}

// NewNatsSubscription creates a new NatsSubscription with the given user store.
func NewNatsSubscription(store common.StorageServices) *NatsSubscription {
	return &NatsSubscription{
		Router: event.Router{
			Store: store,
		},
	}
}

// Connect sets up the connection to the nats server and a consumer for the 
// "Jetstream".
//
// Fails if connection fails or the "Jetstream" consumer cannot be set up.
//
// Once connected, it will attempt to reconnect if disconnected.
func(n *NatsSubscription) Connect(ctx context.Context, connStr string) error {
	var err error

	n.conn, err = nats.Connect(connStr)
	if err != nil {
		return err
	}
	n.conn.SetDisconnectErrHandler(disconnectHandler)
	n.conn.SetReconnectHandler(reconnectHandler)
	n.js, err = jetstream.New(n.conn)
	if err != nil {
		return err
	}
	n.cs, err = n.js.CreateConsumer(ctx, "TRACKER", jetstream.ConsumerConfig{
		Name: config.JetstreamClientName,
		Durable: config.JetstreamClientName,
		FilterSubjects: []string{"TRACKER.*"},
	})
	if err != nil {
		return err
	}

	serverInfo := toServerInfo(n.conn)
	logg.DebugCtxf(ctx, "nats connected, starting consumer", "status", n.conn.Status(), "server", serverInfo)
	n.ctx = ctx
	n.cctx, err = n.cs.Consume(n.handleEvent)
	if err != nil {
		return err		
	}

	return nil
}

// Close cleanly brings down the nats and jetstream connection.
func(n *NatsSubscription) Close() error {
	n.cctx.Stop()
	select {
	case <-n.cctx.Closed():
		n.conn.Close()
	}
	return nil
}

// jetstream message handler and acknowledger.
func(n *NatsSubscription) handleEvent(m jetstream.Msg) {
	var ev geEvent.Event

	logg.DebugCtxf(n.ctx, "have msg", "err", m)
	b := m.Data()
	err := json.Unmarshal(b, &ev)
	if err != nil {
		logg.ErrorCtxf(n.ctx, "nats msg deserialize fail", "err", err)
		//fail(m)
	} else {
		err = n.Route(n.ctx, &ev)
		if err != nil {
			logg.ErrorCtxf(n.ctx, "handler route fail", "err", err)
			//fail(m)
		}
	}
	err = m.Ack()
	if err != nil {
		logg.ErrorCtxf(n.ctx, "ack fail", "err", err)
		panic("ack fail")
	}
	logg.DebugCtxf(n.ctx, "handle msg complete")
}

// used if message should be retried.
func fail(m jetstream.Msg) {
	err := m.Nak()
	if err != nil {
		logg.Errorf("nats nak fail", "err", err)
	}
}

// server info string for debug.
func toServerInfo(conn *nats.Conn) string {
	return fmt.Sprintf("%s@%s (v%s)", conn.ConnectedServerName(), conn.ConnectedUrlRedacted(), conn.ConnectedServerVersion())
}

// on nats disconnection.
func disconnectHandler(conn *nats.Conn, err error) {
	logg.Errorf("nats disconnected", "status", conn.Status(), "reconnects", conn.Stats().Reconnects, "err", err)
}

// on nats reconnection.
func reconnectHandler(conn *nats.Conn) {
	serverInfo := toServerInfo(conn)
	logg.Errorf("nats reconnected", "status", conn.Status(), "reconnects", conn.Stats().Reconnects, "server", serverInfo)
}
