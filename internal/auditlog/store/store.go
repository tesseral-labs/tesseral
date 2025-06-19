package store

import (
	"time"

	"github.com/tesseral-labs/tesseral/internal/auditlog/store/queries"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Store struct {
	Q *queries.Queries
}

func timestampOrNil(t *time.Time) *timestamppb.Timestamp {
	if t == nil || t.IsZero() {
		return nil
	}
	return timestamppb.New(*t)
}

func derefOrEmpty[T any](t *T) T {
	var z T
	if t == nil {
		return z
	}
	return *t
}
