package tests

import (
	"context"
	pb "files_test_client/pb/files"
	"io"
	"os"
	"reflect"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func connectToServer(t *testing.T, target string) pb.FilesClient {
	conn, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	client := pb.NewFilesClient(conn)

	return client
}

func TestUploadFileGetFilesByID(t *testing.T) {
	cl := connectToServer(t, "localhost:1239")

	ctx := context.Background()

	type testCase struct {
		name    string
		file    *pb.UploadFileRequest
		wantErr bool
	}

	f, err := os.Open("a.jpg")
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	defer f.Close()

	bytes, err := io.ReadAll(f)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	testCases := []testCase{
		{
			name: "ok",
			file: &pb.UploadFileRequest{
				UserId: "user1",
				File: &pb.File{
					Name:      "file_name_1",
					CreatedAt: timestamppb.Now(),
					Content:   bytes,
				},
			},
			wantErr: false,
		},
		{
			name: "ok",
			file: &pb.UploadFileRequest{
				UserId: "user1",
				File: &pb.File{
					Name:      "file_name_2",
					CreatedAt: timestamppb.Now(),
					Content:   bytes,
				},
			},
			wantErr: false,
		},
		{
			name: "empty",
			file: &pb.UploadFileRequest{
				UserId: "",
				File: &pb.File{
					Name:      "",
					CreatedAt: nil,
					Content:   nil,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range testCases {
		resp, err := cl.UploadFile(ctx, tt.file)
		if err != nil {
			if !tt.wantErr {
				t.Errorf("error in upload file: %v, wantErr = %v", err, tt.wantErr)
				continue
			}
		}
		if resp.File.Name != tt.file.File.Name || resp.File.CreatedAt == tt.file.File.CreatedAt {
			t.Error("file metadata not match")
			continue
		}

		resp2, err := cl.GetFileById(ctx, &pb.GetFileByIdRequest{
			UserId: tt.file.UserId,
			Id:     resp.File.Id,
		})
		if err != nil {
			t.Errorf("error in get file by id: %v", err)
			continue

		}

		if !reflect.DeepEqual(resp2.File, resp.File) {
			t.Error("resp from GetFileById is not equal to resp from UploadFile")
			continue

		}
	}
}
