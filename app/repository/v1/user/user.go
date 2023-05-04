package user

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/opensaucerer/barf/app/database"
	"github.com/opensaucerer/barf/app/global"
	"github.com/opensaucerer/barf/app/reflection"
)

// Validate validates the user struct
func (u *User) Validate() error {
	if u.FirstName == "" {
		return errors.New("first name is required")
	}
	if u.LastName == "" {
		return errors.New("last name is required")
	}
	if u.Email == "" {
		return errors.New("email is required")
	}
	if u.Age < global.MinAge {
		return fmt.Errorf("age must be greater than %d", global.MinAge-1)
	}
	u.Key = ""
	u.Email = strings.ToLower(u.Email)
	u.FirstName = strings.ToUpper(string(u.FirstName[0])) + strings.ToLower(u.FirstName[1:])
	u.LastName = strings.ToUpper(string(u.LastName[0])) + strings.ToLower(u.LastName[1:])
	return nil
}

// Fields returns the struct fields as a slice of interface{} values
func (u *User) Fields() []interface{} {
	return reflection.ReturnStructFields(u)
}

// Create inserts a new user into the database.
func (u *User) Create() error {

	u.time(true)

	if err := u.key(); err != nil {
		return err
	}

	_, err := database.PostgreSQLDB.Exec(context.Background(), `INSERT INTO users (first_name, last_name, email, age, key, role, active, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`, u.FirstName, u.LastName, u.Email, u.Age, u.Key, u.Role, u.Active, u.CreatedAt, u.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

// key generates a new key if one does not currently exist on the struct
func (u *User) key() error {
	if u.Key == "" {
		salt := u.FirstName + u.LastName + u.Email + strconv.Itoa(int(u.Role))
		hash := sha256.New()
		_, err := hash.Write([]byte(salt))
		if err != nil {
			return err
		}
		u.Key = fmt.Sprintf("%x", hash.Sum(nil))
	}
	return nil
}

// time updates the CreatedAt and UpdatedAt fields on the struct
func (u *User) time(new ...bool) {
	if len(new) > 0 && new[0] {
		u.CreatedAt = time.Now().UTC()
	}
	u.UpdatedAt = time.Now().UTC()
}

// FindByEmail finds a user by their email address
func (u *User) FindByEmail() error {
	err := database.PostgreSQLDB.QueryRow(context.Background(), `SELECT id, key, first_name, last_name, email, age, role, active, created_at, updated_at FROM users WHERE email = $1`, u.Email).Scan(u.Fields()...)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil
		}
		return err
	}
	return nil
}

// FindByKey finds a user by their key
func (u *User) FindByKey() error {
	err := database.PostgreSQLDB.QueryRow(context.Background(), `SELECT id, key, first_name, last_name, email, age, role, active, created_at, updated_at FROM users WHERE key = $1`, u.Key).Scan(u.Fields()...)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil
		}
		return err
	}
	return nil
}

// Delete deletes a user from the database. This is only used for testing.
func (u *User) Delete() error {
	_, err := database.PostgreSQLDB.Exec(context.Background(), `DELETE FROM users WHERE email = $1`, u.Email)
	if err != nil {
		return err
	}
	return nil
}
