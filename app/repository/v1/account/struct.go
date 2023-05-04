package account

import (
	"time"

	"github.com/opensaucerer/barf/app/global"
	"github.com/opensaucerer/barf/app/repository/v1/user"
)

type Account struct {
	Id            int64              `json:"-"`
	Owner         int64              `json:"-" rsf:"false"` // map to the id of the user
	Type          global.AccountType `json:"type"`
	Number        string             `json:"number"`
	LockedBalance float64            `json:"locked_balance"`
	LedgerBalance float64            `json:"ledger_balance"` // locked + available
	Balance       float64            `json:"balance"`        // available
	Active        bool               `json:"active"`
	CreatedAt     time.Time          `json:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at"`
	// DeletedAt string    `json:"deleted_at"`
	User user.User `json:"user"`
}

type Accounts []Account
