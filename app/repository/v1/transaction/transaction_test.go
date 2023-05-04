package transaction

import (
	"reflect"
	"testing"
	"time"
)

// go test -v -run TestTransactionRepositoryUnit ./...
func TestTransactionRepositoryUnit(t *testing.T) {

	t.Run("Should return the required struct fields", func(t *testing.T) {

		data := Transaction{}

		fields := data.Fields()

		// 1 field from the transaction struct 'account' should be ignore - 2 fields from the account struct 'owner' and 'user' should be ignored
		expected := reflect.TypeOf(data).NumField() + reflect.TypeOf(data.Account).NumField() + reflect.TypeOf(data.Account.User).NumField() - 3

		if len(fields) != expected {
			t.Fatalf("unexpected number of fields: got %v want %v", len(fields), expected)
		}

	})

	t.Run("Should add a session id to the transaction", func(t *testing.T) {

		data := Transaction{}

		data.session()

		if data.SessionId == "" {
			t.Fatalf("SessionID should not be empty")
		}

	})

	t.Run("Should not update the session id if it is already set", func(t *testing.T) {

		data := Transaction{
			SessionId: "123",
		}

		data.session()

		if data.SessionId != "123" {
			t.Fatalf("SessionID should not be updated: expected %v got %v", "123", data.SessionId)
		}

	})

	t.Run("Should append time to the CreatedAt and UpdatedAt fields", func(t *testing.T) {

		data := Transaction{}

		data.time(true)

		if data.CreatedAt == (time.Time{}) {
			t.Fatalf("CreatedAt should not be empty: got %v", data.CreatedAt)
		}

		if data.UpdatedAt == (time.Time{}) {
			t.Fatalf("UpdatedAt should not be empty: got %v", data.UpdatedAt)
		}

	})

	t.Run("Should append time to the UpdatedAt field", func(t *testing.T) {

		data := Transaction{}

		data.time()

		if data.UpdatedAt == (time.Time{}) {
			t.Fatalf("UpdatedAt should not be empty: got %v", data.UpdatedAt)
		}

	})
}
