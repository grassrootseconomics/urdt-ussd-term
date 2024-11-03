package event

import (
	"git.grassecon.net/urdt/ussd/remote"
)

var (
	Api remote.AccountServiceInterface = &remote.AccountService{}
)
