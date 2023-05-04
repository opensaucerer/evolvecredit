package user

import (
	"errors"

	"github.com/opensaucerer/barf/app/global"
	userr "github.com/opensaucerer/barf/app/repository/v1/user"
)

// Register registers a new user
func Register(user *userr.User) (*userr.User, error) {

	// this is not good enough, we can improve this with a validation middleware. Ideally, I would create one from scratch as most validators have used in golang are just not sufficient for me. I have a plan to work on one here https://github.com/opensaucerer/vibranium
	if err := user.Validate(); err != nil {
		return nil, err
	}

	// ensure email is unique
	if err := user.FindByEmail(); err != nil {
		return nil, errors.New("we are having issues verifying this email address. Please try again later")
	}

	if user.Key != "" {
		return nil, errors.New("a user with this email address already exists")
	}

	user.Role = global.Customer
	user.Active = true

	// create user
	if err := user.Create(); err != nil {
		return nil, errors.New("we are having issues creating your account. Please try again later")
	}

	return user, nil
}
