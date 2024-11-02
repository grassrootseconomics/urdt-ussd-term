package event

import (
	"context"

	"git.defalsify.org/vise.git/db"
	"git.grassecon.net/urdt/ussd/common"
	"git.grassecon.net/urdt/ussd/remote"
	"git.grassecon.net/term/lookup"
)

const (
	evTransfer = "TOKEN_TRANSFER"
)

type eventTokenTransfer struct {
	From string
	To string
	Value string
}

//func updateTokenTransferList(ctx context.Context, api remote.AccountServiceInterface, store common.UserDataStore, sessionId string) error {
//	return nil
//}

func updateTokenList(ctx context.Context, api remote.AccountServiceInterface, store common.UserDataStore, identity lookup.Identity) error {
	api.FetchVouchers(ctx, identity.ChecksumAddress)
	return nil
}

//func updateTokenBalance(ctx context.Context, api remote.AccountServiceInterface, store common.UserDataStore, sessionId string) error {
//	r, err := api.CheckBalance(ctx, sessionId)
//	if err != nil {
//		return err
//	}
//	//store.WriteEntry()
//	return nil
//}
//
//func updateDefaultToken(ctx context.Context, store common.UserDataStore, sessionId string, activeSym string) {
//
//}

func updateToken(ctx context.Context, store common.UserDataStore, identity lookup.Identity) error {
	var api remote.AccountService

	err := updateTokenList(ctx, &api, store, identity)
	if err != nil {
		return err
	}

//	activeSym, err := store.ReadEntry(common.DATA_ACTIVE_ADDRESS)
//	if err == nil {
//		return nil
//	}
//	if !db.IsNotFound(err) {
//		return err
//	}
//
//	err = updateDefaultToken(ctx, store, sessionId, string(activeSym))
//	if err != nil {
//		return err
//	}
//	err = updateTokenBalance(ctx, &api, store, sessionId)
//	if err != nil {
//		return err
//	}
//	err = updateTokenTransferList(ctx, &api, store, sessionId)
//	if err != nil {
//		return err
//	}
//	
	return nil
}

func handleTokenTransfer(ctx context.Context, store common.UserDataStore, ev *eventTokenTransfer) error {
	identity, err := lookup.IdentityFromAddress(ctx, store, ev.From)
	if err != nil {
		if !db.IsNotFound(err) {
			return err
		}
	} else {
		err = updateToken(ctx, store, identity)
		if err != nil {
			return err
		}
	}
	identity, err = lookup.IdentityFromAddress(ctx, store, ev.To)
	if err != nil {
		if !db.IsNotFound(err) {
			return err
		}
	} else {
		err = updateToken(ctx, store, identity)
		if err != nil {
			return err
		}
	}

	return nil
}
