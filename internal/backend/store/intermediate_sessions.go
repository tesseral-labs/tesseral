package store

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	backendv1 "github.com/openauth/openauth/internal/backend/gen/openauth/backend/v1"
	"github.com/openauth/openauth/internal/backend/projectid"
	"github.com/openauth/openauth/internal/backend/store/queries"
	"github.com/openauth/openauth/internal/store/idformat"
)

func (s *Store) ListIntermediateSessions(ctx context.Context, req *backendv1.ListIntermediateSessionsRequest) (*backendv1.ListIntermediateSessionsResponse, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	var startID uuid.UUID
	if err := s.pageEncoder.Unmarshal(req.PageToken, &startID); err != nil {
		return nil, fmt.Errorf("unmarshal page token: %w", err)
	}

	limit := 10
	qIntermediateSessions, err := q.ListIntermediateSessions(ctx, queries.ListIntermediateSessionsParams{
		ProjectID: projectid.ProjectID(ctx),
		ID:        startID,
		Limit:     int32(limit + 1),
	})
	if err != nil {
		return nil, fmt.Errorf("list intermediate sessions: %w", err)
	}

	var intermediateSessions []*backendv1.IntermediateSession
	for _, qIntermediateSession := range qIntermediateSessions {
		intermediateSessions = append(intermediateSessions, parseIntermediateSession(qIntermediateSession))
	}

	var nextPageToken string
	if len(intermediateSessions) == limit+1 {
		nextPageToken = s.pageEncoder.Marshal(qIntermediateSessions[limit].ID)
		intermediateSessions = intermediateSessions[:limit]
	}

	return &backendv1.ListIntermediateSessionsResponse{
		IntermediateSessions: intermediateSessions,
		NextPageToken:        nextPageToken,
	}, nil
}

func (s *Store) GetIntermediateSession(ctx context.Context, req *backendv1.GetIntermediateSessionRequest) (*backendv1.GetIntermediateSessionResponse, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	intermediateSessionID, err := idformat.IntermediateSession.Parse(req.Id)
	if err != nil {
		return nil, fmt.Errorf("parse intermediate session id: %w", err)
	}

	qIntermediateSession, err := q.GetIntermediateSession(ctx, queries.GetIntermediateSessionParams{
		ProjectID: projectid.ProjectID(ctx),
		ID:        intermediateSessionID,
	})
	if err != nil {
		return nil, fmt.Errorf("get intermediate session: %w", err)
	}

	return &backendv1.GetIntermediateSessionResponse{IntermediateSession: parseIntermediateSession(qIntermediateSession)}, nil
}

func parseIntermediateSession(qIntermediateSession queries.IntermediateSession) *backendv1.IntermediateSession {
	return &backendv1.IntermediateSession{
		Id:                 idformat.IntermediateSession.Format(qIntermediateSession.ID),
		Email:              derefOrEmpty(qIntermediateSession.Email),
		GoogleHostedDomain: derefOrEmpty(qIntermediateSession.GoogleHostedDomain),
		GoogleUserId:       derefOrEmpty(qIntermediateSession.GoogleUserID),
		MicrosoftTenantId:  derefOrEmpty(qIntermediateSession.MicrosoftTenantID),
		MicrosoftUserId:    derefOrEmpty(qIntermediateSession.MicrosoftUserID),
		Revoked:            qIntermediateSession.Revoked,
	}
}
