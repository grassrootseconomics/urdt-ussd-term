package event

import (
	"context"
	"fmt"

	geEvent "github.com/grassrootseconomics/eth-tracker/pkg/event"

	"git.defalsify.org/vise.git/logging"
	"git.grassecon.net/urdt/ussd/common"
)

// TODO: this vocabulary should be public in and provided by the eth-tracker repo
const (
	evGive = "FAUCET_GIVE"
)

var (
	logg = logging.NewVanilla().WithDomain("term-event")
)

type Router struct {
	Store *common.UserDataStore
}

func(r *Router) Route(ctx context.Context, gev *geEvent.Event) error {
	logg.DebugCtxf(ctx, "have event", "ev", gev)
	evCC, ok := asCustodialRegistrationEvent(gev)
	if ok {
		return handleCustodialRegistration(ctx, r.Store, evCC)
	}
	evTT, ok := asTokenTransferEvent(gev)
	if ok {
		return handleTokenTransfer(ctx, r.Store, evTT)
	}

	return fmt.Errorf("unexpected message")
}
