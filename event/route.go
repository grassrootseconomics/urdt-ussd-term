package event

import (
	"fmt"
	"log/slog"
	"os"

	geEvent "github.com/grassrootseconomics/eth-tracker/pkg/event"
)

// TODO: this vocabulary should be public in and provided by the eth-tracker repo
const (
	evGive = "FAUCET_GIVE"
)

var (
	logg = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))
)

type Router struct {
}

func(r *Router) Route(gev *geEvent.Event) error {
	logg.Debug("have event", "ev", gev)
	evCC, ok := asCustodialRegistrationEvent(gev)
	if ok {
		return handleCustodialRegistration(evCC)
	}
	return fmt.Errorf("unexpected message")
}
