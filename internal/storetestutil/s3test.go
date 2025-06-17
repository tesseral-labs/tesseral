package storetestutil

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/testcontainers/testcontainers-go"
)

func NewS3(t *testing.T) *s3.Client {
	container, err := testcontainers.GenericContainer(
		t.Context(),
		testcontainers.GenericContainerRequest{
			ContainerRequest: testcontainers.ContainerRequest{
				Name:         "tesseral-test-local-s3",
				Image:        "adobe/s3mock:latest",
				ExposedPorts: []string{"9090/tcp"},
				Env: map[string]string{
					"initialBuckets":    "tesseral-user-content",
					"root":              "containers3root",
					"retainFilesOnExit": "false",
				},
			},
			Started: true,
		},
	)
	testcontainers.CleanupContainer(t, container)
	if err != nil {
		t.Fatalf("failed to start S3 mock container: %v", err)
	}
	endpoint, err := container.PortEndpoint(t.Context(), "9090/tcp", "http")
	if err != nil {
		t.Fatalf("failed to get S3 mock endpoint: %v", err)
	}
	cfg := s3.Options{
		Region:           awsTestRegion,
		EndpointResolver: s3.EndpointResolverFromURL(endpoint),
	}
	return s3.New(cfg)
}
