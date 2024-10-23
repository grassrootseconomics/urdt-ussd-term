package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"git.grassecon.net/term/event/nats"
)

func main() {
	ctx := context.Background()
	n := nats.NewNatsSubscription()
	err := n.Connect(ctx, "localhost:4222")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Connect err: %v", err)
		os.Exit(1)
	}
	defer n.Close()

	cint := make(chan os.Signal)
	cterm := make(chan os.Signal)
	signal.Notify(cint, os.Interrupt, syscall.SIGINT)
	signal.Notify(cterm, os.Interrupt, syscall.SIGTERM)
	select {
	case _ = <-cint:
	case _ = <-cterm:
	}
}
