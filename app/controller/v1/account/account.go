package account

import (
	"net/http"

	"github.com/opensaucerer/barf"
	accountl "github.com/opensaucerer/barf/app/logic/v1/account"
	accountr "github.com/opensaucerer/barf/app/repository/v1/account"
	"github.com/opensaucerer/barf/app/repository/v1/transaction"
	userr "github.com/opensaucerer/barf/app/repository/v1/user"
)

func Create(w http.ResponseWriter, r *http.Request) {

	var data userr.User
	if err := barf.Request(r).Body().Format(&data); err != nil {
		barf.Response(w).Status(http.StatusBadRequest).JSON(barf.Res{
			Status:  false,
			Message: err.Error(),
		})
		return
	}

	account, err := accountl.Create(&data)
	if err != nil {
		barf.Response(w).Status(http.StatusBadRequest).JSON(barf.Res{
			Status:  false,
			Message: err.Error(),
		})
		return
	}

	barf.Response(w).Status(http.StatusCreated).JSON(barf.Res{
		Status:  true,
		Data:    account,
		Message: "account created successfully",
	})
}

func Search(w http.ResponseWriter, r *http.Request) {

	var data accountr.Account
	if err := barf.Request(r).Query().Format(&data); err != nil {
		barf.Response(w).Status(http.StatusBadRequest).JSON(barf.Res{
			Status:  false,
			Message: err.Error(),
		})
		return
	}

	account, err := accountl.Search(data.Number)
	if err != nil {
		barf.Response(w).Status(http.StatusBadRequest).JSON(barf.Res{
			Status:  false,
			Message: err.Error(),
		})
		return
	}

	barf.Response(w).Status(http.StatusOK).JSON(barf.Res{
		Status:  true,
		Data:    account,
		Message: "account retrieved",
	})
}

func Deposit(w http.ResponseWriter, r *http.Request) {

	var data transaction.Transaction
	if err := barf.Request(r).Body().Format(&data); err != nil {
		barf.Response(w).Status(http.StatusBadRequest).JSON(barf.Res{
			Status:  false,
			Message: err.Error(),
		})
		return
	}

	tx, err := accountl.Deposit(&data)
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
		Message: "deposit successful",
	})
}

func Lock(w http.ResponseWriter, r *http.Request) {

	var data transaction.Transaction
	if err := barf.Request(r).Body().Format(&data); err != nil {
		barf.Response(w).Status(http.StatusBadRequest).JSON(barf.Res{
			Status:  false,
			Message: err.Error(),
		})
		return
	}

	tx, err := accountl.Lock(&data)
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
		Message: "money locked",
	})
}

func Unlock(w http.ResponseWriter, r *http.Request) {

	var data transaction.Transaction
	if err := barf.Request(r).Body().Format(&data); err != nil {
		barf.Response(w).Status(http.StatusBadRequest).JSON(barf.Res{
			Status:  false,
			Message: err.Error(),
		})
		return
	}

	tx, err := accountl.Unlock(&data)
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
		Message: "money unlocked",
	})
}

func Withdraw(w http.ResponseWriter, r *http.Request) {

	var data transaction.Transaction
	if err := barf.Request(r).Body().Format(&data); err != nil {
		barf.Response(w).Status(http.StatusBadRequest).JSON(barf.Res{
			Status:  false,
			Message: err.Error(),
		})
		return
	}

	tx, err := accountl.Withdraw(&data)
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
		Message: "withdrawal processed",
	})
}

func Transactions(w http.ResponseWriter, r *http.Request) {

	var data accountr.Account
	if err := barf.Request(r).Query().Format(&data); err != nil {
		barf.Response(w).Status(http.StatusBadRequest).JSON(barf.Res{
			Status:  false,
			Message: err.Error(),
		})
		return
	}

	txs, err := accountl.Transactions(data.Number)
	if err != nil {
		barf.Response(w).Status(http.StatusBadRequest).JSON(barf.Res{
			Status:  false,
			Message: err.Error(),
		})
		return
	}

	barf.Response(w).Status(http.StatusOK).JSON(barf.Res{
		Status:  true,
		Data:    txs,
		Message: "transactions retrieved",
	})
}
