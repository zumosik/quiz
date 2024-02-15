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

	return &pb.GetFileByIdResponse{File: FileToPb(file)}, nil
}

func (s *serverAPI) GetFilesByName(ctx context.Context, in *pb.GetFilesByNameRequest) (*pb.GetFilesByNameResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *serverAPI) GetFilesByUser(ctx context.Context, in *pb.GetFilesByUserRequest) (*pb.GetFilesByUserRequest, error) {
	//TODO implement me
	panic("implement me")
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
