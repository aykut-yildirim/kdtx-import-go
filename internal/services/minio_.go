package services

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"myapp/internal/models"
	"os"

	"github.com/joho/godotenv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type MinioService struct {
	client *s3.Client
	bucket string
}

func NewMinioService() (*MinioService, error) {

	Logger().STATUS("-- services / NewMinioService")
	godotenv.Load(".env")
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	secretKey := os.Getenv("MINIO_SECRET_KEY")
	bucket := os.Getenv("MINIO_BUCKET")
	Logger().STATUS(endpoint)
	Logger().STATUS(accessKey)
	Logger().STATUS(secretKey)
	Logger().STATUS(bucket)

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-east-1"),
		config.WithCredentialsProvider(
			aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
				return aws.Credentials{
					AccessKeyID:     accessKey,
					SecretAccessKey: secretKey,
				}, nil
			}),
		),
	)

	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String("http://" + endpoint)
		o.UsePathStyle = true
	})

	return &MinioService{
		client: client,
		bucket: bucket,
	}, nil
}

func (m *MinioService) EnsureBucket(ctx context.Context) error {

	Logger().STATUS("-- minio / EnsureBucket")

	_, err := m.client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(m.bucket),
	})

	if err == nil {
		return nil
	}

	_, err = m.client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(m.bucket),
	})

	return err
}

func (m *MinioService) GetFileByPath(
	ctx models.Context,
	bucket string,
	key string,
) ([]byte, error) {

	Logger().STATUS("-- minio / GetFileByPath")

	obj, err := m.client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		return []byte{}, fmt.Errorf("minio get_file failed: %w", err)
	}

	defer obj.Body.Close()

	data, err := io.ReadAll(obj.Body)
	if err != nil {
		return []byte{}, err
	}

	return data, nil
}

func (m *MinioService) SavePDF(
	ctx context.Context,
	task models.Task,
	pdfID string,
	pdfBytes []byte,
) (string, error) {

	Logger().STATUS("-- minio / SavePDF")

	err := m.EnsureBucket(ctx)
	if err != nil {
		return "", err
	}

	key := fmt.Sprintf(
		"%d/%d/%d/pdf/%s.pdf",
		task.ClientID,
		task.AccountID,
		task.TaskID,
		pdfID,
	)

	_, err = m.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(m.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(pdfBytes),
		ContentType: aws.String("application/pdf"),
	})

	if err != nil {
		return "", err
	}

	return key, nil
}

func (m *MinioService) SaveMinio(
	ctx context.Context,
	key string,
	body []byte,
) (string, error) {

	Logger().STATUS("-- minio / SaveMinio")

	err := m.EnsureBucket(ctx)
	if err != nil {
		return "", err
	}

	_, err = m.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(m.bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(body),
	})

	if err != nil {
		return "", err
	}

	return key, nil
}
