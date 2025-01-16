package store

import (
	"bytes"
	"context"
	"fmt"
	"mime"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	backendv1 "github.com/openauth/openauth/internal/backend/gen/openauth/backend/v1"
	"github.com/openauth/openauth/internal/common/apierror"
	"github.com/openauth/openauth/internal/store/idformat"
)

func (s *Store) getURLForFileKey(fileKey string) string {
	return fmt.Sprintf("%s/%s", s.userContentUrl, fileKey)
}

func (s *Store) uploadToS3(ctx context.Context, file *backendv1.ImageUploadRequest, fileKey string) error {
	// Create an S3 PutObject input
	input := &s3.PutObjectInput{
		ACL:         "public-read", // Optional: Make the file publicly accessible
		Body:        bytes.NewReader(file.Data),
		Bucket:      aws.String(s.s3BucketName),
		ContentType: aws.String(file.MimeType),
		Key:         aws.String(fileKey),
	}

	// Upload the file
	_, err := s.s3.PutObject(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to upload file to S3: %w", err)
	}

	return nil
}

func getFileKeyForImageType(projectID uuid.UUID, imageType string, mimeType string) (string, error) {
	ext, err := mime.ExtensionsByType(mimeType)
	if err != nil {
		return "", apierror.NewInvalidArgumentError("invalid logo mimetype", fmt.Errorf("failed to get extension by mimetype: %w", err))
	}

	fileName := fmt.Sprintf("projects/%s/%s.%s", idformat.Project.Format(projectID), imageType, ext[0])

	return fileName, nil
}
