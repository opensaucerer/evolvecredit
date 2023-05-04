package account

import (
	"reflect"
	"testing"
	"time"
)

// go test -v -run TestAccountRepositoryUnit ./...
func TestAccountRepositoryUnit(t *testing.T) {

	t.Run("Should return the required struct fields", func(t *testing.T) {

		data := Account{}

		fields := data.Fields()

		// 2 fields from the account struct 'owner' and 'user' should be ignored
		if len(fields) != reflect.TypeOf(data).NumField()+reflect.TypeOf(data.User).NumField()-2 {
			t.Fatalf("unexpected number of fields: got %v want %v", len(fields), reflect.TypeOf(data).NumField()+reflect.TypeOf(data.User).NumField()-2)
		}

	})

	t.Run("Should append time to the CreatedAt and UpdatedAt fields", func(t *testing.T) {

		data := Account{}

		data.time(true)

		if data.CreatedAt == (time.Time{}) {
			t.Fatalf("CreatedAt should not be empty: got %v", data.CreatedAt)
		}

		if data.UpdatedAt == (time.Time{}) {
			t.Fatalf("UpdatedAt should not be empty: got %v", data.UpdatedAt)
		}

	})

	t.Run("Should append time to the UpdatedAt field", func(t *testing.T) {

		data := Account{}

		data.time()

		if data.UpdatedAt == (time.Time{}) {
			t.Fatalf("UpdatedAt should not be empty: got %v", data.UpdatedAt)
		}

	})
}
