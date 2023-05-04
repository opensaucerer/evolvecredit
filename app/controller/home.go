package controller

import (
	"net/http"

	"github.com/opensaucerer/barf"
	"github.com/opensaucerer/barf/app/types"
)

func Home(w http.ResponseWriter, r *http.Request) {

	data := types.Home{
		Status:      true,
		Version:     "1.0.0",
		Name:        "Zeina MFI",
		Description: "Banking as a service.",
		Website:     "https://zeinamfibyopensaucerer.onrender.com",
	}

	barf.Response(w).Status(http.StatusOK).JSON(barf.Res{
		Status:  true,
		Data:    data,
		Message: "Zeina MFI",
	})
}
