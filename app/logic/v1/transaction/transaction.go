package transaction

import (
	"errors"

	"github.com/opensaucerer/barf/app/repository/v1/transaction"
)

// Transaction returns a list of transactions for the given account.
func Transaction(sessionId string) (*transaction.Transaction, error) {

	if sessionId == "" {
		return nil, errors.New("please provide a valid session id")
	}

	// find account
	tx := &transaction.Transaction{SessionId: sessionId}
	err := tx.FindBySessionId()
	if err != nil {
		return nil, errors.New("we are having issues finding your transaction. Please try again later")
	}

	if tx.Number == "" {
		return nil, errors.New("transaction not found")
	}

	return tx, nil
}
