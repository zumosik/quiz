package auth

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"time"
	"user_service/internal/domain/models"
	"user_service/internal/storage"
	"user_service/lib/utils"
	pb "user_service/pb/auth"
)

type serverAPI struct {
	pb.UnimplementedAuthServer
	storage  Storage
	l        *slog.Logger
	secret   string
	tokenTTL time.Duration
}

type Storage interface {
	SaveUser(ctx context.Context, u models.User) (models.User, error)
	FindUserByEmail(ctx context.Context, email string) (models.User, error)
}

func Register(grpcServer *grpc.Server, storage Storage, logger *slog.Logger, secret string, tokenTTL time.Duration) {
	pb.RegisterAuthServer(grpcServer, &serverAPI{
		storage:  storage,
		l:        logger,
		secret:   secret,
		tokenTTL: tokenTTL,
	})
}

func (s *serverAPI) Login(
	ctx context.Context,
	in *pb.LoginRequest,
) (*pb.LoginResponse, error) {
	const op = "internal/grpc/auth/server/Login()"
	log := s.l.With(slog.String("op", op))

	log.Debug("starting login")

	if !validateLogin(in) {
		log.Error("haven't passed validation")
		return nil, status.Error(codes.InvalidArgument, "incorrect request")
	}

	log.Debug("passed validation")

	u, err := s.storage.FindUserByEmail(ctx, in.Email)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			log.Error("user not found")
			return nil, status.Error(codes.InvalidArgument, "incorrect request")
		}
		log.Error("error in FindUserByEmail", utils.WrapErr(err))
		return nil, status.Error(codes.Internal, "internal error")
	}

	log.Debug("found user")

	if !u.ComparePassword(in.Password) {
		log.Error("invalid password")
		return nil, status.Error(codes.InvalidArgument, "incorrect request")
	}

	log.Debug("creating jwt...")

	payload := jwt.MapClaims{
		"id":    u.ID,
		"email": u.Email,
		"exp":   time.Now().Add(s.tokenTTL).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	t, err := token.SignedString([]byte(s.secret))
	if err != nil {
		log.Error("cant create string token", utils.WrapErr(err))
		return nil, status.Error(codes.Internal, "internal error")
	}

	log.Debug(t)

	return &pb.LoginResponse{Token: t}, nil

}

func (s *serverAPI) Register(
	ctx context.Context,
	in *pb.RegisterRequest,
) (*pb.RegisterResponse, error) {
	// TODO: add unique email validation (error)
	const op = "internal/grpc/auth/server/Login()"
	log := s.l.With(slog.String("op", op))

	if !validateReq(in) {
		log.Error("haven't passed validation")
		return nil, status.Error(codes.InvalidArgument, "incorrect request")
	}

	user, err := s.storage.SaveUser(ctx, models.User{
		Email:    in.Email,
		Password: in.Password,
	})
	if err != nil {
		log.Error("error in SaveUser", utils.WrapErr(err))
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &pb.RegisterResponse{UserId: user.ID}, nil
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
