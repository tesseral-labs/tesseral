package store

import (
	backendv1 "github.com/openauth/openauth/internal/backend/gen/openauth/backend/v1"
	"github.com/openauth/openauth/internal/backend/store/queries"
	"github.com/openauth/openauth/internal/store/idformat"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func parseProjectAPIKey(qProjectAPIKey queries.ProjectApiKey) *backendv1.ProjectAPIKey {
	return &backendv1.ProjectAPIKey{
		Id:          idformat.ProjectAPIKey.Format(qProjectAPIKey.ID),
		ProjectId:   idformat.Project.Format(qProjectAPIKey.ProjectID),
		CreateTime:  timestamppb.New(*qProjectAPIKey.CreateTime),
		Revoked:     qProjectAPIKey.Revoked,
		SecretToken: "", // intentionally left blank
	}
}
