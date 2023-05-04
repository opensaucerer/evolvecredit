package transaction

import (
	"net/http"

	"github.com/opensaucerer/barf"
	transactionl "github.com/opensaucerer/barf/app/logic/v1/transaction"
	transactionr "github.com/opensaucerer/barf/app/repository/v1/transaction"
)

func Transaction(w http.ResponseWriter, r *http.Request) {

	var data transactionr.Transaction
	if err := barf.Request(r).Query().Format(&data); err != nil {
		barf.Response(w).Status(http.StatusBadRequest).JSON(barf.Res{
			Status:  false,
			Message: err.Error(),
		})
		return
	}

	tx, err := transactionl.Transaction(data.SessionId)
	if err != nil {
		barf.Response(w).Status(http.StatusBadRequest).JSON(barf.Res{
			Status:  false,
			Message: err.Error(),
		})
		return
	}

	barf.Response(w).Status(http.StatusOK).JSON(barf.Res{
		Status:  true,
		Data:    tx,
		Message: "transaction retrieved",
	})
}
