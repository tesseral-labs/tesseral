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
)

func (s *Store) CreateCustomAuditLogEvent(ctx context.Context, req *backendv1.CreateAuditLogEventRequest) (*backendv1.CreateAuditLogEventResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	fmt.Println("🚨 request:", req)

	projectID := authn.ProjectID(ctx)

	if err := enforceSingleActor(req); err != nil {
		return nil, apierror.NewInvalidArgumentError("exactly one of organizationId, userId, sessionId, or apiKeyId must be provided", fmt.Errorf("enforce single actor: %w", err))
	}

	eventTime := time.Now()
	if req.AuditLogEvent.EventTime != nil {
		eventTime = req.AuditLogEvent.EventTime.AsTime()
	}
	// Generate the UUIDv7 based on the event time.
	eventID := uuidv7.NewWithTime(eventTime)
	eventName := req.AuditLogEvent.EventName
	if eventName == "" {
		return nil, apierror.NewInvalidArgumentError("", errors.New("missing event name"))
	}
	if err := validateEventName(eventName); err != nil {
		return nil, apierror.NewInvalidArgumentError("invalide event name", fmt.Errorf("validate event name: %w", err))
	}

	orgID, userID, sessionID, apiKeyID, err := deriveEventContextForRequest(ctx, q, req)
	if err != nil {
		return nil, fmt.Errorf("derive event context for actor: %w", err)
	}

	// Marshal the details to JSON if provided.
	var eventDetailsJSON []byte
	if eventDetails := req.AuditLogEvent.EventDetails; eventDetails != nil {
		json, err := eventDetails.MarshalJSON()
		if err != nil {
			return nil, fmt.Errorf("create audit log event: failed to marshal event details JSON: %w", err)
		}
		eventDetailsJSON = json
	}

	qEvent, err := q.CreateAuditLogEvent(ctx, queries.CreateAuditLogEventParams{
		ID:             eventID,
		ProjectID:      projectID,
		OrganizationID: orgID,
		UserID:         userID,
		SessionID:      sessionID,
		ApiKeyID:       apiKeyID,
		EventName:      eventName,
		EventTime:      &eventTime,
		EventDetails:   eventDetailsJSON,
	})
	if err != nil {
		return nil, fmt.Errorf("create audit log event: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	pEvent, err := parseAuditLogEvent(qEvent)
	if err != nil {
		return nil, fmt.Errorf("parse audit log event: %w", err)
	}

	return &backendv1.CreateAuditLogEventResponse{
		AuditLogEvent: pEvent,
	}, nil
}

type logAuditEventParams struct {
	OrganizationID *uuid.UUID

	EventName    string
	EventDetails *structpb.Value
	ResourceType queries.AuditLogEventResourceType
	ResourceID   *uuid.UUID
}

func (s *Store) logAuditEvent(ctx context.Context, q *queries.Queries, data logAuditEventParams) (queries.AuditLogEvent, error) {
	// Generate the UUIDv7 based on the event time.
	eventTime := time.Now()
	eventID := uuidv7.NewWithTime(eventTime)

	eventDetailsBytes, err := protojson.Marshal(data.EventDetails)
	if err != nil {
		return queries.AuditLogEvent{}, fmt.Errorf("failed to marshal event details: %w", err)
	}

	var (
		dogfoodUserID    *uuid.UUID
		dogfoodSessionID *uuid.UUID
		backendApiKeyID  *uuid.UUID
	)
	contextData := authn.GetContextData(ctx)
	switch {
	case contextData.ProjectAPIKey != nil:
		backendApiKeyUUID, err := idformat.BackendAPIKey.Parse(contextData.ProjectAPIKey.BackendAPIKeyID)
		if err != nil {
			return queries.AuditLogEvent{}, fmt.Errorf("parse backend api key id: %w", err)
		}
		backendApiKeyID = (*uuid.UUID)(&backendApiKeyUUID)
	case contextData.DogfoodSession != nil:
		dogfoodUserUUID, err := idformat.User.Parse(contextData.DogfoodSession.UserID)
		if err != nil {
			return queries.AuditLogEvent{}, fmt.Errorf("parse dogfood session user id: %w", err)
		}
		dogfoodUserID = (*uuid.UUID)(&dogfoodUserUUID)
		dogfoodSessionUUID, err := idformat.Session.Parse(contextData.DogfoodSession.SessionID)
		if err != nil {
			return queries.AuditLogEvent{}, fmt.Errorf("parse dogfood session project id: %w", err)
		}
		dogfoodSessionID = (*uuid.UUID)(&dogfoodSessionUUID)
	}

	qEventParams := queries.CreateAuditLogEventParams{
		ID:               eventID,
		ProjectID:        authn.ProjectID(ctx),
		OrganizationID:   data.OrganizationID,
		DogfoodUserID:    dogfoodUserID,
		DogfoodSessionID: dogfoodSessionID,
		BackendApiKeyID:  backendApiKeyID,
		ResourceType:     refOrNil(data.ResourceType),
		ResourceID:       data.ResourceID,
		EventName:        data.EventName,
		EventTime:        &eventTime,
		EventDetails:     eventDetailsBytes,
	}

	qEvent, err := q.CreateAuditLogEvent(ctx, qEventParams)
	if err != nil {
		return queries.AuditLogEvent{}, err
	}

	return qEvent, nil
}

func parseAuditLogEvent(qAuditLogEvent queries.AuditLogEvent) (*backendv1.AuditLogEvent, error) {
	var eventDetails structpb.Struct
	if err := protojson.Unmarshal(qAuditLogEvent.EventDetails, &eventDetails); err != nil {
		return nil, fmt.Errorf("unmarshal event details: %w", err)
	}

	var (
		organizationID        *string
		userID                *string
		sessionID             *string
		apiKeyID              *string
		dogfoodUserID         *string
		dogfoodSessionID      *string
		backendApiKeyID       *string
		intermediateSessionID *string
	)
	if orgUUID := qAuditLogEvent.OrganizationID; orgUUID != nil {
		organizationID = refOrNil(idformat.Organization.Format(*orgUUID))
	}
	if userUUID := qAuditLogEvent.UserID; userUUID != nil {
		userID = refOrNil(idformat.User.Format(*userUUID))
	}
	if sessionUUID := qAuditLogEvent.SessionID; sessionUUID != nil {
		sessionID = refOrNil(idformat.Session.Format(*sessionUUID))
	}
	if apiKeyUUID := qAuditLogEvent.ApiKeyID; apiKeyUUID != nil {
		apiKeyID = refOrNil(idformat.APIKey.Format(*apiKeyUUID))
	}
	if dogfoodUserUUID := qAuditLogEvent.DogfoodUserID; dogfoodUserUUID != nil {
		dogfoodUserID = refOrNil(idformat.Session.Format(*dogfoodUserUUID))
	}
	if dogfoodSessionUUID := qAuditLogEvent.DogfoodSessionID; dogfoodSessionUUID != nil {
		dogfoodSessionID = refOrNil(idformat.Session.Format(*dogfoodSessionUUID))
	}
	if backendApiKeyUUID := qAuditLogEvent.BackendApiKeyID; backendApiKeyUUID != nil {
		backendApiKeyID = refOrNil(idformat.BackendAPIKey.Format(*backendApiKeyUUID))
	}
	if intermediateSessionUUID := qAuditLogEvent.IntermediateSessionID; intermediateSessionUUID != nil {
		intermediateSessionID = refOrNil(idformat.IntermediateSession.Format(*intermediateSessionUUID))
	}

	return &backendv1.AuditLogEvent{
		Id:                    idformat.AuditLogEvent.Format(qAuditLogEvent.ID),
		OrganizationId:        derefOrEmpty(organizationID),
		UserId:                derefOrEmpty(userID),
		SessionId:             derefOrEmpty(sessionID),
		ApiKeyId:              derefOrEmpty(apiKeyID),
		DogfoodUserId:         derefOrEmpty(dogfoodUserID),
		DogfoodSessionId:      derefOrEmpty(dogfoodSessionID),
		BackendApiKeyId:       derefOrEmpty(backendApiKeyID),
		IntermediateSessionId: derefOrEmpty(intermediateSessionID),
		EventName:             qAuditLogEvent.EventName,
		EventTime:             timestampOrNil(qAuditLogEvent.EventTime),
		EventDetails:          &eventDetails,
	}, nil
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

func isJWTFormat(value string) bool {
	return jwtRegex.MatchString(value)
}

func isAPIKeyFormat(value string) bool {
	return apiKeyRegex.MatchString(value)
}

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
	if req.AuditLogEvent.UserId != "" {
		actorCount++
	}
	if req.AuditLogEvent.SessionId != "" {
		actorCount++
	}
	if req.AuditLogEvent.ApiKeyId != "" {
		actorCount++
	}
	if req.AuditLogEvent.Credentials != "" {
		actorCount++
	}

	if actorCount != 1 {
		return fmt.Errorf("exactly one of organizationId, userId, sessionId, or apiKeyId must be provided")
	}

	return nil
}

func deriveEventContextForRequest(ctx context.Context, q *queries.Queries, req *backendv1.CreateAuditLogEventRequest) (orgID, userID, sessionID, apiKeyID *uuid.UUID, err error) {
	projectID := authn.ProjectID(ctx)

	if req.AuditLogEvent.Credentials != "" {
		if isJWTFormat(req.AuditLogEvent.Credentials) {
			parsedAccessToken, err := parseAccessTokenNoValidate(req.AuditLogEvent.Credentials)
			if err != nil {
				return nil, nil, nil, nil, apierror.NewInvalidArgumentError("invalid credential", fmt.Errorf("parse access token: %w", err))
			}

			parsedOrgID, err := idformat.Organization.Parse(parsedAccessToken.Organization.Id)
			if err != nil {
				return nil, nil, nil, nil, apierror.NewInvalidArgumentError("invalid organization_id in credential", fmt.Errorf("parse organization id: %w", err))
			}
			orgID = (*uuid.UUID)(&parsedOrgID)

			parsedUserID, err := idformat.User.Parse(parsedAccessToken.User.Id)
			if err != nil {
				return nil, nil, nil, nil, apierror.NewInvalidArgumentError("invalid user_id in credential", fmt.Errorf("parse user id: %w", err))
			}
			userID = (*uuid.UUID)(&parsedUserID)

			parsedSessionID, err := idformat.Session.Parse(parsedAccessToken.Session.Id)
			if err != nil {
				return nil, nil, nil, nil, apierror.NewInvalidArgumentError("invalid session_id in credential", fmt.Errorf("parse session id: %w", err))
			}
			sessionID = (*uuid.UUID)(&parsedSessionID)
		} else if isAPIKeyFormat(req.AuditLogEvent.Credentials) {
			qProject, err := q.GetProjectByID(ctx, authn.ProjectID(ctx))
			if err != nil {
				return nil, nil, nil, nil, fmt.Errorf("get project by id: %w", err)
			}

			if qProject.ApiKeySecretTokenPrefix == nil {
				return nil, nil, nil, nil, apierror.NewPermissionDeniedError("api key secret token prefix is not set for this project", fmt.Errorf("api key secret token prefix not set for project"))
			}

			secretTokenBytes, err := prettysecret.Parse(*qProject.ApiKeySecretTokenPrefix, req.AuditLogEvent.Credentials)
			if err != nil {
				return nil, nil, nil, nil, apierror.NewUnauthenticatedApiKeyError("malformed_api_key_secret_token", fmt.Errorf("parse secret token: %w", err))
			}
			secretTokenSHA256 := sha256.Sum256(secretTokenBytes[:])

			qApiKeyDetails, err := q.GetAPIKeyDetailsBySecretTokenSHA256(ctx, queries.GetAPIKeyDetailsBySecretTokenSHA256Params{
				SecretTokenSha256: secretTokenSHA256[:],
				ProjectID:         authn.ProjectID(ctx),
			})
			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					return nil, nil, nil, nil, apierror.NewInvalidArgumentError("invalid credential", fmt.Errorf("get api key details: %w", err))
				}

				return nil, nil, nil, nil, fmt.Errorf("get api key details: %w", err)
			}

			apiKeyID = (*uuid.UUID)(&qApiKeyDetails.ID)
			orgID = (*uuid.UUID)(&qApiKeyDetails.OrganizationID)
		}
	}

	if req.AuditLogEvent.SessionId != "" {
		parsedSessionID, err := idformat.Session.Parse(req.AuditLogEvent.SessionId)
		if err != nil {
			return nil, nil, nil, nil, apierror.NewInvalidArgumentError("invalid session_id", fmt.Errorf("parse session id: %w", err))
		}

		eventContext, err := q.DeriveAuditLogEventContextForSessionID(ctx, queries.DeriveAuditLogEventContextForSessionIDParams{
			ID:        parsedSessionID,
			ProjectID: projectID,
		})
		if err != nil {
			return nil, nil, nil, nil, apierror.NewInvalidArgumentError("invalid session_id", fmt.Errorf("derive audit log event context from session id: %w", err))
		}

		sessionID = (*uuid.UUID)(&parsedSessionID)
		orgID = refOrNil(eventContext.OrganizationID)
		userID = refOrNil(eventContext.UserID)
	}

	if req.AuditLogEvent.UserId != "" {
		parsedUserID, err := idformat.User.Parse(req.AuditLogEvent.UserId)
		if err != nil {
			return nil, nil, nil, nil, apierror.NewInvalidArgumentError("invalid user_id", fmt.Errorf("parse user id: %w", err))
		}

		eventContext, err := q.DeriveAuditLogEventContextForUserID(ctx, queries.DeriveAuditLogEventContextForUserIDParams{
			ID:        parsedUserID,
			ProjectID: projectID,
		})
		if err != nil {
			return nil, nil, nil, nil, apierror.NewInvalidArgumentError("invalid user_id", fmt.Errorf("derive audit log event context from user id: %w", err))
		}
		userID = (*uuid.UUID)(&parsedUserID)
		orgID = refOrNil(eventContext.OrganizationID)
	}

	if req.AuditLogEvent.ApiKeyId != "" {
		parsedApiKeyID, err := idformat.APIKey.Parse(req.AuditLogEvent.ApiKeyId)
		if err != nil {
			return nil, nil, nil, nil, apierror.NewInvalidArgumentError("invalid api_key_id", fmt.Errorf("parse api key id: %w", err))
		}

		eventContext, err := q.DeriveAuditLogEventContextForAPIKeyID(ctx, queries.DeriveAuditLogEventContextForAPIKeyIDParams{
			ID:        parsedApiKeyID,
			ProjectID: projectID,
		})
		if err != nil {
			return nil, nil, nil, nil, apierror.NewInvalidArgumentError("invalid api_key_id", fmt.Errorf("derive audit log event context from api key id: %w", err))
		}
		apiKeyID = (*uuid.UUID)(&parsedApiKeyID)
		orgID = refOrNil(eventContext.OrganizationID)
	}

	if req.AuditLogEvent.OrganizationId != "" {
		parsedOrgID, err := idformat.Organization.Parse(req.AuditLogEvent.OrganizationId)
		if err != nil {
			return nil, nil, nil, nil, apierror.NewInvalidArgumentError("invalid organization_id", fmt.Errorf("parse organization id: %w", err))
		}

		// Ensure the organization exists in the project.
		if _, err := q.GetOrganizationByProjectIDAndID(ctx, queries.GetOrganizationByProjectIDAndIDParams{
			ID:        derefOrEmpty((*uuid.UUID)(&parsedOrgID)),
			ProjectID: projectID,
		}); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, nil, nil, nil, apierror.NewInvalidArgumentError("organization_id not found", fmt.Errorf("get organization by id: %w", err))
			}
			return nil, nil, nil, nil, fmt.Errorf("get organization by id: %w", err)
		}

		orgID = (*uuid.UUID)(&parsedOrgID)
	}

	return
}
