package transaction

import (
	"time"

	"github.com/opensaucerer/barf/app/global"
	"github.com/opensaucerer/barf/app/repository/v1/account"
)

type Transaction struct {
	Id        int64         `json:"-"`
	Number    string        `json:"number"`
	Amount    float64       `json:"amount"`
	SessionId string        `json:"session_id"`
	Type      global.Type   `json:"type"`
	Status    global.Status `json:"status"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
	// DeletedAt string    `json:"deleted_at"`
	Account account.Account `json:"account"`
}

type Transactions []Transaction
