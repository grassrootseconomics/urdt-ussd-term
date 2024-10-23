package nats

import (
	"context"
	"encoding/json"

	nats "github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"

	"git.defalsify.org/vise.git/logging"

	geEvent "github.com/grassrootseconomics/eth-tracker/pkg/event"

	"git.grassecon.net/term/event"
)

const (
	evGive = "TRACKER.FAUCET_GIVE"
	evReg = "TRACKER.CUSTODIAL_REGISTRATION"
)

var (
	logg = logging.NewVanilla().WithDomain("nats")
)

type NatsSubscription struct {
	*event.Router
	ctx context.Context
	conn *nats.Conn
	js jetstream.JetStream
	cs jetstream.Consumer
	cctx jetstream.ConsumeContext
}

func NewNatsSubscription() *NatsSubscription {
	return &NatsSubscription{}
}

func(n *NatsSubscription) Connect(ctx context.Context, connStr string) error {
	var err error

	n.conn, err = nats.Connect(connStr)
	if err != nil {
		return err
	}
	n.js, err = jetstream.New(n.conn)
	if err != nil {
		return err
	}
	n.cs, err = n.js.OrderedConsumer(ctx, "TRACKER", jetstream.OrderedConsumerConfig{
		//FilterSubjects: []string{"TRACKER.*"},
	})
	if err != nil {
		return err
	}

	n.ctx = ctx
	n.cctx, err = n.cs.Consume(n.handleEvent)
	if err != nil {
		return err		
	}

	return nil
}

func(n *NatsSubscription) Close() error {
	n.cctx.Stop()
	select {
	case <-n.cctx.Closed():
		n.conn.Close()
	}
	return nil
}

func fail(m jetstream.Msg) {
	err := m.Nak()
	if err != nil {
		logg.Errorf("nats nak fail", "err", err)
	}
}

func(n *NatsSubscription) handleEvent(m jetstream.Msg) {
	var ev geEvent.Event

	logg.Tracef("have msg", "err", m)
	b := m.Data()
	err := json.Unmarshal(b, &ev)
	if err != nil {
		logg.Errorf("nats msg deserialize fail", "err", err)
		fail(m)
		return
	}
	err = n.Route(&ev)
	if err != nil {
		logg.Errorf("handler route fail", "err", err)
		fail(m)
		return
	}
	err = m.Ack()
//	err = m.DoubleAck(n.ctx)
	if err != nil {
		logg.Errorf("ack fail", "err", err)
		panic("ack fail")
	}
}
