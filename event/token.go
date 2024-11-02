package event

import (
	"context"
	"fmt"
	"strings"
	"strconv"

	geEvent "github.com/grassrootseconomics/eth-tracker/pkg/event"

	"git.defalsify.org/vise.git/db"
	"git.grassecon.net/urdt/ussd/common"
	"git.grassecon.net/urdt/ussd/remote"
	"git.grassecon.net/term/lookup"
)

const (
	evTokenTransfer = "TOKEN_TRANSFER"
	// TODO: use export from urdt storage
	DATATYPE_USERSUB = 64
)

func renderTx() {

}

type eventTokenTransfer struct {
	From string
	To string
	Value int
	TxHash string
}

func updateTokenTransferList(ctx context.Context, api remote.AccountServiceInterface, store common.UserDataStore, identity lookup.Identity) error {
	var r []string

	txs, err := api.FetchTransactions(ctx, identity.ChecksumAddress)
	if err != nil {
		return err
	}

	for i, tx := range(txs) {
		r = append(r, fmt.Sprintf("%d %s %s", i, tx.DateBlock, tx.TxHash[:10]))
	}

	s := strings.Join(r, "\n")
	return store.WriteEntry(ctx, identity.SessionId, common.DATA_TRANSACTIONS, []byte(s))
}

func updateTokenList(ctx context.Context, api remote.AccountServiceInterface, store *common.UserDataStore, identity lookup.Identity) error {
	holdings, err := api.FetchVouchers(ctx, identity.ChecksumAddress)
	if err != nil {
		return err
	}
	metadata := common.ProcessVouchers(holdings)
	_ = metadata

	// TODO: export subprefixdb and use that instead
	// TODO: make sure subprefixdb is thread safe when using gdbm
	store.Db.SetPrefix(DATATYPE_USERSUB)

	k := append([]byte("vouchers"), []byte("sym")...)
	err = store.Db.Put(ctx, k, []byte(metadata.Symbols))
	if err != nil {
		return err
	}

	k = append([]byte("vouchers"), []byte("bal")...)
	err = store.Db.Put(ctx, k, []byte(metadata.Balances))
	if err != nil {
		return err
	}

	k = append([]byte("vouchers"), []byte("deci")...)
	err = store.Db.Put(ctx, k, []byte(metadata.Decimals))
	if err != nil {
		return err
	}

	k = append([]byte("vouchers"), []byte("addr")...)
	err = store.Db.Put(ctx, k, []byte(metadata.Addresses))
	if err != nil {
		return err
	}

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
func updateDefaultToken(ctx context.Context, store *common.UserDataStore, identity lookup.Identity, activeSym string) error {
	return nil
}

func updateWait(ctx context.Context, api remote.AccountServiceInterface) error {
	return nil
}

func updateToken(ctx context.Context, store *common.UserDataStore, identity lookup.Identity) error {
	var api remote.AccountService

	err := updateTokenList(ctx, &api, store, identity)
	if err != nil {
		return err
	}

	activeSym, err := store.ReadEntry(ctx, identity.SessionId, common.DATA_ACTIVE_ADDRESS)
	if err == nil {
		return nil
	}
	if !db.IsNotFound(err) {
		return err
	}

	err = updateDefaultToken(ctx, store, identity, string(activeSym))
	if err != nil {
		return err
	}

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

func asTokenTransferEvent(gev *geEvent.Event) (*eventTokenTransfer, bool) {
	var err error
	var ok bool
	var ev eventTokenTransfer

	if gev.TxType != evTokenTransfer {
		return nil, false
	}

	pl := gev.Payload
	// assuming from and to are checksum addresses
	ev.From, ok = pl["from"].(string)
	if !ok {
		return nil, false
	}
	ev.To, ok = pl["to"].(string)
	if !ok {
		return nil, false
	}
	ev.TxHash, err = common.NormalizeHex(gev.TxHash)
	if err != nil {
		logg.Error("could not decode tx hash", "tx", gev.TxHash, "err", err)
		return nil, false
	}

	value, ok := pl["value"].(string)
	if !ok {
		return nil, false
	}
	ev.Value, err = strconv.Atoi(value)
	if err != nil {
		logg.Error("could not decode value", "value", value, "err", err)
		return nil, false
	}
	return &ev, true
}

func handleTokenTransfer(ctx context.Context, store *common.UserDataStore, ev *eventTokenTransfer) error {
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
