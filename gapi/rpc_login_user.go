package gapi

import (
	"context"
	"database/sql"
	"errors"
	"log"

	db "github.com/akshay237/backend-with-go/database/sqlc"
	"github.com/akshay237/backend-with-go/pb"
	"github.com/akshay237/backend-with-go/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {

	// 1. check if the user exists or not
	user, err := s.store.GetUser(ctx, req.GetUsername())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "user not found: %s", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to get user: %s", err)
	}

	// 2. check if the password of the user is same or not
	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		log.Printf("Check password err: %v", err)
		return nil, status.Errorf(codes.Unauthenticated, "wrong password: %s", err)
	}

	// 3.1 create a access token for the user
	accessToken, accessPayload, err := s.tokenMaker.CreateToken(req.Username, s.config.AccessTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal server error access token creation failed: %s", err)
	}

	// 3.2 create a refresh token
	refreshToken, refreshPayload, err := s.tokenMaker.CreateToken(req.Username, s.config.RefreshTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal server error refresh token failed: %s", err)
	}

	// 3.3 create the session and store to DB
	session, err := s.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    "",
		ClientIp:     "192.128.10.21",
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAT,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create session: %s", err)
	}

	// 4. send the response to user
	response := &pb.LoginUserResponse{
		SessionId:             session.ID.String(),
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  timestamppb.New(accessPayload.ExpiredAT),
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpiredAT),
		User:                  convertUser(user),
	}

	return response, nil
}
