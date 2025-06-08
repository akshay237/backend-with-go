package api

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	db "github.com/akshay237/backend-with-go/database/sqlc"
	"github.com/akshay237/backend-with-go/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

var (
	ErrWrongPassword = errors.New("please provide valid password for login, password is incorrect")
)

// Create User
type CreateUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type userResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
}

func (s *Server) CreateUser(ctx *gin.Context) {
	// 0. create an req variable
	var req CreateUserRequest

	// 1. if the request is not valid
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// 2. hash the password before saving it to database
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// 2. if the request is valid create the request to store acc into db
	createUserReq := db.CreateUSerParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	}

	// 3. calls the store account func of database to create an account
	user, err := s.store.CreateUSer(ctx, createUserReq)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case UniqueKeyConstraint:
				err := fmt.Errorf("user already exists with this name [%s]", req.Username)
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// 5. create user response
	response := newUserResponse(user)

	// 4. return the account details to the end user
	ctx.JSON(http.StatusOK, response)
}

// GetUser API
type GetUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
}

func (s *Server) GetUser(ctx *gin.Context) {

	// 1. create an request variable
	var req GetUserRequest

	// 2. check for the request
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadGateway, errorResponse(err))
		return
	}

	// 3. make the call to get user
	user, err := s.store.GetUser(ctx, req.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// 3. return the user response
	response := newUserResponse(user)

	ctx.JSON(http.StatusOK, response)
}

// Login User
type loginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginUserResponse struct {
	Session_Id            uuid.UUID    `json:"session_id"`
	AccessToken           string       `json:"access_token"`
	AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time    `json:"refresh_token_expires_at"`
	User                  userResponse `json:"user"`
}

func (s *Server) loginUser(ctx *gin.Context) {

	// 1. check the valid request
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// 2. get the user from the db
	user, err := s.store.GetUser(ctx, req.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// 3. check if the password of the user is same or not
	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		log.Printf("Check password err: %v", err)
		ctx.JSON(http.StatusUnauthorized, errorResponse(ErrWrongPassword))
		return
	}

	// 4. create a access token for the user
	accessToken, accessPayload, err := s.tokenMaker.CreateToken(req.Username, s.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// 4.1 create a refresh token
	refreshToken, refreshPayload, err := s.tokenMaker.CreateToken(req.Username, s.config.RefreshTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// 4.2 create the session and store to DB
	session, err := s.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.Request.RemoteAddr,
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAT,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// 5. send the response to user with
	response := loginUserResponse{
		Session_Id:            session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAT,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAT,
		User:                  newUserResponse(user),
	}
	ctx.JSON(http.StatusOK, response)
}
