package main

import (
	"context"
	pb "files_test_client/pb/files"
	"fmt"
	"io"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func main() {
	conn, err := grpc.Dial("localhost:1239", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	client := pb.NewFilesClient(conn)

	ctx := context.Background()

	f, err := os.Open("a.jpg")
	if err != nil {
		panic(err)
	}

	content, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}

	req := &pb.UploadFileRequest{
		UserId: "id_user_111",
		File: &pb.File{
			Content:   content,
			Name:      "name_11",
			CreatedAt: timestamppb.Now(),
		},
	}

	resp, err := client.UploadFile(ctx, req)

	if err != nil {
		panic(err)
	}

	f2, err := os.Create("res.jpg")
	if err != nil {
		panic(err)
	}

	f2.Write(resp.File.Content)

	fmt.Println(resp.File.CreatedAt)
	fmt.Println(resp.File.Id)
	fmt.Println(resp.File.Name)

}
