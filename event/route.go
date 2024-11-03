package event

import (
	"context"
	"fmt"

	geEvent "github.com/grassrootseconomics/eth-tracker/pkg/event"

	"git.defalsify.org/vise.git/logging"
	"git.grassecon.net/urdt/ussd/common"
)

var (
	logg = logging.NewVanilla().WithDomain("term-event")
)

// Router is responsible for invoking handlers corresponding to events.
type Router struct {
	// User data store abstraction over application data.
	Store *common.UserDataStore
}

// Route parses an event from the event stream, and resolves the handler
// corresponding to the event.
//
// An error will be returned if no handler can be found, or if the resolved
// handler fails to successfully execute.
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
