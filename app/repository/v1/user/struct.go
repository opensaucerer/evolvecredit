package user

import (
	"time"

	"github.com/opensaucerer/barf/app/global"
)

type User struct {
	Id        int64       `json:"-"`
	Key       string      `json:"key"`
	FirstName string      `json:"first_name"`
	LastName  string      `json:"last_name"`
	Email     string      `json:"email"`
	Age       int         `json:"age"`
	Role      global.Role `json:"role"`
	Active    bool        `json:"active"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}
