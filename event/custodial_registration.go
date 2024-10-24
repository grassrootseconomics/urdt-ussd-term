package event

import (
	geEvent "github.com/grassrootseconomics/eth-tracker/pkg/event"
)

const (
	evReg = "CUSTODIAL_REGISTRATION"
)

type eventCustodialRegistration struct {
	account string
}

func asCustodialRegistrationEvent(gev *geEvent.Event) (*eventCustodialRegistration, bool) {
	var ok bool
	var ev eventCustodialRegistration
	if gev.TxType != evReg {
		return nil, false
	}
	pl := gev.Payload
	ev.account, ok = pl["account"].(string)
	if !ok {
		return nil, false
	}
	logg.Debug("parsed ev", "pl", pl, "ev", ev)
	return &ev, true
}

func handleCustodialRegistration(ev *eventCustodialRegistration) error {
	return nil
}
