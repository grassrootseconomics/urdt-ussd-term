package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"git.grassecon.net/urdt/ussd/initializers"
	"git.grassecon.net/urdt/ussd/common"
	"git.grassecon.net/term/config"
	"git.grassecon.net/term/event/nats"
)

func init() {
	initializers.LoadEnvVariables()
}

func main() {
	config.LoadConfig()

	var dbDir string
	flag.StringVar(&dbDir, "dbdir", ".state", "database dir to read from")
	flag.Parse()

	ctx := context.Background()
//	db := mem.NewMemDb()
//	err := db.Connect(ctx, "")
//	if err != nil {
//		fmt.Fprintf(os.Stderr, "Db connect err: %v", err)
//		os.Exit(1)
//	}

	menuStorageService := common.NewStorageService(dbDir)
	err := menuStorageService.EnsureDbDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}

	n := nats.NewNatsSubscription(menuStorageService)
	err = n.Connect(ctx, config.JetstreamURL)
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
