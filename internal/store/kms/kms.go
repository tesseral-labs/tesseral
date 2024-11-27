package kms

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsKms "github.com/aws/aws-sdk-go-v2/service/kms"
)

type KeyManagementService struct {
	kms *awsKms.Client
}

type KeyManagementServiceDecryptResult struct {
	KeyID string
	Value []byte
}

type KeyManagementServiceEncryptResult struct {
	KeyID          string
	CipherTextBlob []byte
}

func NewKeyManagementServiceFromConfig(cfg *aws.Config, endpoint *string) *KeyManagementService {
	return &KeyManagementService{
		kms: awsKms.NewFromConfig(*cfg, func(o *awsKms.Options) {
			if endpoint != nil {
				o.BaseEndpoint = endpoint
			}
		}),
	}
}

func (k *KeyManagementService) Decrypt(ctx context.Context, params *awsKms.DecryptInput) (*KeyManagementServiceDecryptResult, error) {
	decryptOutput, err := k.kms.Decrypt(ctx, params)
	if err != nil {
		return nil, err
	}

	return &KeyManagementServiceDecryptResult{
		KeyID: *decryptOutput.KeyId,
		Value: decryptOutput.Plaintext,
	}, nil
}

func (k *KeyManagementService) Encrypt(ctx context.Context, params *awsKms.EncryptInput) (*KeyManagementServiceEncryptResult, error) {
	encryptOutput, err := k.kms.Encrypt(ctx, params)
	if err != nil {
		return &KeyManagementServiceEncryptResult{}, err
	}

	return &KeyManagementServiceEncryptResult{
		KeyID:          *encryptOutput.KeyId,
		CipherTextBlob: encryptOutput.CiphertextBlob,
	}, nil
}
