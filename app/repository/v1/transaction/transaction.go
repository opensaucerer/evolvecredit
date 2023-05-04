package transaction

import (
	"context"
	"crypto/sha1"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/opensaucerer/barf/app/database"
	"github.com/opensaucerer/barf/app/reflection"
)

// Fields returns the struct fields as a slice of interface{} values
func (t *Transaction) Fields() []interface{} {
	return reflection.ReturnStructFields(t)
}

// Create inserts a new transaction into the database transactionally. This means a database.PostgreSQLDBTx must have been started before else the function will panic on a nil pointer dereference.
func (t *Transaction) Create() error {

	t.time(true)

	if err := t.session(); err != nil {
		return err
	}

	_, err := database.PostgreSQLDBTx.Exec(context.Background(), `INSERT INTO transactions (number, amount, session_id, type, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`, t.Number, t.Amount, t.SessionId, t.Type, t.Status, t.CreatedAt, t.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

// time updates the CreatedAt and UpdatedAt fields on the struct
func (t *Transaction) time(new ...bool) {
	if len(new) > 0 && new[0] {
		t.CreatedAt = time.Now().UTC()
	}
	t.UpdatedAt = time.Now().UTC()
}

// Session add a session id to the struct if it is not already set
func (t *Transaction) session() error {
	if t.SessionId == "" {
		salt := strconv.FormatInt(time.Now().UnixNano(), 10) + t.Number + strconv.FormatFloat(t.Amount, 'f', 2, 64) + strconv.Itoa(int(t.Type)) + strconv.Itoa(int(t.Status)) + strconv.FormatInt(t.CreatedAt.UnixNano(), 10) + strconv.FormatInt(t.UpdatedAt.UnixNano(), 10)
		hash := sha1.New()
		_, err := hash.Write([]byte(salt))
		if err != nil {
			return err
		}
		t.SessionId = strings.ToUpper(fmt.Sprintf("%x", hash.Sum(nil)))
	}
	return nil
}

// FindByAccountNumber finds all transactions by the account field
func (t *Transactions) FindByAccountNumber(number string) error {
	rows, err := database.PostgreSQLDB.Query(context.Background(), `SELECT "transactions".id, "transactions".number, amount, session_id, "transactions".type, status, "transactions".created_at, "transactions".updated_at, a.id, a.type, a.number, a.locked_balance, a.ledger_balance, a.balance, a.active, a.created_at, a.updated_at, u.id, u.key, u.first_name, u.last_name, u.email, u.age, u.role, u.active, u.created_at, u.updated_at FROM transactions LEFT JOIN accounts as a ON "transactions".number = a.number LEFT JOIN users as u ON a.owner = u.id WHERE "transactions".number = $1`, number)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var transaction Transaction
		err := rows.Scan(transaction.Fields()...)
		if err != nil {
			return err
		}
		*t = append(*t, transaction)
	}
	return nil
}

// FindBySessionId find a transaction by the session_id field
func (t *Transaction) FindBySessionId() error {
	err := database.PostgreSQLDB.QueryRow(context.Background(), `SELECT "transactions".id, "transactions".number, amount, session_id, "transactions".type, status, "transactions".created_at, "transactions".updated_at, a.id, a.type, a.number, a.locked_balance, a.ledger_balance, a.balance, a.active, a.created_at, a.updated_at, u.id, u.key, u.first_name, u.last_name, u.email, u.age, u.role, u.active, u.created_at, u.updated_at FROM transactions LEFT JOIN accounts as a ON "transactions".number = a.number LEFT JOIN users as u ON a.owner = u.id WHERE session_id = $1`, t.SessionId).Scan(t.Fields()...)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil
		}
		return err
	}
	return nil
}

// Delete deletes a transaction from the database. This is only used for testing.
func (t *Transaction) Delete() error {
	_, err := database.PostgreSQLDB.Exec(context.Background(), `DELETE FROM transactions WHERE session_id = $1`, t.SessionId)
	if err != nil {
		return err
	}
	return nil
}
