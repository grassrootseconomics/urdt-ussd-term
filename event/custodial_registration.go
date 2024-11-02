package event

import (
	"context"

	geEvent "github.com/grassrootseconomics/eth-tracker/pkg/event"

	"git.grassecon.net/urdt/ussd/common"
	"git.grassecon.net/term/lookup"
)

const (
	evReg = "CUSTODIAL_REGISTRATION"
)

type eventCustodialRegistration struct {
	Account string
}

func asCustodialRegistrationEvent(gev *geEvent.Event) (*eventCustodialRegistration, bool) {
	var ok bool
	var ev eventCustodialRegistration
	if gev.TxType != evReg {
		return nil, false
	}
	pl := gev.Payload
	ev.Account, ok = pl["account"].(string)
	if !ok {
		return nil, false
	}
	logg.Debug("parsed ev", "pl", pl, "ev", ev)
	return &ev, true
}

func handleCustodialRegistration(ctx context.Context, store *common.UserDataStore, ev *eventCustodialRegistration) error {
	identity, err := lookup.IdentityFromAddress(ctx, store, ev.Account)
	if err != nil {
		return err
	}
	_ = identity
	return nil
}
