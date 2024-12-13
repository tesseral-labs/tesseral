package projectid

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/openauth/openauth/internal/store/idformat"
)

type ctxData struct {
	projectID uuid.UUID
}

type ctxKey struct{}

var ErrProjectIDHeaderRequired = errors.New("X-TODO-OpenAuth-Project-ID header is required")

func NewHttpHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		projectIDHeader := "project_0cyqkcnidpxvcq4hr3gazl1tu" // r.Header.Get("X-TODO-OpenAuth-Project-ID")
		if projectIDHeader == "" {
			http.Error(w, ErrProjectIDHeaderRequired.Error(), http.StatusBadRequest)
			return
		}

		projectID, err := idformat.Project.Parse(projectIDHeader)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), ctxKey{}, ctxData{
			projectID,
		})
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func ProjectID(ctx context.Context) uuid.UUID {
	v, ok := ctx.Value(ctxKey{}).(ctxData)
	if !ok {
		panic(errors.New("ctx does not carry project ID data"))
	}

	return v.projectID
}
