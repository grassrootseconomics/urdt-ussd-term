package handler

import (
	"git.defalsify.org/vise.git/db"

	"git.grassecon.net/term/common"
	"git.grassecon.net/term/lookup"
)

func HandleEvReg(store db.Db, payload map[string]any) error {
	var err error
	address, ok := payload["address"].(string)
	if !ok {
		return ErrInvalidPayload
	}
	address, err = common.NormalizeHex(address)
	if err != nil {
		return err
	}
	sessionId, err := lookup.GetSessionIdByAddress(store, address)
	if err != nil {
		return err
	}
	_ = sessionId
	return nil
}
