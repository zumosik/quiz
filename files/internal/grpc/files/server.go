package files

import (
	"context"
	"errors"
	"files/internal/domain/models"
	"files/internal/storage"
	pb "files/pb/files"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log/slog"
)

type serverAPI struct {
	pb.UnimplementedFilesServer
	storage Storage
	l       *slog.Logger
}

type Storage interface {
	UploadFile(ctx context.Context, f models.File) (models.File, error)
	GetFileById(ctx context.Context, id string) (models.File, error)
	GetFilesByName(ctx context.Context, name string, limit int) ([]models.File, error)
	GetFilesByUser(ctx context.Context, userId string, limit int) ([]models.File, error)
}

func Register(grpcServer *grpc.Server, storage Storage, logger *slog.Logger) {
	pb.RegisterFilesServer(grpcServer, &serverAPI{
		storage: storage,
		l:       logger,
	})
}

func (s *serverAPI) UploadFile(ctx context.Context, in *pb.UploadFileRequest) (*pb.UploadFileResponse, error) {
	const op = "internal/grpc/auth/server/Login()"
	log := s.l.With(slog.String("op", op))

	if !validateUploadFile(in) {
		log.Error("haven't passed validation")
		return nil, status.Error(codes.InvalidArgument, "incorrect request")
	}

	f := PbToFile(in.File)
	f.UserID = in.UserId

	file, err := s.storage.UploadFile(ctx, f)
	if err != nil {
		return nil, status.Error(codes.Internal, "incorrect request")
	}

	return &pb.UploadFileResponse{File: FileToPb(file)}, nil
}

// GetFileById will return only for this user
func (s *serverAPI) GetFileById(ctx context.Context, in *pb.GetFileByIdRequest) (*pb.GetFileByIdResponse, error) {
	const op = "internal/grpc/auth/server/GetFileById()"
	log := s.l.With(slog.String("op", op))
	if !validateGetFileById(in) {
		log.Error("haven't passed validation")
		return nil, status.Error(codes.InvalidArgument, "incorrect request")
	}

	file, err := s.storage.GetFileById(ctx, in.UserId)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "not found")
		}
		return nil, status.Error(codes.Internal, "incorrect request")
	}

	if file.UserID != in.UserId {
		// TODO: replace with permission checking
		return nil, status.Error(codes.PermissionDenied, "permission denied")
	}

	return &pb.GetFileByIdResponse{File: FileToPb(file)}, nil
}

// GetFilesByName will return files by name for ALL users
func (s *serverAPI) GetFilesByName(ctx context.Context, in *pb.GetFilesByNameRequest) (*pb.GetFilesByNameResponse, error) {
	const limit = 10

	const op = "internal/grpc/auth/server/GetFilesByName()"
	log := s.l.With(slog.String("op", op))
	if !validateGetFilesByName(in) {
		log.Error("haven't passed validation")
		return nil, status.Error(codes.InvalidArgument, "incorrect request")
	}

	files, err := s.storage.GetFilesByName(ctx, in.Name, limit)
	if err != nil {
		return nil, err
	}

	var filesPb []*pb.File

	for _, f := range files {
		filesPb = append(filesPb, FileToPb(f))
	}

	return &pb.GetFilesByNameResponse{Files: filesPb}, err
}

func (s *serverAPI) GetFilesByUser(ctx context.Context, in *pb.GetFilesByUserRequest) (*pb.GetFilesByUserResponse, error) {
	const limit = 10

	const op = "internal/grpc/auth/server/GetFilesByName()"
	log := s.l.With(slog.String("op", op))
	if !validateGetFilesByUser(in) {
		log.Error("haven't passed validation")
		return nil, status.Error(codes.InvalidArgument, "incorrect request")
	}

	files, err := s.storage.GetFilesByUser(ctx, in.UserId, limit)
	if err != nil {
		return nil, err
	}

	var filesPb []*pb.File

	for _, f := range files {
		filesPb = append(filesPb, FileToPb(f))
	}

	return &pb.GetFilesByUserResponse{Files: filesPb}, err
}

// validateUploadFile returns true if all data is correct
func validateUploadFile(in *pb.UploadFileRequest) bool {
	if len(in.UserId) < 3 || len(in.File.Content) < 1 || len(in.File.Name) < 1 {
		return false
	}

	return true

}

// validateGetFileById returns true if all data is correct
func validateGetFileById(in *pb.GetFileByIdRequest) bool {
	return !(len(in.UserId) < 3 || len(in.Id) < 3)
}

// validateGetFilesByName returns true if all data is correct
func validateGetFilesByName(in *pb.GetFilesByNameRequest) bool {
	return !(len(in.UserId) < 3 || len(in.Name) < 3)
}

// validateGetFilesByUser returns true if all data is correct
func validateGetFilesByUser(in *pb.GetFilesByUserRequest) bool {
	return !(len(in.UserId) < 3)
}

func FileToPb(file models.File) *pb.File {
	return &pb.File{
		Id:        file.ID,
		Content:   file.Bytes,
		Name:      file.Name,
		CreatedAt: timestamppb.New(file.CreatedAt),
	}
}

// PbToFile WILL NOT set userID
func PbToFile(file *pb.File) models.File {
	return models.File{
		Name:      file.Name,
		CreatedAt: file.CreatedAt.AsTime(),
		Bytes:     file.Content,
	}
}
