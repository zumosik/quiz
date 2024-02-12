package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type File struct {
	ID   string
	Name string
}

func main() {
	ctx := context.Background()
	opt := option.WithCredentialsFile("opt.json")
	config := &firebase.Config{
		DatabaseURL:   "gs://files-saver-2233.appspot.com",
		StorageBucket: "files-saver-2233.appspot.com",
	}

	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		panic(err)
	}

	client, err := app.Storage(ctx)
	if err != nil {
		panic(err)
	}

	bucket, err := client.DefaultBucket()
	if err != nil {
		panic(err)
	}

	// f, err := os.Open("download.png")
	// if err != nil {
	// 	panic(err)
	// }
	// defer f.Close()

	// f1, err := os.Open("r.jpg")
	// if err != nil {
	// 	panic(err)
	// }
	// defer f1.Close()

	// ff := []*os.File{f, f1}
	// var ids []string

	// time.Sleep(2 * time.Second)

	// for i, file := range ff {
	// 	fmt.Println(i)

	// 	uid := uuid.New()
	// 	id := uid.String()
	// 	ids = append(ids, id)

	// 	obj := bucket.Object(file.Name())
	// 	w := obj.NewWriter(ctx)
	// 	w.ObjectAttrs.Metadata = map[string]string{"firebaseStorageDownloadTokens": id}
	// 	defer func(w *storage.Writer) {
	// 		err := w.Close()
	// 		if err != nil {
	// 			panic(err)
	// 		}
	// 	}(w)

	// 	if _, err := io.Copy(w, file); err != nil {
	// 		panic(err)
	// 	}
	// }

	var ids []string

	files, err := listFiles(ctx, bucket)
	if err != nil {
		fmt.Printf("Error listing files: %v", err)
	}

	// Print the list of files
	for _, file := range files {
		ids = append(ids, file.ID)
		fmt.Println(file)
	}

	for _, id := range ids {
		err = downloadFileByID(ctx, bucket, id)
		if err != nil {
			panic(err)

		}
	}

	fmt.Print("success")
}
func downloadFileByID(ctx context.Context, bucket *storage.BucketHandle, fileID string) error {
	// Iterate through objects in the bucket
	it := bucket.Objects(ctx, nil)
	for {
		attrs, err := it.Next()

		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}

		// Check if metadata contains the file ID
		if attrs.Metadata["firebaseStorageDownloadTokens"] == fileID {
			// Open the object
			reader, err := bucket.Object(attrs.Name).NewReader(ctx)
			if err != nil {
				return err
			}
			defer reader.Close()

			// Create a file to write the object contents
			file, err := os.Create(fileID + ".png") // Change the file name as needed
			if err != nil {
				return err
			}
			defer file.Close()

			// Copy object contents to the file
			if _, err := io.Copy(file, reader); err != nil {
				return err
			}

			return nil // File downloaded successfully
		}
	}

	return fmt.Errorf("file with ID %s not found", fileID)
}

func listFiles(ctx context.Context, bucket *storage.BucketHandle) ([]File, error) {
	var fileList []File

	// Iterate through objects in the bucket
	it := bucket.Objects(ctx, nil)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		// Append file information to the list
		fileList = append(fileList, File{
			ID:   attrs.Metadata["firebaseStorageDownloadTokens"],
			Name: attrs.Name,
		})
	}

	return fileList, nil
}
