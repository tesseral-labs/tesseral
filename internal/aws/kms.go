package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
)

type KeyManagementService struct {
	kms *kms.Client
}

type KeyManagementServiceDecryptResult struct {
	KeyID string
	PlainText string
}

type KeyManagementServiceEncryptResult struct {
	KeyID string
	CipherTextBlob []byte
}

func newKeyManagementService(options *kms.Options) *KeyManagementService {
	return &KeyManagementService{
		kms: kms.New(*options),
	}
}

func newKeyManagementServiceFromConfig(cfg *aws.Config) *KeyManagementService {
	return &KeyManagementService{
		kms: kms.NewFromConfig(*cfg),
	}
}

func (k *KeyManagementService) CreateKey(ctx context.Context, params *kms.CreateKeyInput) (*kms.CreateKeyOutput, error) {
	createKeyOutput, err := k.kms.CreateKey(ctx, &kms.CreateKeyInput{
		// TODO: Make sure our description is appropriately set for the project the key belongs to
		Description: aws.String("Example KMS Key with auto-rotation enabled"),
		KeyUsage:    types.KeyUsageTypeEncryptDecrypt,
		CustomerMasterKeySpec: types.CustomerMasterKeySpecSymmetricDefault,
	})
	if err != nil {
		return nil, err
	}

	keyID := *createKeyOutput.KeyMetadata.KeyId
	k.kms.EnableKeyRotation(context.TODO(), &kms.EnableKeyRotationInput{
		KeyId: &keyID,
	})

	return createKeyOutput, nil
}

func (k *KeyManagementService) Decrypt(ctx context.Context, params *kms.DecryptInput) (*KeyManagementServiceDecryptResult, error) {
	decryptOutput, err := k.kms.Decrypt(ctx, params)
	if err != nil {
		return nil, err
	}

	return &KeyManagementServiceDecryptResult{
		KeyID: *decryptOutput.KeyId,
		PlainText: string(decryptOutput.Plaintext),
	}, nil
}

func (k *KeyManagementService) Encrypt(ctx context.Context, params *kms.EncryptInput) (*KeyManagementServiceEncryptResult, error) {
	encryptOutput, err := k.kms.Encrypt(ctx, params)
	if err != nil {
		return &KeyManagementServiceEncryptResult{}, err
	}

	return &KeyManagementServiceEncryptResult{
		KeyID: *encryptOutput.KeyId,
		CipherTextBlob: encryptOutput.CiphertextBlob,
	}, nil
}
