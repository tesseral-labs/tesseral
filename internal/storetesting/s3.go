package storetesting

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type testS3 struct {
	Client                *s3.Client
	UserContentBucketName string
}

func newS3() (*testS3, func()) {
	const bucketName = "tesseral-user-content"
	container, err := testcontainers.GenericContainer(
		context.Background(),
		testcontainers.GenericContainerRequest{
			ContainerRequest: testcontainers.ContainerRequest{
				Image:        "adobe/s3mock:latest",
				ExposedPorts: []string{"9090/tcp"},
				Env: map[string]string{
					"initialBuckets":    bucketName,
					"root":              "containers3root",
					"retainFilesOnExit": "false",
				},
				WaitingFor: wait.ForLog("Started S3MockApplication"),
			},
			Started: true,
		},
	)
	cleanup := func() {
		_ = testcontainers.TerminateContainer(container)
	}
	if err != nil {
		cleanup()
		log.Panicf("failed to start S3 mock container: %v", err)
	}
	endpoint, err := container.PortEndpoint(context.Background(), "9090/tcp", "http")
	if err != nil {
		cleanup()
		log.Panicf("failed to get S3 mock endpoint: %v", err)
	}
	cfg := s3.Options{
		Region:       awsTestRegion,
		BaseEndpoint: &endpoint,
		UsePathStyle: true,
		Credentials:  credentials.NewStaticCredentialsProvider("foo", "bar", ""),
	}
	return &testS3{
		Client:                s3.New(cfg),
		UserContentBucketName: bucketName,
	}, cleanup
}
