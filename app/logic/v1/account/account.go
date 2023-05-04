package account

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/opensaucerer/barf/app/database"
	"github.com/opensaucerer/barf/app/global"
	"github.com/opensaucerer/barf/app/repository"
	accountr "github.com/opensaucerer/barf/app/repository/v1/account"
	"github.com/opensaucerer/barf/app/repository/v1/transaction"
	userr "github.com/opensaucerer/barf/app/repository/v1/user"
)

// Create adds a new account for the given user.
func Create(user *userr.User) (*accountr.Account, error) {

	// again, this is not good enough and should be improved with a validation middleware
	if user.Key == "" {
		return nil, errors.New("please provide a valid user")
	}

	// validate user's existence
	user.FindByKey()

	if user.Email == "" {
		return nil, errors.New("user not found")
	}

	number, err := repository.GenerateAccountNumber()
	if err != nil {
		return nil, errors.New("we are having issues creating your account number. Please try again later")
	}

	// create account
	account := accountr.Account{
		Owner:  user.Id,
		Type:   global.Savings,
		Number: number,
		Active: true,
		User:   *user,
	}

	if err := account.Create(); err != nil {
		return nil, errors.New("we are having issues creating your account. Please try again later")
	}

	return &account, nil
}

// Search returns an account for the given account number.
func Search(number string) (*accountr.Account, error) {

	if number == "" {
		return nil, errors.New("please provide a valid account number")
	}

	account := accountr.Account{
		Number: number,
	}

	if err := account.FindByNumber(); err != nil {
		return nil, errors.New("we are having issues finding your account. Please try again later")
	}

	if account.User.Key == "" {
		return nil, errors.New("account not found")
	}

	return &account, nil
}

// Deposit adds the given amount to the account's balance and records the transaction.
func Deposit(tx *transaction.Transaction) (*transaction.Transaction, error) {

	if tx.Number == "" {
		return nil, errors.New("please provide a valid account number")
	}

	// find account
	tx.Account.Number = tx.Number
	err := tx.Account.FindByNumber()
	if err != nil {
		return nil, errors.New("we are having issues finding your account. Please try again later")
	}

	if tx.Account.User.Key == "" {
		return nil, errors.New("account not found")
	}

	// prepare the deposit transaction
	tx.Type = global.Deposit
	tx.Status = global.Completed

	// update account balance
	database.PostgreSQLDBTx, err = database.PostgreSQLDB.Begin(context.Background())
	if err != nil {
		return nil, errors.New("we are having issues processing your deposit. Please try again later")
	}

	if err := tx.Account.Deposit(tx.Amount); err != nil {
		database.PostgreSQLDBTx.Rollback(context.Background())
		return nil, errors.New("we are having issues processing your deposit. Please try again later")
	}

	// create transaction
	if err := tx.Create(); err != nil {
		database.PostgreSQLDBTx.Rollback(context.Background())
		return nil, errors.New("we are having issues processing your deposit. Please try again later")
	}

	// commit transaction
	if err := database.PostgreSQLDBTx.Commit(context.Background()); err != nil {
		database.PostgreSQLDBTx.Rollback(context.Background())
		return nil, errors.New("we are having issues processing your deposit. Please try again later")
	}

	// dispose of the transaction
	database.PostgreSQLDBTx = nil

	return tx, nil
}

// Lock moves the given amount from the account's balance to the locked balance and records the transaction.
func Lock(tx *transaction.Transaction) (*transaction.Transaction, error) {

	if tx.Number == "" {
		return nil, errors.New("please provide a valid account number")
	}

	// find account
	tx.Account.Number = tx.Number
	err := tx.Account.FindByNumber()
	if err != nil {
		return nil, errors.New("we are having issues finding your account. Please try again later")
	}

	if tx.Account.User.Key == "" {
		return nil, errors.New("account not found")
	}

	// prepare the lock transaction
	tx.Type = global.Lock
	tx.Status = global.Completed

	// update account balance
	database.PostgreSQLDBTx, err = database.PostgreSQLDB.Begin(context.Background())
	if err != nil {
		return nil, errors.New("we are having issues processing your lock. Please try again later")
	}

	if err := tx.Account.Lock(tx.Amount); err != nil {
		database.PostgreSQLDBTx.Rollback(context.Background())
		if err == pgx.ErrNoRows {
			return nil, errors.New("insufficient funds in account's available balance")
		}
		return nil, errors.New("we are having issues processing your lock. Please try again later")
	}

	// create transaction
	if err := tx.Create(); err != nil {
		database.PostgreSQLDBTx.Rollback(context.Background())
		return nil, errors.New("we are having issues processing your lock. Please try again later")
	}

	// commit transaction
	if err := database.PostgreSQLDBTx.Commit(context.Background()); err != nil {
		database.PostgreSQLDBTx.Rollback(context.Background())
		return nil, errors.New("we are having issues processing your lock. Please try again later")
	}

	// dispose of the transaction
	database.PostgreSQLDBTx = nil

	return tx, nil
}

// Unlock moves the given amount from the account's locked balance to the balance and records the transaction.
func Unlock(tx *transaction.Transaction) (*transaction.Transaction, error) {

	if tx.Number == "" {
		return nil, errors.New("please provide a valid account number")
	}

	// find account
	tx.Account.Number = tx.Number
	err := tx.Account.FindByNumber()
	if err != nil {
		return nil, errors.New("we are having issues finding your account. Please try again later")
	}

	if tx.Account.User.Key == "" {
		return nil, errors.New("account not found")
	}

	// prepare the unlock transaction
	tx.Type = global.Unlock
	tx.Status = global.Completed

	// update account balance
	database.PostgreSQLDBTx, err = database.PostgreSQLDB.Begin(context.Background())
	if err != nil {
		return nil, errors.New("we are having issues processing your unlock. Please try again later")
	}

	if err := tx.Account.Unlock(tx.Amount); err != nil {
		database.PostgreSQLDBTx.Rollback(context.Background())
		if err == pgx.ErrNoRows {
			return nil, errors.New("insufficient funds in account's locked balance")
		}
		return nil, errors.New("we are having issues processing your unlock. Please try again later")
	}

	// create transaction
	if err := tx.Create(); err != nil {
		database.PostgreSQLDBTx.Rollback(context.Background())
		return nil, errors.New("we are having issues processing your unlock. Please try again later")
	}

	// commit transaction
	if err := database.PostgreSQLDBTx.Commit(context.Background()); err != nil {
		database.PostgreSQLDBTx.Rollback(context.Background())
		return nil, errors.New("we are having issues processing your unlock. Please try again later")
	}

	// dispose of the transaction
	database.PostgreSQLDBTx = nil

	return tx, nil
}

// Withdraw moves the given amount from the account's balance to the balance and records the transaction.
func Withdraw(tx *transaction.Transaction) (*transaction.Transaction, error) {

	if tx.Number == "" {
		return nil, errors.New("please provide a valid account number")
	}

	// find account
	tx.Account.Number = tx.Number
	err := tx.Account.FindByNumber()
	if err != nil {
		return nil, errors.New("we are having issues finding your account. Please try again later")
	}

	if tx.Account.User.Key == "" {
		return nil, errors.New("account not found")
	}

	// prepare the withdraw transaction
	tx.Type = global.Withdrawal
	tx.Status = global.Completed

	// update account balance
	database.PostgreSQLDBTx, err = database.PostgreSQLDB.Begin(context.Background())
	if err != nil {
		return nil, errors.New("we are having issues processing your withdraw. Please try again later")
	}

	if err := tx.Account.Withdraw(tx.Amount); err != nil {
		database.PostgreSQLDBTx.Rollback(context.Background())
		if err == pgx.ErrNoRows {
			return nil, errors.New("insufficient funds in account's available balance")
		}
		return nil, errors.New("we are having issues processing your withdraw. Please try again later")
	}

	// create transaction
	if err := tx.Create(); err != nil {
		database.PostgreSQLDBTx.Rollback(context.Background())
		return nil, errors.New("we are having issues processing your withdraw. Please try again later")
	}

	// commit transaction
	if err := database.PostgreSQLDBTx.Commit(context.Background()); err != nil {
		database.PostgreSQLDBTx.Rollback(context.Background())
		return nil, errors.New("we are having issues processing your withdraw. Please try again later")
	}

	// dispose of the transaction
	database.PostgreSQLDBTx = nil

	return tx, nil
}

// Transactions returns a list of transactions for the given account.
func Transactions(number string) (transaction.Transactions, error) {

	if number == "" {
		return nil, errors.New("please provide a valid account number")
	}

	// find account
	account := &accountr.Account{Number: number}
	err := account.FindByNumber()
	if err != nil {
		return nil, errors.New("we are having issues finding your account. Please try again later")
	}

	if account.User.Key == "" {
		return nil, errors.New("account not found")
	}

	// find transactions
	txs := transaction.Transactions{}
	err = txs.FindByAccountNumber(account.Number)
	if err != nil {
		log.Println(err)
		return nil, errors.New("we are having issues finding your transactions. Please try again later")
	}

	return txs, nil
}
