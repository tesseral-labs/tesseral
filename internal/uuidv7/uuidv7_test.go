package uuidv7_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/tesseral-labs/tesseral/internal/uuidv7"
)

func TestUUIDv7(t *testing.T) {
	tt := []time.Time{
		time.UnixMilli(0),
		time.Now(),
	}
	for _, test := range tt {
		uuid, err := uuidv7.NewWithTime(test)
		assert.NoError(t, err)

		ts := tsFromUUIDv7(uuid)
		assert.Equal(t, ts.UnixMilli(), test.UnixMilli())
	}
}

// tsFromUUIDv7 reconstructs a timestamp from a UUIDv7, accurate to millisecond precision.
func tsFromUUIDv7(id uuid.UUID) time.Time {
	milli := (int64(id[0]) << 40) | (int64(id[1]) << 32) | (int64(id[2]) << 24) | (int64(id[3]) << 16) | (int64(id[4]) << 8) | int64(id[5])
	return time.UnixMilli(milli).UTC()
}
