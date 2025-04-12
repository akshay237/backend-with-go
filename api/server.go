package api

import (
	db "github.com/akshay237/backend-with-go/database/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server Serves HTTP requests for our banking service.
type Server struct {
	store  db.Store
	Router *gin.Engine
}

// New Server creates a new HTTP server and setup routing.
func NewServerHandler(store db.Store) *Server {
	server := &Server{
		store: store,
	}
	// routes
	router := gin.Default()

	// add the validator middleware
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	// account apis
	router.POST("/accounts", server.CreateAccount)
	router.GET("/accounts/:id", server.GetAccount)
	router.GET("/accounts", server.ListAccounts)
	router.PUT("/accounts", server.UpdateAccount)
	router.DELETE("/accounts/:id", server.DeleteAccount)

	// transfer api
	router.POST("/transfers", server.createTransfer)

	// user apis
	router.POST("/users", server.CreateUser)
	router.POST("/users/byusername", server.GetUser)

	server.Router = router
	return server
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
