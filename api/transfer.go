package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	db "github.com/AnkitNayan83/houseBank/db/sqlc"
	"github.com/gin-gonic/gin"
)

type transferMoneyRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"` //custom validator
}

func (server *Server) TransferMoney(ctx *gin.Context) {
	var req transferMoneyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account1, err := server.store.GetAccountById(ctx, req.FromAccountID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if account1.Currency != req.Currency {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("account: [%d] currency mismatch: %s vs %s", account1.ID, account1.Currency, req.Currency)))
		return
	}

	if account1.Balance < req.Amount {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("account: [%d] insufficient balance: %d < %d", account1.ID, account1.Balance, req.Amount)))
		return
	}

	account2, err := server.store.GetAccountById(ctx, req.ToAccountID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if account1.Currency != account2.Currency {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("account: [%d] currency mismatch: %s vs %s", account2.ID, account2.Currency, req.Currency)))
		return
	}

	arg := db.TransferMoneyTxParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        req.Amount,
	}

	TransferMoney, err := server.store.TransferMoneyTx(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, TransferMoney)
}
