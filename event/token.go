package event

import (
	"context"
	"fmt"
	"strings"
	"strconv"

	geEvent "github.com/grassrootseconomics/eth-tracker/pkg/event"

	"git.defalsify.org/vise.git/db"
	"git.grassecon.net/urdt/ussd/common"
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
	VoucherAddress string
}

func updateTokenTransferList(ctx context.Context, store common.UserDataStore, identity lookup.Identity) error {
	var r []string

	txs, err := lookup.Api.FetchTransactions(ctx, identity.ChecksumAddress)
	if err != nil {
		return err
	}

	for i, tx := range(txs) {
		r = append(r, fmt.Sprintf("%d %s %s", i, tx.DateBlock, tx.TxHash[:10]))
	}

	s := strings.Join(r, "\n")
	return store.WriteEntry(ctx, identity.SessionId, common.DATA_TRANSACTIONS, []byte(s))
}

func updateTokenList(ctx context.Context, store *common.UserDataStore, identity lookup.Identity) error {
	holdings, err := lookup.Api.FetchVouchers(ctx, identity.ChecksumAddress)
	if err != nil {
		return err
	}
	metadata := common.ProcessVouchers(holdings)
	_ = metadata

	// TODO: export subprefixdb and use that instead
	// TODO: make sure subprefixdb is thread safe when using gdbm
	// TODO: why is address session here unless explicitly set
	store.Db.SetSession(identity.SessionId)
	store.Db.SetPrefix(DATATYPE_USERSUB)

	k := append([]byte("vouchers"), []byte("sym")...)
	err = store.Db.Put(ctx, k, []byte(metadata.Symbols))
	if err != nil {
		return err
	}
	logg.TraceCtxf(ctx, "processvoucher", "key", k)
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

func updateDefaultToken(ctx context.Context, store *common.UserDataStore, identity lookup.Identity, activeSym string) error {
	pfxDb := common.StoreToPrefixDb(store, []byte("vouchers"))
	// TODO: the activeSym input should instead be newline separated list?
	tokenData, err := common.GetVoucherData(ctx, pfxDb, activeSym)
	if err != nil {
		return err
	}
	logg.TraceCtxf(ctx, "tokendaa", "d", tokenData)
	return common.UpdateVoucherData(ctx, store, identity.SessionId, tokenData)
}

func updateWait(ctx context.Context) error {
	return nil
}

func toSym(ctx context.Context, address string) ([]byte, error) {
	voucherData, err := lookup.Api.VoucherData(ctx, address)
	if err != nil {
		return nil, err
	}
	return []byte(voucherData.TokenSymbol), nil
}

func updateToken(ctx context.Context, store *common.UserDataStore, identity lookup.Identity, tokenAddress string) error {
	err := updateTokenList(ctx, store, identity)
	if err != nil {
		return err
	}

	store.Db.SetSession(identity.SessionId)
	activeSym, err := store.ReadEntry(ctx, identity.SessionId, common.DATA_ACTIVE_SYM)
	if err == nil {
		return nil
	}
	if !db.IsNotFound(err) {
		return err
	}
	if activeSym == nil {
		activeSym, err = toSym(ctx, tokenAddress)
		if err != nil {
			return err
		}
	}
	logg.Debugf("barfoo")

	err = updateDefaultToken(ctx, store, identity, string(activeSym))
	if err != nil {
		return err
	}

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
		logg.Errorf("could not decode tx hash", "tx", gev.TxHash, "err", err)
		return nil, false
	}

	value, ok := pl["value"].(string)
	if !ok {
		return nil, false
	}
	ev.Value, err = strconv.Atoi(value)
	if err != nil {
		logg.Errorf("could not decode value", "value", value, "err", err)
		return nil, false
	}

	ev.VoucherAddress = gev.ContractAddress

	return &ev, true
}

func handleTokenTransfer(ctx context.Context, store *common.UserDataStore, ev *eventTokenTransfer) error {
	identity, err := lookup.IdentityFromAddress(ctx, store, ev.From)
	if err != nil {
		if !db.IsNotFound(err) {
			return err
		}
	} else {
		err = updateToken(ctx, store, identity, ev.VoucherAddress)
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
		err = updateToken(ctx, store, identity, ev.VoucherAddress)
		if err != nil {
			return err
		}
	}

	return nil
}
