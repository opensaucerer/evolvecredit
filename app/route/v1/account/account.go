package account

import (
	"github.com/opensaucerer/barf"
	accountc "github.com/opensaucerer/barf/app/controller/v1/account"
)

func RegisterAccountRoutes() {
	barf.Post("/v1/account/create", accountc.Create)
	barf.Get("/v1/account/search", accountc.Search)
	barf.Patch("/v1/account/deposit", accountc.Deposit)
	barf.Patch("/v1/account/lock", accountc.Lock)
	barf.Patch("/v1/account/unlock", accountc.Unlock)
	barf.Patch("/v1/account/withdraw", accountc.Withdraw)
	barf.Get("/v1/account/transactions", accountc.Transactions)
}
