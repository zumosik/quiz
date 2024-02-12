package files

import (
	"context"
	"files/internal/domain/models"
	pb "files/pb/files"
	"log/slog"
)

type serverAPI struct {
	pb.UnimplementedFilesServer
	storage Storage
	l       *slog.Logger
}

type Storage interface {
	UploadFile(ctx context.Context, f models.File) (models.File, error)
}

func (s *serverAPI) UploadFile(ctx context.Context, in *pb.UploadFileRequest) (*pb.UploadFileResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *serverAPI) GetFileById(ctx context.Context, in *pb.GetFileByIdRequest) (*pb.GetFileByIdResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *serverAPI) GetFilesByName(ctx context.Context, in *pb.GetFilesByNameRequest) (*pb.GetFilesByNameResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *serverAPI) GetFilesByUser(ctx context.Context, in *pb.GetFilesByUserRequest) (*pb.GetFilesByUserRequest, error) {
	//TODO implement me
	panic("implement me")
}
