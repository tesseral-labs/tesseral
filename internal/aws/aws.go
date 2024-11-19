package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
)

type AWSClient struct {
	KMS *KeyManagementService
}

type NewAWSClientParams struct {

}

func NewAWSClient(params *NewAWSClientParams) *AWSClient {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}

	return &AWSClient{
		KMS: newKeyManagementServiceFromConfig(&cfg),
	}
}