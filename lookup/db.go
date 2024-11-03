package lookup

import (
	"context"

	"git.defalsify.org/vise.git/db"
	"git.defalsify.org/vise.git/logging"
	"git.grassecon.net/urdt/ussd/common"
)

var (
	logg = logging.NewVanilla().WithDomain("term-lookup")
)

type Identity struct {
	NormalAddress string
	ChecksumAddress string
	SessionId string
}

func IdentityFromAddress(ctx context.Context, store *common.UserDataStore, address string) (Identity, error) {
	var err error
	var ident Identity

	ident.ChecksumAddress = address
	ident.NormalAddress, err = common.NormalizeHex(ident.ChecksumAddress)
	if err != nil {
		return ident, err
	}
	ident.SessionId, err = getSessionIdByAddress(ctx, store, ident.NormalAddress)
	if err != nil {
		return ident, err
	}
	return ident, nil
}

func getSessionIdByAddress(ctx context.Context, store *common.UserDataStore, address string) (string, error) {
	logg.Debugf("fooar")
	// TODO: replace with userdatastore when double sessionid issue fixed
	//r, err := store.ReadEntry(ctx, address, common.DATA_PUBLIC_KEY_REVERSE)
	store.Db.SetPrefix(db.DATATYPE_USERDATA)
	store.Db.SetSession(address)
	r, err := store.Db.Get(ctx, common.PackKey(common.DATA_PUBLIC_KEY_REVERSE, []byte{}))
	if err != nil {
		return "", err
	}
	return string(r), nil
}
