package api

import (
	"fmt"

	db "github.com/akshay237/backend-with-go/database/sqlc"
	"github.com/akshay237/backend-with-go/token"
	"github.com/akshay237/backend-with-go/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server Serves HTTP requests for our banking service.
type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	Router     *gin.Engine
}

// New Server creates a new HTTP server and setup routing.
func NewServerHandler(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasteoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %v", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	// add the validator middleware
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	// routes
	router := gin.Default()

	// user apis
	router.POST("/users", server.CreateUser)
	router.POST("/user/login", server.loginUser)
	router.POST("/token/renew_access", server.renewAccessToken)

	// add the middlewares to all other routes
	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	// account apis
	authRoutes.POST("/accounts", server.CreateAccount)
	authRoutes.GET("/accounts/:id", server.GetAccount)
	authRoutes.GET("/accounts", server.ListAccounts)
	authRoutes.PUT("/accounts", server.UpdateAccount)
	authRoutes.DELETE("/accounts/:id", server.DeleteAccount)

	// transfer api
	authRoutes.POST("/transfers", server.createTransfer)

	server.Router = router
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
