package version

import (
	"github.com/opensaucerer/barf/app/route"
	"github.com/opensaucerer/barf/app/route/v1/account"
	"github.com/opensaucerer/barf/app/route/v1/user"
)

// Version1Routes registers all routes for the v1 version
func V1() {
	route.RegisterHomeRoutes()
	user.RegisterUserRoutes()
	account.RegisterAccountRoutes()
	// transaction.RegisterTransactionRoutes()
}
