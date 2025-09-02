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
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

type FileService interface {
	UploadToS3(file io.Reader, filename, contentType string) (string, error)
}

type fileService struct {
	s3Client   *s3.Client
	bucketName string
}

func NewFileService() FileService {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(fmt.Sprintf("Failed to load AWS config: %v", err))
	}

	s3Client := s3.NewFromConfig(cfg)
	bucketName := os.Getenv("AWS_S3_BUCKET")
	if bucketName == "" {
		panic("AWS_S3_BUCKET environment variable is required")
	}

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

	// Generate pre-signed URL (valid for 24 hours)
	presignClient := s3.NewPresignClient(s.s3Client)
	request, err := presignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(24 * time.Hour)
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate pre-signed URL: %w", err)
	}

	return request.URL, nil
}
