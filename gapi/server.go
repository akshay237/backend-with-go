package gapi

import (
	"fmt"

	db "github.com/akshay237/backend-with-go/database/sqlc"
	"github.com/akshay237/backend-with-go/pb"
	"github.com/akshay237/backend-with-go/token"
	"github.com/akshay237/backend-with-go/util"
)

// Server serves gRPC requests for our banking service.
type Server struct {
	pb.UnimplementedSimpleBankServer
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
}

// New Server creates a new gRPC server.
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

	return server, nil
}
