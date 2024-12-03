package store

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/openauth/openauth/internal/store/idformat"
	"github.com/openauth/openauth/internal/store/queries"
)

type VerifiedEmail struct {
	ID              string
	ProjectID       string
	CreateTime      time.Time
	Email           string
	GoogleUserID    string
	MicrosoftUserID string
}

type CreateVerifiedEmailParams struct {
	ProjectID       string
	Email           string
	GoogleUserID    string
	MicrosoftUserID string
}

func (s *Store) CreateVerifiedEmail(ctx context.Context, params *CreateVerifiedEmailParams) (*VerifiedEmail, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	projectID, err := uuid.Parse(params.ProjectID)
	if err != nil {
		return nil, err
	}

	verifiedEmail, err := q.CreateVerifiedEmail(ctx, queries.CreateVerifiedEmailParams{
		ID:              uuid.New(),
		ProjectID:       projectID,
		Email:           params.Email,
		GoogleUserID:    &params.GoogleUserID,
		MicrosoftUserID: &params.MicrosoftUserID,
	})
	if err != nil {
		return nil, err
	}

	if err := commit(); err != nil {
		return nil, err
	}

	return parseVerfifiedEmail(&verifiedEmail), nil
}

func parseVerfifiedEmail(v *queries.VerifiedEmail) *VerifiedEmail {
	return &VerifiedEmail{
		ID:              idformat.VerifiedEmail.Format(v.ID),
		ProjectID:       idformat.Project.Format(v.ProjectID),
		CreateTime:      *v.CreateTime,
		Email:           v.Email,
		GoogleUserID:    *v.GoogleUserID,
		MicrosoftUserID: *v.MicrosoftUserID,
	}
}
