package storetestutil

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func NewS3Client(t *testing.T) *s3.Client {
	return s3.NewFromConfig(*aws.NewConfig())
}
