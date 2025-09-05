package service

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

type FileService interface {
	UploadToS3(file io.Reader, filename, contentType string) (string, error)
	CheckConnectivity(ctx context.Context) error
}

type fileService struct {
	s3Client   *s3.Client
	bucketName string
}

func NewFileService() FileService {
	// Get MinIO configuration from environment
	minioEndpoint := os.Getenv("MINIO_ENDPOINT")
	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	secretKey := os.Getenv("MINIO_SECRET_KEY")
	bucketName := os.Getenv("MINIO_BUCKET")
	useSSL := os.Getenv("MINIO_USE_SSL") == "true"

	// Validate required environment variables
	if minioEndpoint == "" {
		panic("MINIO_ENDPOINT environment variable is required")
	}
	if accessKey == "" {
		panic("MINIO_ACCESS_KEY environment variable is required")
	}
	if secretKey == "" {
		panic("MINIO_SECRET_KEY environment variable is required")
	}
	if bucketName == "" {
		panic("MINIO_BUCKET environment variable is required")
	}

	// Configure AWS SDK to work with MinIO
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-east-1"), // MinIO doesn't care about region, but AWS SDK requires it
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
	)
	if err != nil {
		panic(fmt.Sprintf("Failed to load MinIO config: %v", err))
	}

	// Create S3 client configured for MinIO
	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(fmt.Sprintf("http%s://%s", map[bool]string{true: "s", false: ""}[useSSL], minioEndpoint))
		o.UsePathStyle = true // MinIO uses path-style URLs
	})

	return &fileService{
		s3Client:   s3Client,
		bucketName: bucketName,
	}
}

func (s *fileService) UploadToS3(file io.Reader, filename, contentType string) (string, error) {
	// Generate unique filename
	ext := filepath.Ext(filename)
	uniqueFilename := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	// Create S3 key with timestamp folder structure
	now := time.Now()
	key := fmt.Sprintf("uploads/%d/%02d/%02d/%s", now.Year(), now.Month(), now.Day(), uniqueFilename)

	// Upload to S3
	_, err := s.s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(key),
		Body:        file,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to S3: %w", err)
	}

	// Generate MinIO URL (for MinIO, we can return the direct URL)
	minioEndpoint := os.Getenv("MINIO_ENDPOINT")
	useSSL := os.Getenv("MINIO_USE_SSL") == "true"
	protocol := "http"
	if useSSL {
		protocol = "https"
	}
	
	// Return direct MinIO URL
	fileURL := fmt.Sprintf("%s://%s/%s/%s", protocol, minioEndpoint, s.bucketName, key)
	return fileURL, nil
}

// CheckConnectivity tests MinIO connectivity by listing buckets
func (s *fileService) CheckConnectivity(ctx context.Context) error {
	// Perform a lightweight operation to test connectivity
	_, err := s.s3Client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return fmt.Errorf("MinIO connectivity check failed: %w", err)
	}
	return nil
}
