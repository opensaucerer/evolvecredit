package user

import (
	"github.com/opensaucerer/barf"
	userc "github.com/opensaucerer/barf/app/controller/v1/user"
)

func RegisterUserRoutes() {
	barf.Post("/v1/user/register", userc.Register)
}
