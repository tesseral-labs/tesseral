package storetestutil

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func NewS3() (*s3.Client, func()) {
	container, err := testcontainers.GenericContainer(
		context.Background(),
		testcontainers.GenericContainerRequest{
			ContainerRequest: testcontainers.ContainerRequest{
				Image:        "adobe/s3mock:latest",
				ExposedPorts: []string{"9090/tcp"},
				Env: map[string]string{
					"initialBuckets":    "tesseral-user-content",
					"root":              "containers3root",
					"retainFilesOnExit": "false",
				},
				WaitingFor: wait.ForLog("Tomcat started on ports"),
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
		Region:           awsTestRegion,
		EndpointResolver: s3.EndpointResolverFromURL(endpoint),
	}
	return s3.New(cfg), cleanup
}
