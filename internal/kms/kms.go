package kms

import (
	"context"
	"fmt"

	gcpkms "cloud.google.com/go/kms/apiv1"
	"cloud.google.com/go/kms/apiv1/kmspb"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	awskms "github.com/aws/aws-sdk-go-v2/service/kms"
	awskmstypes "github.com/aws/aws-sdk-go-v2/service/kms/types"
)

type KMS struct {
	config Config

	awskmsClient *awskms.Client
	gcpkmsClient *gcpkms.KeyManagementClient
}

type Config struct {
	Backend string `conf:"backend,noredact"`

	AWSKMSV1KeyID           string `conf:"aws_kms_v1_key_id,noredact"`
	AWSKMSV1KMSBaseEndpoint string `conf:"aws_kms_v1_kms_base_endpoint,noredact"`

	GCPKMSV1KeyName string `conf:"gcp_kms_v1_key_name,noredact"`
}

func New(ctx context.Context, config Config) (*KMS, error) {
	switch config.Backend {
	case "aws_kms_v1":
		awsConfig, err := awsconfig.LoadDefaultConfig(ctx)
		if err != nil {
			return nil, fmt.Errorf("load aws config: %w", err)
		}

		awskmsClient := awskms.NewFromConfig(awsConfig, func(o *awskms.Options) {
			if config.AWSKMSV1KMSBaseEndpoint != "" {
				o.BaseEndpoint = &config.AWSKMSV1KMSBaseEndpoint
			}
		})

		return &KMS{
			config:       config,
			awskmsClient: awskmsClient,
		}, nil
	case "gcp_kms_v1":
		gcpkmsClient, err := gcpkms.NewKeyManagementClient(ctx)
		if err != nil {
			return nil, fmt.Errorf("create gcp kms client: %w", err)
		}

		return &KMS{
			config:       config,
			gcpkmsClient: gcpkmsClient,
		}, nil
	default:
		return nil, fmt.Errorf("unknown backend: %q", config.Backend)
	}
}

func (k *KMS) Encrypt(ctx context.Context, plaintext []byte) ([]byte, error) {
	switch k.config.Backend {
	case "aws_kms_v1":
		encryptRes, err := k.awskmsClient.Encrypt(ctx, &awskms.EncryptInput{
			KeyId:               &k.config.AWSKMSV1KeyID,
			EncryptionAlgorithm: awskmstypes.EncryptionAlgorithmSpecRsaesOaepSha256,
			Plaintext:           plaintext,
		})
		if err != nil {
			return nil, fmt.Errorf("aws kms: encrypt: %w", err)
		}

		return encryptRes.CiphertextBlob, nil
	case "gcp_kms_v1":
		encryptRes, err := k.gcpkmsClient.Encrypt(ctx, &kmspb.EncryptRequest{
			Name:      k.config.GCPKMSV1KeyName,
			Plaintext: plaintext,
		})
		if err != nil {
			return nil, fmt.Errorf("gcp kms: encrypt: %w", err)
		}

		return encryptRes.Ciphertext, nil
	default:
		panic("unreachable")
	}
}

func (k *KMS) Decrypt(ctx context.Context, ciphertext []byte) ([]byte, error) {
	switch k.config.Backend {
	case "aws_kms_v1":
		decryptRes, err := k.awskmsClient.Decrypt(ctx, &awskms.DecryptInput{
			KeyId:               &k.config.AWSKMSV1KeyID,
			EncryptionAlgorithm: awskmstypes.EncryptionAlgorithmSpecRsaesOaepSha256,
			CiphertextBlob:      ciphertext,
		})
		if err != nil {
			return nil, fmt.Errorf("aws kms: decrypt: %w", err)
		}

		return decryptRes.Plaintext, nil
	case "gcp_kms_v1":
		decryptRes, err := k.gcpkmsClient.Decrypt(ctx, &kmspb.DecryptRequest{
			Name:       k.config.GCPKMSV1KeyName,
			Ciphertext: ciphertext,
		})
		if err != nil {
			return nil, fmt.Errorf("gcp kms: decrypt: %w", err)
		}

		return decryptRes.Plaintext, nil
	default:
		panic("unreachable")
	}
}
