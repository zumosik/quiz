package auth

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pb "user_service/pb/auth"
)

type serverAPI struct {
	pb.UnimplementedAuthServer
	auth Auth
}

type Auth interface {
	Login(
		ctx context.Context,
		email string,
		password string,
		appID int,
	) (token string, err error)
	RegisterNewUser(
		ctx context.Context,
		email string,
		password string,
	) (userID int64, err error)
}

func Register(grpcServer *grpc.Server, auth Auth) {
	pb.RegisterAuthServer(grpcServer, &serverAPI{
		auth: auth,
	})
}

func (s *serverAPI) Login(
	ctx context.Context,
	in *pb.LoginRequest,
) (*pb.LoginResponse, error) {
	if !validateLogin(in) {
		return nil, status.Error(codes.InvalidArgument, "incorrect request")
	}
	return nil, nil

}

func (s *serverAPI) Register(
	ctx context.Context,
	in *pb.RegisterRequest,
) (*pb.RegisterResponse, error) {
	if !validateReq(in) {
		return nil, status.Error(codes.InvalidArgument, "incorrect request")
	}

	return nil, nil
}

// validateLogin returns true if all data is correct
func validateLogin(request *pb.LoginRequest) bool {
	if request == nil {
		return false
	}

	if len(request.Password) < 3 || len(request.Email) < 3 || request.AppId == 0 {
		return false
	}

	return true
}

// validateReq returns true if all data is correct
func validateReq(request *pb.RegisterRequest) bool {
	if request == nil {
		return false
	}

	if len(request.Password) < 3 || len(request.Email) < 3 {
		return false
	}

	return true
}
