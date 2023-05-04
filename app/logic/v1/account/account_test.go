package account

import (
	"testing"

	userl "github.com/opensaucerer/barf/app/logic/v1/user"
	accountr "github.com/opensaucerer/barf/app/repository/v1/account"
	"github.com/opensaucerer/barf/app/repository/v1/transaction"
	userr "github.com/opensaucerer/barf/app/repository/v1/user"
	"github.com/opensaucerer/barf/app/test"
)

// go test -v -run TestAccountLogicIntegration ./...
func TestAccountLogicIntegration(t *testing.T) {

	test.Setup()

	usr := userr.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "johndoe@email.com",
		Age:       30,
	}
	var acc *accountr.Account
	var dtx *transaction.Transaction
	var ltx *transaction.Transaction
	var utx *transaction.Transaction

	t.Run("Should create an account for user", func(t *testing.T) {

		_, err := userl.Register(&usr)
		if err != nil {
			t.Fatal(err)
		}

		acc, err = Create(&usr)
		if err != nil {
			t.Fatal(err)
		}

		if acc.Number == "" {
			t.Fatalf("account number should not be empty: got %v", acc.Number)
		}

		if acc.User.Email != usr.Email {
			t.Fatalf("logic returned unexpected email: got %v want %v", acc.User.Email, usr.Email)
		}

	})

	t.Run("Should deposit the specified amount into the user account", func(t *testing.T) {

		amount := 1000.0

		var err error

		dtx, err = Deposit(&transaction.Transaction{
			Number: acc.Number,
			Amount: amount,
		})
		if err != nil {
			t.Fatal(err)
		}

		if dtx.Number != acc.Number {
			t.Fatalf("logic returned an unexpected account number: got %v want %v",
				dtx.Number, acc.Number)
		}

		if dtx.Amount != amount {
			t.Fatalf("logic returned an unexpected account balance: got %v want %v",
				dtx.Amount, amount)
		}

	})

	t.Run("Should move the specified amount to the locked balance", func(t *testing.T) {

		amount := 200.0

		var err error

		ltx, err = Lock(&transaction.Transaction{
			Number: acc.Number,
			Amount: amount,
		})
		if err != nil {
			t.Fatal(err)
		}

		if ltx.Number != acc.Number {
			t.Fatalf("logic returned an unexpected account number: got %v want %v",
				ltx.Number, acc.Number)
		}

		if ltx.Amount != amount {
			t.Fatalf("handler returned an unexpected account balance: got %v want %v",
				ltx.Amount, amount)
		}

		ltx.Account.FindByNumber()
		if ltx.Account.LockedBalance != amount {
			t.Fatalf("handler returned an unexpected locked balance: got %v want %v",
				ltx.Account.LockedBalance, amount)
		}
	})

	t.Run("Should move the specified amount out of the locked balance", func(t *testing.T) {

		amount := 200.0

		var err error

		utx, err = Unlock(&transaction.Transaction{
			Number: acc.Number,
			Amount: amount,
		})
		if err != nil {
			t.Fatal(err)
		}

		if utx.Number != acc.Number {
			t.Fatalf("logic returned an unexpected account number: got %v want %v",
				ltx.Number, acc.Number)
		}

		if utx.Amount != amount {
			t.Fatalf("handler returned an unexpected account balance: got %v want %v",
				ltx.Amount, amount)
		}

		utx.Account.FindByNumber()
		if utx.Account.LockedBalance != 0 {
			t.Fatalf("handler returned an unexpected locked balance: got %v want %v",
				utx.Account.LockedBalance, 0)
		}

	})

	t.Run("Should retrieve the transactions associated for the account", func(t *testing.T) {

		txs, err := Transactions(acc.Number)
		if err != nil {
			t.Fatal(err)
		}

		if len(txs) != 3 {
			t.Fatalf("logic returned an unexpected number of transactions: got %v want %v",
				len(txs), 3)
		}
	})

	// clean up
	usr.Delete()
	acc.Delete()
	dtx.Delete()
	ltx.Delete()
	utx.Delete()
}
