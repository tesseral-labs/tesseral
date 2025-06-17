package storetestutil

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StoreDependencies struct {
	DB      *pgxpool.Pool
	KMS     *KMS
	S3      *s3.Client
	Console *Console
}

// NewStoreDependencies initializes and returns test dependencies for testing store layers.
func NewStoreDependencies() (*StoreDependencies, func()) {
	db, cleanupDB := NewDB()
	kms, cleanupKms := NewKMS()
	s3, cleanupS3 := NewS3()
	console := NewConsole(db, kms)

	cleanup := func() {
		cleanupS3()
		cleanupKms()
		cleanupDB()
	}

	return &StoreDependencies{
		DB:      db,
		S3:      s3,
		KMS:     kms,
		Console: console,
	}, cleanup
}
