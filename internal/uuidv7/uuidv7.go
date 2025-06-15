package uuidv7

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// NewWithTime creates a new UUID at the given time.
//
// Allowing control over the time bits (as opposed to using the current time)
// means that the bits typically reserved for a sequence number are left random.
func NewWithTime(ts time.Time) uuid.UUID {
	id, err := uuid.NewRandom()
	if err != nil {
		panic(fmt.Errorf("generate random uuid: %w", err))
	}
	return makeUUIDv7(ts, id)
}

// makeUUIDv7 copies google/uuid for constructing a UUIDv7 at a point in time
// given a base UUID data.
func makeUUIDv7(ts time.Time, id uuid.UUID) uuid.UUID {
	nano := ts.UnixNano()
	const nanoPerMilli = 1_000_000
	milli := nano / nanoPerMilli

	// Sequence number is not used since there is no accurate way to establish one.
	// Instead we leave the random bits from the V4 in place.

	/*
		 0                   1                   2                   3
		 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
		+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
		|                           unix_ts_ms                          |
		+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
		|          unix_ts_ms           |  ver  |  rand_a (12 bit seq)  |
		+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
		|var|                        rand_b                             |
		+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
		|                            rand_b                             |
		+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	*/
	_ = id[15] // bounds check

	id[0] = byte(milli >> 40)
	id[1] = byte(milli >> 32)
	id[2] = byte(milli >> 24)
	id[3] = byte(milli >> 16)
	id[4] = byte(milli >> 8)
	id[5] = byte(milli)

	id[6] = 0x70 | (0x0F & id[6]) // Version is 7 (0b0111)
	id[8] = 0x80 | (0x3F & id[8]) // Variant is 0b10

	return id
}
