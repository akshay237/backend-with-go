package gapi

import (
	"context"

	db "github.com/akshay237/backend-with-go/database/sqlc"

	"github.com/akshay237/backend-with-go/pb"
	"github.com/akshay237/backend-with-go/util"
	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	UniqueKeyConstraint  = "unique_voilation"
	ForeignKeyConstraint = "foreign_key_violation"
)

func (s *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {

	// 1. hash the password before saving it to database
	hashedPassword, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %v", err)
	}

	// 2. if the request is valid create the request to store acc into db
	createUserReq := db.CreateUSerParams{
		Username:       req.GetUsername(),
		HashedPassword: hashedPassword,
		FullName:       req.GetFullName(),
		Email:          req.GetEmail(),
	}

	// 3. calls the store account func of database to create an account
	user, err := s.store.CreateUSer(ctx, createUserReq)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case UniqueKeyConstraint:
				return nil, status.Errorf(codes.AlreadyExists, "username already exists: %s", err)
			}
		}
		return nil, status.Errorf(codes.Internal, "failed to create user: %s", err)
	}

	// 5. create user response
	response := &pb.CreateUserResponse{
		User: convertUser(user),
	}

	// 4. return the account details to the end user
	return response, nil
}
