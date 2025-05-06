package authgrpc

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/neandrson/protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverAPI struct {
	sso.UnimplementedAuthServer
	auth Auth
}

// Interface for service layer
type Auth interface {
	Login(ctx context.Context, login string, password string, appID int) (token string, err error)
	RegisterNewUser(ctx context.Context, login string, password string) (userID int64, err error)
}

func Register(gRPCServer *grpc.Server, auth Auth) {
	sso.RegisterAuthServer(gRPCServer, &serverAPI{auth: auth})
}

// Implementations of gRPC handlers
func (s *serverAPI) Login(ctx context.Context, in *sso.LoginRequest) (*sso.LoginResponse, error) {
	if in.Login == "" {
		return nil, status.Error(codes.InvalidArgument, "login is required")
	}
	if in.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}
	if in.GetAppId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "app_id is required")
	}

	token, err := s.auth.Login(ctx, in.Login, in.Password, int(in.GetAppId()))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, status.Error(codes.NotFound, "user not found")
	}
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to login user")
	}

	return &sso.LoginResponse{Token: token}, nil
}

func (s *serverAPI) Register(ctx context.Context, in *sso.RegisterRequest) (*sso.RegisterResponse, error) {
	if in.Login == "" {
		return nil, status.Error(codes.InvalidArgument, "login is required")
	}
	if in.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	userID, err := s.auth.RegisterNewUser(ctx, in.Login, in.Password)
	if errors.Is(err, errors.New("user exists")) {
		return nil, status.Error(codes.AlreadyExists, "user already exists")
	}
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to register user")
	}

	return &sso.RegisterResponse{UserId: userID}, nil
}
