package user

import (
	"reflect"
	"testing"
	"time"
)

// go test -v -run TestUserRepositoryUnit ./...
func TestUserRepositoryUnit(t *testing.T) {

	t.Run("Should return the required struct fields", func(t *testing.T) {

		data := User{}

		fields := data.Fields()

		if len(fields) != reflect.TypeOf(data).NumField() {
			t.Fatalf("unexpected number of fields: got %v want %v", len(fields), reflect.TypeOf(data).NumField())
		}

		// pointers are returned but the underlying type should be int64 (.Elem())
		if reflect.TypeOf(fields[0]).Elem().Kind() != reflect.Int64 {
			t.Fatalf("unexpected type of field: got %v want %v", reflect.TypeOf(fields[0]).Elem().Kind(), reflect.Int64)
		}

	})

	t.Run("Should append time to the CreatedAt and UpdatedAt fields", func(t *testing.T) {

		data := User{}

		data.time(true)

		if data.CreatedAt == (time.Time{}) {
			t.Fatalf("CreatedAt should not be empty: got %v", data.CreatedAt)
		}

		if data.UpdatedAt == (time.Time{}) {
			t.Fatalf("UpdatedAt should not be empty: got %v", data.UpdatedAt)
		}

	})

	t.Run("Should append time to the UpdatedAt field", func(t *testing.T) {

		data := User{}

		data.time()

		if data.UpdatedAt == (time.Time{}) {
			t.Fatalf("UpdatedAt should not be empty: got %v", data.UpdatedAt)
		}

	})

	t.Run("Should fail validation due to missing email", func(t *testing.T) {

		data := User{
			FirstName: "John",
			LastName:  "Doe",
			Age:       30,
		}

		err := data.Validate()

		if err == nil {
			t.Fatal("expected error but got none")
		}

		if err.Error() != "email is required" {
			t.Fatalf("unexpected error: got %v want %v", err.Error(), "email is required")
		}

	})

	t.Run("Should generate a key for the user", func(t *testing.T) {

		data := User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "johndoe@email.com",
			Age:       30,
		}

		data.key()

		if data.Key == "" {
			t.Fatal("expected key but got none")
		}

	})
}
