package store

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/tesseral-labs/tesseral/internal/backend/authn"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	"github.com/tesseral-labs/tesseral/internal/backend/store/queries"
	"github.com/tesseral-labs/tesseral/internal/common/apierror"
	commonv1 "github.com/tesseral-labs/tesseral/internal/common/gen/tesseral/common/v1"
	"github.com/tesseral-labs/tesseral/internal/prettysecret"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
	"github.com/tesseral-labs/tesseral/internal/uuidv7"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Store) CreateCustomAuditLogEvent(ctx context.Context, req *backendv1.CreateAuditLogEventRequest) (*backendv1.CreateAuditLogEventResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	if err := enforceSingleActor(req); err != nil {
		return nil, fmt.Errorf("enforce single actor: %w", err)
	}

	eventTime := time.Now()
	if req.AuditLogEvent.EventTime != nil {
		eventTime = req.AuditLogEvent.EventTime.AsTime()
	}

	// Generate the UUIDv7 based on the event time.
	eventID := uuidv7.NewWithTime(eventTime)

	eventName := req.AuditLogEvent.EventName
	if err := validateEventName(eventName); err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid event name", fmt.Errorf("validate event name: %w", err))
	}

	deriveEventContextRes, err := deriveEventContextForRequest(ctx, q, req)
	if err != nil {
		return nil, fmt.Errorf("derive event context for actor: %w", err)
	}

	// Marshal the details to JSON if provided.
	var eventDetails []byte
	if req.AuditLogEvent.EventDetails != nil {
		b, err := req.AuditLogEvent.EventDetails.MarshalJSON()
		if err != nil {
			return nil, fmt.Errorf("create audit log event: failed to marshal event details JSON: %w", err)
		}
		eventDetails = b
	}

	qEvent, err := q.CreateAuditLogEvent(ctx, queries.CreateAuditLogEventParams{
		ID:             eventID,
		ProjectID:      authn.ProjectID(ctx),
		OrganizationID: deriveEventContextRes.OrganizationID,
		ActorUserID:    deriveEventContextRes.UserID,
		ActorSessionID: deriveEventContextRes.SessionID,
		ActorApiKeyID:  deriveEventContextRes.APIKeyID,
		EventName:      eventName,
		EventTime:      &eventTime,
		EventDetails:   eventDetails,
	})
	if err != nil {
		return nil, fmt.Errorf("create audit log event: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	return &backendv1.CreateAuditLogEventResponse{
		AuditLogEvent: parseAuditLogEvent(qEvent),
	}, nil
}

func parseAuditLogEvent(qAuditLogEvent queries.AuditLogEvent) *backendv1.AuditLogEvent {
	var eventDetails structpb.Struct
	if err := protojson.Unmarshal(qAuditLogEvent.EventDetails, &eventDetails); err != nil {
		panic(fmt.Errorf("unmarshal event details: %w", err))
	}

	var organizationID string
	if qAuditLogEvent.OrganizationID != nil {
		organizationID = idformat.Organization.Format(*qAuditLogEvent.OrganizationID)
	}

	var userID string
	if qAuditLogEvent.ActorUserID != nil {
		userID = idformat.User.Format(*qAuditLogEvent.ActorUserID)
	}

	var sessionID string
	if qAuditLogEvent.ActorSessionID != nil {
		sessionID = idformat.Session.Format(*qAuditLogEvent.ActorSessionID)
	}

	var apiKeyID string
	if qAuditLogEvent.ActorApiKeyID != nil {
		apiKeyID = idformat.APIKey.Format(*qAuditLogEvent.ActorApiKeyID)
	}

	var backendApiKeyID string
	if qAuditLogEvent.ActorBackendApiKeyID != nil {
		backendApiKeyID = idformat.BackendAPIKey.Format(*qAuditLogEvent.ActorBackendApiKeyID)
	}

	var intermediateSessionID string
	if qAuditLogEvent.ActorIntermediateSessionID != nil {
		intermediateSessionID = idformat.IntermediateSession.Format(*qAuditLogEvent.ActorIntermediateSessionID)
	}

	return &backendv1.AuditLogEvent{
		Id:                         idformat.AuditLogEvent.Format(qAuditLogEvent.ID),
		OrganizationId:             organizationID,
		ActorUserId:                userID,
		ActorSessionId:             sessionID,
		ActorApiKeyId:              apiKeyID,
		ActorBackendApiKeyId:       backendApiKeyID,
		ActorIntermediateSessionId: intermediateSessionID,
		EventName:                  qAuditLogEvent.EventName,
		EventTime:                  timestamppb.New(*qAuditLogEvent.EventTime),
		EventDetails:               &eventDetails,
	}
}

var eventNamePattern = regexp.MustCompile(`^[a-z0-9_]+\.[a-z0-9_]+\.[a-z0-9_]+`)

func validateEventName(eventName string) error {
	if !eventNamePattern.MatchString(eventName) {
		return apierror.NewInvalidArgumentError("event names must be of the form x.y.z, only containing a-z0-9_", nil)
	}
	if strings.HasPrefix(eventName, "tesseral") {
		return apierror.NewInvalidArgumentError("event names must not start with 'tesseral'", nil)
	}
	return nil
}

var (
	jwtRegex    = regexp.MustCompile(`^[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+$`)
	apiKeyRegex = regexp.MustCompile(`^[a-z0-9_]+$`)
)

func parseAccessTokenNoValidate(accessToken string) (*commonv1.AccessTokenData, error) {
	jwtParts := strings.Split(accessToken, ".")
	if len(jwtParts) != 3 {
		return nil, apierror.NewInvalidArgumentError("invalid credential", fmt.Errorf("invalid access token format: expected 3 parts, got %d", len(jwtParts)))
	}

	payloadSegment := jwtParts[1]
	decoded, err := base64.RawURLEncoding.DecodeString(payloadSegment)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid credential", fmt.Errorf("failed to decode access token payload: %w", err))
	}

	var data commonv1.AccessTokenData

	if err := json.Unmarshal(decoded, &data); err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid credential", fmt.Errorf("failed to unmarshal access token claims: %w", err))
	}

	return &data, nil
}

func enforceSingleActor(req *backendv1.CreateAuditLogEventRequest) error {
	actorCount := 0
	if req.AuditLogEvent.OrganizationId != "" {
		actorCount++
	}
	if req.AuditLogEvent.ActorUserId != "" {
		actorCount++
	}
	if req.AuditLogEvent.ActorSessionId != "" {
		actorCount++
	}
	if req.AuditLogEvent.ActorApiKeyId != "" {
		actorCount++
	}
	if req.AuditLogEvent.ActorCredentials != "" {
		actorCount++
	}

	if actorCount != 1 {
		return apierror.NewInvalidArgumentError("exactly one of actor_credentials, actor_api_key_id, actor_session_id, actor_user_id, or organization_id must be provided", nil)
	}

	return nil
}

type deriveEventContextForRequestResponse struct {
	OrganizationID *uuid.UUID
	UserID         *uuid.UUID
	SessionID      *uuid.UUID
	APIKeyID       *uuid.UUID
}

func deriveEventContextForRequest(ctx context.Context, q *queries.Queries, req *backendv1.CreateAuditLogEventRequest) (*deriveEventContextForRequestResponse, error) {
	switch {
	case jwtRegex.MatchString(req.AuditLogEvent.ActorCredentials):
		// If ActorCredentials looks like a jwt, parse it as such. We
		// deliberately don't validate the JWT here; in the special case of
		// ActorCredentials in CreateAuditLogEvent in the Backend API, this JWT
		// is a convenience shorthand, not an authentication scheme.
		//
		// Moreover, we don't want callers to have to deal with the possibility
		// of an audit log call failing because the user's JWT expired after
		// some auditable work was done, but before calling CreateAuditLogEvent.
		parsedAccessToken, err := parseAccessTokenNoValidate(req.AuditLogEvent.ActorCredentials)
		if err != nil {
			return nil, apierror.NewInvalidArgumentError("invalid credential", fmt.Errorf("parse access token: %w", err))
		}

		orgID, err := idformat.Organization.Parse(parsedAccessToken.Organization.Id)
		if err != nil {
			return nil, apierror.NewInvalidArgumentError("invalid organization_id in credential", fmt.Errorf("parse organization id: %w", err))
		}

		userID, err := idformat.User.Parse(parsedAccessToken.User.Id)
		if err != nil {
			return nil, apierror.NewInvalidArgumentError("invalid user_id in credential", fmt.Errorf("parse user id: %w", err))
		}

		sessionID, err := idformat.Session.Parse(parsedAccessToken.Session.Id)
		if err != nil {
			return nil, apierror.NewInvalidArgumentError("invalid session_id in credential", fmt.Errorf("parse session id: %w", err))
		}

		return &deriveEventContextForRequestResponse{
			OrganizationID: (*uuid.UUID)(&orgID),
			UserID:         (*uuid.UUID)(&userID),
			SessionID:      (*uuid.UUID)(&sessionID),
		}, nil
	case req.AuditLogEvent.ActorCredentials != "" && apiKeyRegex.MatchString(req.AuditLogEvent.ActorCredentials):
		// If ActorCredentials looks like an API key, parse it as such.
		qProject, err := q.GetProjectByID(ctx, authn.ProjectID(ctx))
		if err != nil {
			return nil, fmt.Errorf("get project by id: %w", err)
		}

		if qProject.ApiKeySecretTokenPrefix == nil {
			return nil, apierror.NewPermissionDeniedError("api key secret token prefix is not set for this project", fmt.Errorf("api key secret token prefix not set for project"))
		}

		secretTokenBytes, err := prettysecret.Parse(*qProject.ApiKeySecretTokenPrefix, req.AuditLogEvent.ActorCredentials)
		if err != nil {
			return nil, apierror.NewUnauthenticatedApiKeyError("malformed_api_key_secret_token", fmt.Errorf("parse secret token: %w", err))
		}
		secretTokenSHA256 := sha256.Sum256(secretTokenBytes[:])

		qApiKeyDetails, err := q.GetAPIKeyDetailsBySecretTokenSHA256(ctx, queries.GetAPIKeyDetailsBySecretTokenSHA256Params{
			ProjectID:         authn.ProjectID(ctx),
			SecretTokenSha256: secretTokenSHA256[:],
		})
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, apierror.NewInvalidArgumentError("invalid credential", fmt.Errorf("get api key details: %w", err))
			}

			return nil, fmt.Errorf("get api key details: %w", err)
		}

		return &deriveEventContextForRequestResponse{
			OrganizationID: &qApiKeyDetails.OrganizationID,
			APIKeyID:       &qApiKeyDetails.ID,
		}, nil
	case req.AuditLogEvent.ActorSessionId != "":
		parsedSessionID, err := idformat.Session.Parse(req.AuditLogEvent.ActorSessionId)
		if err != nil {
			return nil, apierror.NewInvalidArgumentError("invalid session_id", fmt.Errorf("parse session id: %w", err))
		}

		eventContext, err := q.DeriveAuditLogEventContextForSessionID(ctx, queries.DeriveAuditLogEventContextForSessionIDParams{
			ProjectID: authn.ProjectID(ctx),
			ID:        parsedSessionID,
		})
		if err != nil {
			return nil, apierror.NewInvalidArgumentError("invalid session_id", fmt.Errorf("derive audit log event context from session id: %w", err))
		}

		return &deriveEventContextForRequestResponse{
			OrganizationID: refOrNil(eventContext.OrganizationID),
			UserID:         refOrNil(eventContext.UserID),
			SessionID:      (*uuid.UUID)(&parsedSessionID),
		}, nil
	case req.AuditLogEvent.ActorUserId != "":
		parsedUserID, err := idformat.User.Parse(req.AuditLogEvent.ActorUserId)
		if err != nil {
			return nil, apierror.NewInvalidArgumentError("invalid user_id", fmt.Errorf("parse user id: %w", err))
		}

		eventContext, err := q.DeriveAuditLogEventContextForUserID(ctx, queries.DeriveAuditLogEventContextForUserIDParams{
			ProjectID: authn.ProjectID(ctx),
			ID:        parsedUserID,
		})
		if err != nil {
			return nil, apierror.NewInvalidArgumentError("invalid user_id", fmt.Errorf("derive audit log event context from user id: %w", err))
		}

		return &deriveEventContextForRequestResponse{
			OrganizationID: refOrNil(eventContext.OrganizationID),
			UserID:         (*uuid.UUID)(&parsedUserID),
		}, nil
	case req.AuditLogEvent.ActorApiKeyId != "":
		parsedApiKeyID, err := idformat.APIKey.Parse(req.AuditLogEvent.ActorApiKeyId)
		if err != nil {
			return nil, apierror.NewInvalidArgumentError("invalid api_key_id", fmt.Errorf("parse api key id: %w", err))
		}

		eventContext, err := q.DeriveAuditLogEventContextForAPIKeyID(ctx, queries.DeriveAuditLogEventContextForAPIKeyIDParams{
			ProjectID: authn.ProjectID(ctx),
			ID:        parsedApiKeyID,
		})
		if err != nil {
			return nil, apierror.NewInvalidArgumentError("invalid api_key_id", fmt.Errorf("derive audit log event context from api key id: %w", err))
		}

		return &deriveEventContextForRequestResponse{
			OrganizationID: refOrNil(eventContext.OrganizationID),
			APIKeyID:       (*uuid.UUID)(&parsedApiKeyID),
		}, nil
	case req.AuditLogEvent.OrganizationId != "":
		orgID, err := idformat.Organization.Parse(req.AuditLogEvent.OrganizationId)
		if err != nil {
			return nil, apierror.NewInvalidArgumentError("invalid organization_id", fmt.Errorf("parse organization id: %w", err))
		}

		// Ensure the organization exists in the project.
		if _, err := q.GetOrganizationByProjectIDAndID(ctx, queries.GetOrganizationByProjectIDAndIDParams{
			ProjectID: authn.ProjectID(ctx),
			ID:        orgID,
		}); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, apierror.NewInvalidArgumentError("organization not found", fmt.Errorf("get organization by id: %w", err))
			}
			return nil, fmt.Errorf("get organization by id: %w", err)
		}

		return &deriveEventContextForRequestResponse{
			OrganizationID: (*uuid.UUID)(&orgID),
		}, nil
	}

	// If no actor/organization could be derived, then we say that this is a
	// free-floating audit log event for the project.
	return &deriveEventContextForRequestResponse{}, nil
}
