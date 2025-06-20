package storetesting

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/require"
)

func TestNewS3(t *testing.T) {
	t.Parallel()

	testS3, cleanup := newS3()
	t.Cleanup(cleanup)

	res, err := testS3.Client.ListBuckets(t.Context(), &s3.ListBucketsInput{})
	require.NoError(t, err, "failed to list S3 buckets")
	require.NotEmpty(t, res.Buckets, "S3 buckets should not be empty")
}
