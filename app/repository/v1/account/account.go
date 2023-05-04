package account

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/opensaucerer/barf/app/database"
	"github.com/opensaucerer/barf/app/reflection"
)

// Fields returns the struct fields as a slice of interface{} values
func (a *Account) Fields() []interface{} {
	return reflection.ReturnStructFields(a)
}

// Create inserts a new user into the database.
func (a *Account) Create() error {

	a.time(true)

	_, err := database.PostgreSQLDB.Exec(context.Background(), `INSERT INTO accounts (owner, type, number, locked_balance, ledger_balance, balance, active, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`, a.Owner, a.Type, a.Number, a.LockedBalance, a.LedgerBalance, a.Balance, a.Active, a.CreatedAt, a.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

// time updates the CreatedAt and UpdatedAt fields on the struct
func (a *Account) time(new ...bool) {
	if len(new) > 0 && new[0] {
		a.CreatedAt = time.Now().UTC()
	}
	a.UpdatedAt = time.Now().UTC()
}

// FindByNumber finds the account by the number field
func (a *Account) FindByNumber() error {
	err := database.PostgreSQLDB.QueryRow(context.Background(), `SELECT "accounts".id, "accounts".type, number, locked_balance, ledger_balance, balance, "accounts".active, "accounts".created_at, "accounts".updated_at, u.id, u.key, u.first_name, u.last_name, u.email, u.age, u.role, u.active, u.created_at, u.updated_at FROM accounts LEFT JOIN users as u ON owner = u.id WHERE number = $1`, a.Number).Scan(a.Fields()...)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil
		}
		return err
	}
	return nil
}

// FindByOwner finds all accounts by the owner field
func (a Accounts) FindByOwner(owner string) error {
	rows, err := database.PostgreSQLDB.Query(context.Background(), `SELECT "accounts".id, "accounts".type, number, locked_balance, ledger_balance, balance, "accounts".active, "accounts".created_at, "accounts".updated_at, u.id, u.key, u.first_name, u.last_name, u.email, u.age, u.role, u.active, u.created_at, u.updated_at FROM accounts LEFT JOIN users as u ON owner = u.id WHERE owner = $1`, owner)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var account Account
		err := rows.Scan(account.Fields()...)
		if err != nil {
			return err
		}
		a = append(a, account)
	}
	return nil
}

// Delete deletes an account from the database. This is only used for testing.
func (a *Account) Delete() error {
	_, err := database.PostgreSQLDB.Exec(context.Background(), `DELETE FROM accounts WHERE number = $1`, a.Number)
	if err != nil {
		return err
	}
	return nil
}

// Deposit adds the amount to the account balance and ledger balance atomically and
// transactionally. This means a database.PostgreSQLDBTx must have been started before
func (a *Account) Deposit(amount float64) error {
	a.time()
	_, err := database.PostgreSQLDBTx.Exec(context.Background(), `UPDATE accounts SET balance = balance + $1, ledger_balance = ledger_balance + $1, updated_at = $2 WHERE number = $3`, amount, a.UpdatedAt, a.Number)
	if err != nil {
		return err
	}
	return nil
}

// Lock locks the amount on the account. This means the amount is subtracted from the balance and added to the locked balance but only if the balance is sufficient. This is done atomically and transactionally. This means a database.PostgreSQLDBTx must have been started before
func (a *Account) Lock(amount float64) error {
	a.time()
	err := database.PostgreSQLDBTx.QueryRow(context.Background(), `UPDATE accounts SET balance = balance - $1, locked_balance = locked_balance + $1, updated_at = $2 WHERE number = $3 AND balance >= $1 RETURNING id`, amount, a.UpdatedAt, a.Number).Scan(&a.Id)
	if err != nil {
		return err
	}
	return nil
}

// Unlock unlocks the amount on the account. This means the amount is subtracted from the locked balance and added to the balance but only if the locked balance is sufficient. This is done atomically and transactionally. This means a database.PostgreSQLDBTx must have been started before
func (a *Account) Unlock(amount float64) error {
	a.time()
	err := database.PostgreSQLDBTx.QueryRow(context.Background(), `UPDATE accounts SET balance = balance + $1, locked_balance = locked_balance - $1, updated_at = $2 WHERE number = $3 AND locked_balance >= $1 RETURNING id`, amount, a.UpdatedAt, a.Number).Scan(&a.Id)
	if err != nil {
		return err
	}
	return nil
}

// Withdraw withdraws the amount from the account balance and ledger balance but only if the balance is sufficient. This is done atomically and transactionally. This means a database.PostgreSQLDBTx must have been started before
func (a *Account) Withdraw(amount float64) error {
	a.time()
	err := database.PostgreSQLDBTx.QueryRow(context.Background(), `UPDATE accounts SET balance = balance - $1, ledger_balance = ledger_balance - $1, updated_at = $2 WHERE number = $3 AND balance >= $1 RETURNING id`, amount, a.UpdatedAt, a.Number).Scan(&a.Id)
	if err != nil {
		return err
	}
	return nil
}
