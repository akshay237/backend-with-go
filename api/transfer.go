package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	db "github.com/akshay237/backend-with-go/db/sqlc"
	"github.com/gin-gonic/gin"
)

type transferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (s *Server) createTransfer(ctx *gin.Context) {
	// 0. create an req variable
	var req transferRequest

	// 1. if the request is not valid
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if !s.validAccount(ctx, req.FromAccountID, req.Currency) {
		return
	}

	if !s.validAccount(ctx, req.ToAccountID, req.Currency) {
		return
	}

	// 2. if the request is valid create the transfer req
	createTransferReq := db.TransferTxParams{
		FromAccountId: req.FromAccountID,
		ToAccountId:   req.ToAccountID,
		Amount:        req.Amount,
	}

	// 3. calls the store account func of database to create an account
	result, err := s.store.TransferTx(ctx, createTransferReq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// 4. return the account details to the end user
	ctx.JSON(http.StatusOK, result)
}

func (s *Server) validAccount(ctx *gin.Context, accountId int64, currency string) bool {

	account, err := s.store.GetAccount(ctx, accountId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch is %s and %s", accountId, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return false
	}

	return true
}
