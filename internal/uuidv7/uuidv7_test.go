package uuidv7

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUUIDv7(t *testing.T) {
	testCases := []struct {
		name string
		now  time.Time
		id   uuid.UUID
		want uuid.UUID
	}{
		{
			name: "make against nil uuid",
			now:  time.Unix(1749684795, 0),
			id:   uuid.Nil,
			want: uuid.MustParse("01976157-3678-7000-8000-000000000000"),
		},
		{
			name: "make against max uuid",
			now:  time.Unix(1749684795, 0),
			id:   uuid.Max,
			want: uuid.MustParse("01976157-3678-7fff-bfff-ffffffffffff"),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got := makeUUIDv7(tt.now, tt.id)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.now.Unix(), tsFromUUIDv7(got).Unix())
		})
	}
}

// tsFromUUIDv7 reconstructs a timestamp from a UUIDv7, accurate to millisecond precision.
func tsFromUUIDv7(id uuid.UUID) time.Time {
	milli := (int64(id[0]) << 40) | (int64(id[1]) << 32) | (int64(id[2]) << 24) | (int64(id[3]) << 16) | (int64(id[4]) << 8) | int64(id[5])
	return time.UnixMilli(milli).UTC()
}
