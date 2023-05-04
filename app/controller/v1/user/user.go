package user

import (
	"net/http"

	"github.com/opensaucerer/barf"
	userl "github.com/opensaucerer/barf/app/logic/v1/user"
	userr "github.com/opensaucerer/barf/app/repository/v1/user"
)

func Register(w http.ResponseWriter, r *http.Request) {

	var data userr.User
	if err := barf.Request(r).Body().Format(&data); err != nil {
		barf.Response(w).Status(http.StatusBadRequest).JSON(barf.Res{
			Status:  false,
			Message: err.Error(),
		})
		return
	}

	user, err := userl.Register(&data)
	if err != nil {
		barf.Response(w).Status(http.StatusBadRequest).JSON(barf.Res{
			Status:  false,
			Message: err.Error(),
		})
		return
	}

	barf.Response(w).Status(http.StatusCreated).JSON(barf.Res{
		Status:  true,
		Data:    user,
		Message: "user created successfully",
	})
}
