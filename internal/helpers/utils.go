package helper

import (
	// "context"
	// "context"
	"errors"
	"fmt"

	// "myapp/internal/models"
	"myapp/internal/models"
	"myapp/internal/services"
	"os"

	"github.com/joho/godotenv"
	// "path/filepath"
	// "strings"
)

// func ReadLocalFile(
// 	path string,
// ) ([]byte, error) {

// 	return os.ReadFile(path)
// }

// func ReadMinioFile(
// 	ctx models.Context,
// 	bucket string,
// 	key string,
// ) ([]byte, error) {
// 	minioSvc, err := services.NewMinioService()
// 	if err != nil {
// 		return []byte{}, fmt.Errorf("failed to create minio service: %w", err)
// 	}
// 	return minioSvc.GetFileByPath(ctx, bucket, key)
// }

func GetFile(ctx models.Context) ([]byte, error) {
	fmt.Println("--GetFile")
	if ctx.Task.FileMinioPath != nil {
		fmt.Println(ctx.Task.FileMinioPath)
		fmt.Println(*ctx.Task.FileMinioPath)
		fmt.Println("--GetFile--FileMinioPath")

		godotenv.Load(".env")
		bucket := os.Getenv("MINIO_BUCKET")
		key :=*ctx.Task.FileMinioPath
		minioSvc, err := services.NewMinioService()
		if err != nil {
			return []byte{}, fmt.Errorf("failed to create minio service: %w", err)
		}
		fileBytes, err := minioSvc.GetFileByPath(ctx, bucket, key)

		if err != nil {
			return nil, err
		}
		return fileBytes, nil
	}

	if ctx.Task.LocalPath != nil {
		fmt.Println("--GetFile--LocalPath")
		fmt.Println(ctx.Task.LocalPath)
		fmt.Println(*ctx.Task.LocalPath)

		fileBytes, err := os.ReadFile(*ctx.Task.LocalPath)

		


		if err != nil {
			return nil, err
		}
		return fileBytes, nil
	}

	return nil, errors.New("missing file information")
}
