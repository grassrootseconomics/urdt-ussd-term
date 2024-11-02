package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"git.defalsify.org/vise.git/db/mem"

	"git.grassecon.net/urdt/ussd/config"
	"git.grassecon.net/urdt/ussd/initializers"
	"git.grassecon.net/term/event/nats"
)

func init() {
	initializers.LoadEnvVariables()
}

func main() {
	config.LoadConfig()

	ctx := context.Background()
	db := mem.NewMemDb()
	err := db.Connect(ctx, "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Db connect err: %v", err)
		os.Exit(1)
	}
	n := nats.NewNatsSubscription(db)
	err = n.Connect(ctx, "localhost:4222")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Stream connect err: %v", err)
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
