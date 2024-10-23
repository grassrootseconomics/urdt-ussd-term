package event

import (
	geEvent "github.com/grassrootseconomics/eth-tracker/pkg/event"
)

type Router struct {
}

func(r *Router) Route(ev *geEvent.Event) error {
	return nil
}
