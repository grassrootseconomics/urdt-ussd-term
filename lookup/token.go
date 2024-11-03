package lookup

import (
	"git.grassecon.net/urdt/ussd/remote"
)

var (
	// Api provides the api implementation for all external lookups.
	Api remote.AccountServiceInterface = &remote.AccountService{}
)
