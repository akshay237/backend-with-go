package api

import (
	db "github.com/akshay237/backend-with-go/database/sqlc"
	"github.com/gin-gonic/gin"
)

// Server Serves HTTP requests for our banking service.
type Server struct {
	store  *db.Store
	router *gin.Engine
}

// New Server creates a new HTTP server and setup routing.
func NewServer(store *db.Store) *Server {
	server := &Server{
		store: store,
	}

	router := gin.Default()

	router.POST("/accounts", server.CreateAccount)
	server.router = router
	return server
}

// Start runs the HTTP server on the address provided.
func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
