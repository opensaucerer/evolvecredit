package account

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/opensaucerer/barf"
	accountc "github.com/opensaucerer/barf/app/controller/v1/account"
	userl "github.com/opensaucerer/barf/app/logic/v1/user"
	accountr "github.com/opensaucerer/barf/app/repository/v1/account"
	"github.com/opensaucerer/barf/app/repository/v1/transaction"
	userr "github.com/opensaucerer/barf/app/repository/v1/user"
	"github.com/opensaucerer/barf/app/test"
)

// go test -v -run TestAccountRouteIntegration ./...
func TestAccountRouteIntegration(t *testing.T) {

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

	t.Run("Should create an account user", func(t *testing.T) {

		userl.Register(&usr)

		// convert struct to bytes
		datab, _ := json.Marshal(usr)

		req, err := http.NewRequest("POST", "/v1/account/create", bytes.NewBuffer(datab))
		if err != nil {
			t.Fatal(err)
		}

		writer := httptest.NewRecorder()
		handler := http.HandlerFunc(accountc.Create)

		handler.ServeHTTP(writer, req)

		if status := writer.Code; status != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusCreated)
		}

		// convert response body to struct
		var res barf.Res
		if err := json.NewDecoder(writer.Body).Decode(&res); err != nil {
			t.Fatal(err)
		}

		if res.Status != true {
			t.Fatalf("handler returned unexpected status: got %v want %v",
				res.Status, true)
		}

		if res.Data.(map[string]interface{})["number"] == "" {
			t.Fatalf("handler returned an empty account number")
		}

		acc = &accountr.Account{
			Number: res.Data.(map[string]interface{})["number"].(string),
		}
	})

	t.Run("Should find and return the created account", func(t *testing.T) {

		req, err := http.NewRequest("GET", fmt.Sprintf("/v1/account/search?number=%s", acc.Number), nil)
		if err != nil {
			t.Fatal(err)
		}

		writer := httptest.NewRecorder()
		handler := http.HandlerFunc(accountc.Search)

		handler.ServeHTTP(writer, req)

		if status := writer.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusCreated)
		}

		// convert response body to struct
		var res barf.Res
		if err := json.NewDecoder(writer.Body).Decode(&res); err != nil {
			t.Fatal(err)
		}

		if res.Status != true {
			t.Fatalf("handler returned unexpected status: got %v want %v",
				res.Status, true)
		}

		if res.Data.(map[string]interface{})["number"] != acc.Number {
			t.Fatalf("handler returned an unexpected account number: got %v want %v",
				res.Data.(map[string]interface{})["number"], acc.Number)
		}
	})

	t.Run("Should deposit the specified amount into the user account", func(t *testing.T) {

		amount := 1000.0

		// convert struct to bytes
		datab, _ := json.Marshal(struct {
			Number string  `json:"number"`
			Amount float64 `json:"amount"`
		}{
			Number: acc.Number,
			Amount: amount,
		})

		req, err := http.NewRequest("PATCH", "/v1/account/deposit", bytes.NewBuffer(datab))
		if err != nil {
			t.Fatal(err)
		}

		writer := httptest.NewRecorder()
		handler := http.HandlerFunc(accountc.Deposit)

		handler.ServeHTTP(writer, req)

		if status := writer.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		// convert response body to struct
		var res barf.Res
		if err := json.NewDecoder(writer.Body).Decode(&res); err != nil {
			t.Fatal(err)
		}

		if res.Status != true {
			t.Fatalf("handler returned unexpected status: got %v want %v",
				res.Status, true)
		}

		if res.Data.(map[string]interface{})["number"] != acc.Number {
			t.Fatalf("handler returned an unexpected account number: got %v want %v",
				res.Data.(map[string]interface{})["number"], acc.Number)
		}

		if res.Data.(map[string]interface{})["amount"].(float64) != amount {
			t.Fatalf("handler returned an unexpected account balance: got %v want %v",
				res.Data.(map[string]interface{})["amount"], amount)
		}

		dtx = &transaction.Transaction{
			SessionId: res.Data.(map[string]interface{})["session_id"].(string),
			Account:   *acc,
		}

		dtx.Account.FindByNumber()
		if dtx.Account.Balance != amount {
			t.Fatalf("handler returned an unexpected account balance: got %v want %v",
				dtx.Account.Balance, amount)
		}
	})

	t.Run("Should move the specified amount to the locked balance", func(t *testing.T) {

		amount := 200.0

		// convert struct to bytes
		datab, _ := json.Marshal(struct {
			Number string  `json:"number"`
			Amount float64 `json:"amount"`
		}{
			Number: acc.Number,
			Amount: amount,
		})

		req, err := http.NewRequest("PATCH", "/v1/account/lock", bytes.NewBuffer(datab))
		if err != nil {
			t.Fatal(err)
		}

		writer := httptest.NewRecorder()
		handler := http.HandlerFunc(accountc.Lock)

		handler.ServeHTTP(writer, req)

		if status := writer.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		// convert response body to struct
		var res barf.Res
		if err := json.NewDecoder(writer.Body).Decode(&res); err != nil {
			t.Fatal(err)
		}

		if res.Status != true {
			t.Fatalf("handler returned unexpected status: got %v want %v",
				res.Status, true)
		}

		if res.Data.(map[string]interface{})["number"] != acc.Number {
			t.Fatalf("handler returned an unexpected account number: got %v want %v",
				res.Data.(map[string]interface{})["number"], acc.Number)
		}

		if res.Data.(map[string]interface{})["amount"].(float64) != amount {
			t.Fatalf("handler returned an unexpected account balance: got %v want %v",
				res.Data.(map[string]interface{})["amount"], amount)
		}

		ltx = &transaction.Transaction{
			SessionId: res.Data.(map[string]interface{})["session_id"].(string),
			Account:   *acc,
		}

		ltx.Account.FindByNumber()
		if ltx.Account.LockedBalance != amount {
			t.Fatalf("handler returned an unexpected locked balance: got %v want %v",
				ltx.Account.LockedBalance, amount)
		}
	})

	t.Run("Should move the specified amount out of the locked balance", func(t *testing.T) {

		amount := 200.0

		// convert struct to bytes
		datab, _ := json.Marshal(struct {
			Number string  `json:"number"`
			Amount float64 `json:"amount"`
		}{
			Number: acc.Number,
			Amount: amount,
		})

		req, err := http.NewRequest("PATCH", "/v1/account/lock", bytes.NewBuffer(datab))
		if err != nil {
			t.Fatal(err)
		}

		writer := httptest.NewRecorder()
		handler := http.HandlerFunc(accountc.Unlock)

		handler.ServeHTTP(writer, req)

		if status := writer.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		// convert response body to struct
		var res barf.Res
		if err := json.NewDecoder(writer.Body).Decode(&res); err != nil {
			t.Fatal(err)
		}

		if res.Status != true {
			t.Fatalf("handler returned unexpected status: got %v want %v",
				res.Status, true)
		}

		if res.Data.(map[string]interface{})["number"] != acc.Number {
			t.Fatalf("handler returned an unexpected account number: got %v want %v",
				res.Data.(map[string]interface{})["number"], acc.Number)
		}

		if res.Data.(map[string]interface{})["amount"].(float64) != amount {
			t.Fatalf("handler returned an unexpected account balance: got %v want %v",
				res.Data.(map[string]interface{})["amount"], amount)
		}

		utx = &transaction.Transaction{
			SessionId: res.Data.(map[string]interface{})["session_id"].(string),
			Account:   *acc,
		}

		utx.Account.FindByNumber()
		if utx.Account.LockedBalance != 0 {
			t.Fatalf("handler returned an unexpected locked balance: got %v want %v",
				ltx.Account.LockedBalance, 0)
		}
	})

	t.Run("Should retrieve the transactions associated for the account", func(t *testing.T) {

		req, err := http.NewRequest("GET", fmt.Sprintf("/v1/account/transactions?number=%s", acc.Number), nil)
		if err != nil {
			t.Fatal(err)
		}

		writer := httptest.NewRecorder()
		handler := http.HandlerFunc(accountc.Transactions)

		handler.ServeHTTP(writer, req)

		if status := writer.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusCreated)
		}

		// convert response body to struct
		var res barf.Res
		if err := json.NewDecoder(writer.Body).Decode(&res); err != nil {
			t.Fatal(err)
		}

		if res.Status != true {
			t.Fatalf("handler returned unexpected status: got %v want %v",
				res.Status, true)
		}

		if len(res.Data.([]interface{})) != 3 {
			t.Fatalf("handler returned an unexpected number of transactions: got %v want %v",
				len(res.Data.([]interface{})), 3)
		}
	})

	// clean up
	usr.Delete()
	acc.Delete()
	dtx.Delete()
	ltx.Delete()
	utx.Delete()

}
