package transaction

import (
	"github.com/opensaucerer/barf"
	"github.com/opensaucerer/barf/app/controller/v1/transaction"
)

func RegisterTransactionRoutes() {
	barf.Get("/v1/transaction", transaction.Transaction)
}
