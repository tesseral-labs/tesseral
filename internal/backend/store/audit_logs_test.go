package store

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUUIDv7(t *testing.T) {
	tt := []time.Time{
		time.UnixMilli(0),
		time.Now(),
	}
	for _, test := range tt {
		uuid, err := makeUUIDv7(test)
		assert.NoError(t, err)

		ts := tsFromUUIDv7(uuid)
		assert.Equal(t, ts.UnixMilli(), test.UnixMilli())
		assert.Equalf(t, time.UTC, ts.Location(), "should be UTC time")
	}
}
