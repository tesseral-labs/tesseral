package kms

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsKms "github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
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

func NewKeyManagementServiceFromConfig(cfg *aws.Config) *KeyManagementService {
	return &KeyManagementService{
		kms: awsKms.NewFromConfig(*cfg),
	}
}

func (k *KeyManagementService) CreateKey(ctx context.Context, params *awsKms.CreateKeyInput) (*awsKms.CreateKeyOutput, error) {
	createKeyOutput, err := k.kms.CreateKey(ctx, &awsKms.CreateKeyInput{
		// TODO: Make sure our description is appropriately set for the project the key belongs to
		Description:           aws.String("Example KMS Key with auto-rotation enabled"),
		KeyUsage:              types.KeyUsageTypeEncryptDecrypt,
		CustomerMasterKeySpec: types.CustomerMasterKeySpecSymmetricDefault,
	})
	if err != nil {
		return nil, err
	}

	keyID := *createKeyOutput.KeyMetadata.KeyId
	if _, err := k.kms.EnableKeyRotation(context.TODO(), &awsKms.EnableKeyRotationInput{
		KeyId: &keyID,
	}); err != nil {
		return nil, err
	}

	return createKeyOutput, nil
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
