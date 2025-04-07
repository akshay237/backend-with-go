package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/lib/pq"

	db "github.com/akshay237/backend-with-go/database/sqlc"
	"github.com/gin-gonic/gin"
)

const (
	UniqueKeyConstraint  = "unique_voilation"
	ForeignKeyConstraint = "foreign_key_violation"
)

// Create Account
type CreateAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,currency"`
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
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case ForeignKeyConstraint:
				err := fmt.Errorf("create user before creating account using this owner name [%s]", req.Owner)
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			case UniqueKeyConstraint:
				err := fmt.Errorf("account already exists with this username [%s] and currency [%s]", req.Owner, req.Currency)
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
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

// List Accounts
type ListAccountRequest struct {
	PageId   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (s *Server) ListAccounts(ctx *gin.Context) {

	// 1. validate the request
	var req ListAccountRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// 2. create the list accounts db functions args
	args := db.ListAccountsParams{
		Limit:  req.PageSize,
		Offset: (req.PageId - 1) * req.PageSize, // it eill used to skip the no of accounts
	}

	// 3. calls the list account db function
	accounts, err := s.store.ListAccounts(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// 4. return the accounts
	ctx.JSON(http.StatusOK, accounts)
}

// Update Account
type UpdateAccountRequest struct {
	Id      int64 `json:"id" binding:"required,min=1"`
	Balance int64 `json:"balance" binding:"required"`
}

func (s *Server) UpdateAccount(ctx *gin.Context) {

	// 1. validate the request
	var req UpdateAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// 2. create args for the update account func
	args := db.UpdateAccountParams{
		ID:      req.Id,
		Balance: req.Balance,
	}

	// 3. call the update account database func
	account, err := s.store.UpdateAccount(ctx, args)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			noAccountError := fmt.Errorf("no account exists for id %d", req.Id)
			ctx.JSON(http.StatusNotFound, errorResponse(noAccountError))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// 4. return the updated account
	ctx.JSON(http.StatusOK, account)
}

// Delete Account
type DeleteAccountRequest struct {
	Id int64 `uri:"id" binding:"required,min=1"`
}

func (s *Server) DeleteAccount(ctx *gin.Context) {

	// 1. validate the request
	var req DeleteAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// 2. calls the delete account func
	err := s.store.DeleteAccount(ctx, req.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			noAccountError := fmt.Errorf("no account exists for id %d", req.Id)
			ctx.JSON(http.StatusNotFound, errorResponse(noAccountError))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// 3. account is deleted
	ctx.JSON(http.StatusNoContent, struct{}{})
}
