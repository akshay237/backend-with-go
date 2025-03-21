package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	db "github.com/akshay237/backend-with-go/database/sqlc"
	"github.com/gin-gonic/gin"
)

// Create Account
type CreateAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD INR"`
}

func (s *Server) CreateAccount(ctx *gin.Context) {
	// 0. create an req variable
	var req CreateAccountRequest

	// 1. if the request is not valid
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// 2. if the request is valid create the request to store acc into db
	createAccountReq := db.CreateAccountParams{
		Owner:    req.Owner,
		Currency: req.Currency,
		Balance:  0,
	}

	// 3. calls the store account func of database to create an account
	account, err := s.store.CreateAccount(ctx, createAccountReq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// 4. return the account details to the end user
	ctx.JSON(http.StatusOK, account)
}

// Get Account
type GetAccountRequest struct {
	Id int64 `uri:"id" binding:"required,min=1"`
}

func (s *Server) GetAccount(ctx *gin.Context) {

	// 1. check if we are getting valid id in request
	var req GetAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// 2. calls the get account databse function by passing the id
	account, err := s.store.GetAccount(ctx, req.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			noAccountError := fmt.Errorf("no account exists for id %d", req.Id)
			ctx.JSON(http.StatusNotFound, errorResponse(noAccountError))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// 3. return the account details
	ctx.JSON(http.StatusOK, account)
}
