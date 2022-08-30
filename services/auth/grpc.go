package main

import (
	"context"

	"github.com/Omar-Belghaouti/pdash/services/auth/data"
	"github.com/Omar-Belghaouti/pdash/services/auth/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	pb.UnimplementedAuthServiceServer
}

// VerifyToken implementation for Auth gRPC server
func (s *server) VerifyToken(ctx context.Context, in *pb.Auth) (*pb.Auth, error) {
	accessToken := in.GetAccessToken()
	_, err := data.TokenMaker.VerifyToken(accessToken)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token: %s", err.Error())
	}
	return &pb.Auth{
		AccessToken: accessToken,
	}, nil
}
