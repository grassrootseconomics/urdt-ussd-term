package lookup

import (
	"context"

	"git.grassecon.net/urdt/ussd/common"
)

type Identity struct {
	NormalAddress string
	ChecksumAddress string
	SessionId string
}

func IdentityFromAddress(ctx context.Context, store common.UserDataStore, address string) (Identity, error) {
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

func getSessionIdByAddress(ctx context.Context, store common.UserDataStore, address string) (string, error) {
	
	r, err := store.ReadEntry(ctx, address, common.DATA_PUBLIC_KEY_REVERSE)
	if err != nil {
		return "", err
	}
	return string(r), nil
}
