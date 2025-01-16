package store

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/openauth/openauth/internal/backend/authn"
	backendv1 "github.com/openauth/openauth/internal/backend/gen/openauth/backend/v1"
	"github.com/openauth/openauth/internal/store/idformat"
)

func (s *Store) UpdateProjectUISettings(ctx context.Context, req *backendv1.UpdateProjectUISettingsRequest) (*backendv1.UpdateProjectUISettingsResponse, error) {
	projectID := authn.ProjectID(ctx)

	err := s.uploadToS3(req.Logo, fmt.Sprintf("project/%s/test.txt", idformat.Project.Format(projectID)))
	if err != nil {
		return nil, fmt.Errorf("failed to upload file to S3: %w", err)
	}

	return &backendv1.UpdateProjectUISettingsResponse{}, nil
}

func (s *Store) uploadToS3(file *backendv1.ImageUploadRequest, fileName string) error {
	slog.InfoContext(context.Background(), "Uploading file to S3", "bucketName", s.s3BucketName, "fileName", fileName)

	// Create an S3 PutObject input
	input := &s3.PutObjectInput{
		Bucket: aws.String(s.s3BucketName),
		Key:    aws.String(fileName),
		Body:   bytes.NewReader([]byte("test")),
		ACL:    "public-read", // Optional: Make the file publicly accessible
	}

	// Upload the file
	_, err := s.s3.PutObject(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("failed to upload file to S3: %w", err)
	}

	return nil
}
