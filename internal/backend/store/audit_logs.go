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
	"github.com/tesseral-labs/tesseral/internal/prettysecret"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
	"github.com/tesseral-labs/tesseral/internal/uuidv7"
	"google.golang.org/protobuf/types/known/structpb"
)

func (s *Store) CreateCustomAuditLogEvent(ctx context.Context, req *backendv1.CreateAuditLogEventRequest) (*backendv1.CreateAuditLogEventResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	projectID := authn.ProjectID(ctx)

	eventTime := time.Now().UTC()
	if req.AuditLogEvent.EventTime != nil {
		eventTime = req.AuditLogEvent.EventTime.AsTime()
	}
	// Generate the UUIDv7 based on the event time.
	eventID, err := uuidv7.NewWithTime(eventTime)
	if err != nil {
		return nil, fmt.Errorf("failed to create UUID: %w", err)
	}
	eventName := req.AuditLogEvent.EventName
	if eventName == "" {
		return nil, apierror.NewInvalidArgumentError("", errors.New("missing event name"))
	} else {
		if err := validateEventName(eventName); err != nil {
			return nil, apierror.NewInvalidArgumentError("invalide event name", fmt.Errorf("validate event name: %w", err))
		}
	}

	// Resolve the actor type/ID from the given inputs.
	var (
		orgID       = req.AuditLogEvent.OrganizationId
		orgUUID     *uuid.UUID
		userID      = req.AuditLogEvent.UserId
		userUUID    *uuid.UUID
		sessionID   = req.AuditLogEvent.SessionId
		sessionUUID *uuid.UUID
		apiKeyID    = req.AuditLogEvent.ApiKeyId
		apiKeyUUID  *uuid.UUID
	)

	if req.AuditLogEvent.Credential != "" {
		if isJWTFormat(req.AuditLogEvent.Credential) {
			parsedAccessToken, err := parseAccessTokenNoValidate(req.AuditLogEvent.Credential)
			if err != nil {
				return nil, apierror.NewInvalidArgumentError("invalid credential", fmt.Errorf("parse access token: %w", err))
			}

			if parsedAccessToken.OrganizationID != nil {
				orgID = idformat.Organization.Format(*parsedAccessToken.OrganizationID)
				orgUUID = parsedAccessToken.OrganizationID
			}

			if parsedAccessToken.UserID != nil {
				userID = idformat.User.Format(*parsedAccessToken.UserID)
				userUUID = parsedAccessToken.UserID
			}

			if parsedAccessToken.SessionID != nil {
				sessionID = idformat.Session.Format(*parsedAccessToken.SessionID)
				sessionUUID = parsedAccessToken.SessionID
			}
		} else if isAPIKeyFormat(req.AuditLogEvent.Credential) {
			qProject, err := q.GetProjectByID(ctx, authn.ProjectID(ctx))
			if err != nil {
				return nil, fmt.Errorf("get project by id: %w", err)
			}

			if qProject.ApiKeySecretTokenPrefix == nil {
				return nil, apierror.NewPermissionDeniedError("api key secret token prefix is not set for this project", fmt.Errorf("api key secret token prefix not set for project"))
			}

			secretTokenBytes, err := prettysecret.Parse(*qProject.ApiKeySecretTokenPrefix, req.AuditLogEvent.Credential)
			if err != nil {
				return nil, apierror.NewUnauthenticatedApiKeyError("malformed_api_key_secret_token", fmt.Errorf("parse secret token: %w", err))
			}
			secretTokenSHA256 := sha256.Sum256(secretTokenBytes[:])

			qApiKeyDetails, err := q.GetAPIKeyDetailsBySecretTokenSHA256(ctx, queries.GetAPIKeyDetailsBySecretTokenSHA256Params{
				SecretTokenSha256: secretTokenSHA256[:],
				ProjectID:         authn.ProjectID(ctx),
			})
			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					return nil, apierror.NewInvalidArgumentError("invalid credential", fmt.Errorf("get api key details: %w", err))
				}

				return nil, fmt.Errorf("get api key details: %w", err)
			}

			apiKeyID = idformat.APIKey.Format(qApiKeyDetails.ID)
			apiKeyUUID = (*uuid.UUID)(&qApiKeyDetails.ID)
			orgID = idformat.Organization.Format(qApiKeyDetails.OrganizationID)
			orgUUID = (*uuid.UUID)(&qApiKeyDetails.OrganizationID)
		}
	}

	if orgID == "" && userID == "" && sessionID == "" && apiKeyID == "" {
		return nil, apierror.NewInvalidArgumentError("", errors.New("either organization_id, user_id, session_id, or api_key_id must be provided"))
	}
	if sessionID != "" {
		id, err := idformat.Session.Parse(sessionID)
		if err != nil {
			return nil, apierror.NewInvalidArgumentError("invalid session_id", err)
		}
		// Lookup session, user, and organization
		session, err := q.GetSession(ctx, queries.GetSessionParams{
			ID:        id,
			ProjectID: projectID,
		})
		if err != nil {
			return nil, apierror.NewInvalidArgumentError("invalid session_id", fmt.Errorf("get session: %w", err))
		}
		sessionUUID = (*uuid.UUID)(&session.ID)
		userID = idformat.User.Format(session.UserID)
	}
	if userID != "" {
		id, err := idformat.User.Parse(userID)
		if err != nil {
			return nil, apierror.NewInvalidArgumentError("invalid user_id", err)
		}
		// Lookup user and organization
		user, err := q.GetUser(ctx, queries.GetUserParams{
			ID:        id,
			ProjectID: projectID,
		})
		if err != nil {
			return nil, apierror.NewInvalidArgumentError("invalid user_id", fmt.Errorf("get user: %w", err))
		}
		userUUID = (*uuid.UUID)(&user.ID)
		orgID = idformat.Organization.Format(user.OrganizationID)
	}
	if apiKeyID != "" {
		id, err := idformat.APIKey.Parse(apiKeyID)
		if err != nil {
			return nil, apierror.NewInvalidArgumentError("invalid api_key_id", err)
		}
		// Lookup API key and organization
		apiKey, err := q.GetAPIKeyByID(ctx, queries.GetAPIKeyByIDParams{
			ID:        id,
			ProjectID: projectID,
		})
		if err != nil {
			return nil, apierror.NewInvalidArgumentError("invalid api_key_id", fmt.Errorf("get api key: %w", err))
		}
		apiKeyUUID = (*uuid.UUID)(&apiKey.ID)
		orgID = idformat.Organization.Format(apiKey.OrganizationID)
	}
	if orgID != "" {
		id, err := idformat.Organization.Parse(orgID)
		if err != nil {
			return nil, apierror.NewInvalidArgumentError("invalid organization_id", err)
		}
		_, err = q.GetOrganizationByProjectIDAndID(ctx, queries.GetOrganizationByProjectIDAndIDParams{
			ID:        id,
			ProjectID: projectID,
		})
		if err != nil {
			return nil, apierror.NewInvalidArgumentError("invalid organization_id", fmt.Errorf("get organization: %w", err))
		}
		orgUUID = (*uuid.UUID)(&id)
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
		OrganizationID: orgUUID,
		UserID:         userUUID,
		SessionID:      sessionUUID,
		ApiKeyID:       apiKeyUUID,
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
	EventDetails map[string]any
	ResourceType queries.AuditLogEventResourceType
	ResourceID   *uuid.UUID
}

func (s *Store) logAuditEvent(ctx context.Context, q *queries.Queries, data logAuditEventParams) (queries.AuditLogEvent, error) {
	// Generate the UUIDv7 based on the event time.
	eventTime := time.Now().UTC()
	eventID, err := uuidv7.NewWithTime(eventTime)
	if err != nil {
		return queries.AuditLogEvent{}, fmt.Errorf("failed to create UUID: %w", err)
	}

	eventDetailsBytes, err := json.Marshal(data.EventDetails)
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

func parseAuditLogEvent(qEvent queries.AuditLogEvent) (*backendv1.AuditLogEvent, error) {
	eventDetailsJSON := qEvent.EventDetails
	var eventDetails structpb.Struct
	if err := eventDetails.UnmarshalJSON(eventDetailsJSON); err != nil {
		return nil, err
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

		resourceType *string
		resourceID   *string
	)
	if orgUUID := qEvent.OrganizationID; orgUUID != nil {
		organizationID = refOrNil(idformat.Organization.Format(*orgUUID))
	}
	if userUUID := qEvent.UserID; userUUID != nil {
		userID = refOrNil(idformat.User.Format(*userUUID))
	}
	if sessionUUID := qEvent.SessionID; sessionUUID != nil {
		sessionID = refOrNil(idformat.Session.Format(*sessionUUID))
	}
	if apiKeyUUID := qEvent.ApiKeyID; apiKeyUUID != nil {
		apiKeyID = refOrNil(idformat.APIKey.Format(*apiKeyUUID))
	}
	if dogfoodUserUUID := qEvent.DogfoodUserID; dogfoodUserUUID != nil {
		dogfoodUserID = refOrNil(idformat.Session.Format(*dogfoodUserUUID))
	}
	if dogfoodSessionUUID := qEvent.DogfoodSessionID; dogfoodSessionUUID != nil {
		dogfoodSessionID = refOrNil(idformat.Session.Format(*dogfoodSessionUUID))
	}
	if backendApiKeyUUID := qEvent.BackendApiKeyID; backendApiKeyUUID != nil {
		backendApiKeyID = refOrNil(idformat.BackendAPIKey.Format(*backendApiKeyUUID))
	}
	if intermediateSessionUUID := qEvent.IntermediateSessionID; intermediateSessionUUID != nil {
		intermediateSessionID = refOrNil(idformat.IntermediateSession.Format(*intermediateSessionUUID))
	}
	if resourceTypeDto := qEvent.ResourceType; resourceTypeDto != nil {
		resourceType = (*string)(resourceTypeDto)
	}
	if resourceUUID := qEvent.ResourceID; resourceUUID != nil {
		switch *qEvent.ResourceType {
		case queries.AuditLogEventResourceTypeAction:
			resourceID = refOrNil(resourceUUID.String()) // TODO: Actions don't have an ID format
		case queries.AuditLogEventResourceTypeApiKey:
			resourceID = refOrNil(idformat.APIKey.Format(*resourceUUID))
		case queries.AuditLogEventResourceTypeApiKeyRoleAssignment:
			resourceID = refOrNil(idformat.APIKeyRoleAssignment.Format(*resourceUUID))
		case queries.AuditLogEventResourceTypeAuditLogEvent:
			resourceID = refOrNil(idformat.AuditLogEvent.Format(*resourceUUID))
		case queries.AuditLogEventResourceTypeBackendApiKey:
			resourceID = refOrNil(idformat.BackendAPIKey.Format(*resourceUUID))
		case queries.AuditLogEventResourceTypeEmailVerificationChallenge:
			resourceID = refOrNil(idformat.EmailVerificationChallenge.Format(*resourceUUID))
		case queries.AuditLogEventResourceTypeIntermediateSession:
			resourceID = refOrNil(idformat.IntermediateSession.Format(*resourceUUID))
		case queries.AuditLogEventResourceTypeOrganization:
			resourceID = refOrNil(idformat.Organization.Format(*resourceUUID))
		case queries.AuditLogEventResourceTypeOrganizationGoogleHostedDomains:
			resourceID = nil // Google hosted domains are not represented as a resource ID
		case queries.AuditLogEventResourceTypeOrganizationMicrosoftTenantIds:
			resourceID = nil // Microsoft tenant IDs are not represented as a resource ID
		case queries.AuditLogEventResourceTypePasskey:
			resourceID = refOrNil(idformat.Passkey.Format(*resourceUUID))
		case queries.AuditLogEventResourceTypePasswordResetCode:
			resourceID = refOrNil(idformat.PasswordResetCode.Format(*resourceUUID))
		case queries.AuditLogEventResourceTypeProject:
			resourceID = refOrNil(idformat.Project.Format(*resourceUUID))
		case queries.AuditLogEventResourceTypeProjectUiSettings:
			resourceID = refOrNil(idformat.ProjectUISettings.Format(*resourceUUID))
		case queries.AuditLogEventResourceTypeProjectWebhookSettings:
			resourceID = refOrNil(idformat.ProjectWebhookSettings.Format(*resourceUUID))
		case queries.AuditLogEventResourceTypePublishableKey:
			resourceID = refOrNil(idformat.PublishableKey.Format(*resourceUUID))
		case queries.AuditLogEventResourceTypeRole:
			resourceID = refOrNil(idformat.Role.Format(*resourceUUID))
		case queries.AuditLogEventResourceTypeSamlConnection:
			resourceID = refOrNil(idformat.SAMLConnection.Format(*resourceUUID))
		case queries.AuditLogEventResourceTypeScimApiKey:
			resourceID = refOrNil(idformat.SCIMAPIKey.Format(*resourceUUID))
		case queries.AuditLogEventResourceTypeUser:
			resourceID = refOrNil(idformat.User.Format(*resourceUUID))
		case queries.AuditLogEventResourceTypeUserAuthenticatorAppChallenge:
			resourceID = refOrNil(idformat.AuthenticatorAppRecoveryCode.Format(*resourceUUID))
		case queries.AuditLogEventResourceTypeUserImpersonationToken:
			resourceID = refOrNil(idformat.UserImpersonationToken.Format(*resourceUUID))
		case queries.AuditLogEventResourceTypeUserInvite:
			resourceID = refOrNil(idformat.UserInvite.Format(*resourceUUID))
		case queries.AuditLogEventResourceTypeUserRoleAssignment:
			resourceID = refOrNil(idformat.UserRoleAssignment.Format(*resourceUUID))
		}
	}

	return &backendv1.AuditLogEvent{
		Id:                    idformat.AuditLogEvent.Format(qEvent.ID),
		OrganizationId:        derefOrEmpty(organizationID),
		UserId:                derefOrEmpty(userID),
		SessionId:             derefOrEmpty(sessionID),
		ApiKeyId:              derefOrEmpty(apiKeyID),
		DogfoodUserId:         derefOrEmpty(dogfoodUserID),
		DogfoodSessionId:      derefOrEmpty(dogfoodSessionID),
		BackendApiKeyId:       derefOrEmpty(backendApiKeyID),
		IntermediateSessionId: derefOrEmpty(intermediateSessionID),
		ResourceType:          derefOrEmpty(resourceType),
		ResourceId:            derefOrEmpty(resourceID),
		EventName:             qEvent.EventName,
		EventTime:             timestampOrNil(qEvent.EventTime),
		EventDetails:          &eventDetails,
	}, nil
}

var eventNamePattern = regexp.MustCompile(`^[a-z0-9_]+\.[a-z0-9_]+\.[a-z0-9_]+`)

func validateEventName(action string) error {
	if !eventNamePattern.MatchString(action) {
		return apierror.NewInvalidArgumentError("event names must be of the form x.y.z, only containing a-z0-9_", nil)
	}
	if strings.HasPrefix(action, "tesseral") {
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

type accessTokenDetails struct {
	ImpersonatorEmail string
	OrganizationID    *uuid.UUID
	SessionID         *uuid.UUID
	UserID            *uuid.UUID
}

func parseAccessTokenNoValidate(accessToken string) (*accessTokenDetails, error) {
	jwtParts := strings.Split(accessToken, ".")
	if len(jwtParts) != 3 {
		return nil, apierror.NewInvalidArgumentError("invalid credential", fmt.Errorf("invalid access token format: expected 3 parts, got %d", len(jwtParts)))
	}

	payloadSegment := jwtParts[1]
	decoded, err := base64.RawURLEncoding.DecodeString(payloadSegment)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid credential", fmt.Errorf("failed to decode access token payload: %w", err))
	}

	var details accessTokenDetails

	var claims map[string]interface{}
	if err := json.Unmarshal(decoded, &claims); err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid credential", fmt.Errorf("failed to unmarshal access token claims: %w", err))
	}

	if claims["session"] != nil && claims["session"].(map[string]interface{})["id"] != nil {
		sessionID, err := uuid.Parse(claims["session"].(map[string]interface{})["id"].(string))
		if err != nil {
			return nil, apierror.NewInvalidArgumentError("invalid credential", fmt.Errorf("failed to parse session ID: %w", err))
		}

		details.SessionID = &sessionID
	}

	if claims["user"] != nil && claims["user"].(map[string]interface{})["id"] != nil {
		userID, err := uuid.Parse(claims["user"].(map[string]interface{})["id"].(string))
		if err != nil {
			return nil, apierror.NewInvalidArgumentError("invalid credential", fmt.Errorf("failed to parse user ID: %w", err))
		}

		details.UserID = &userID
	}

	if claims["organization"] != nil && claims["organization"].(map[string]interface{})["id"] != nil {
		orgID, err := uuid.Parse(claims["organization"].(map[string]interface{})["id"].(string))
		if err != nil {
			return nil, apierror.NewInvalidArgumentError("invalid credential", fmt.Errorf("failed to parse organization ID: %w", err))
		}

		details.OrganizationID = &orgID
	}

	if claims["impersonator"] != nil && claims["impersonator"].(map[string]interface{})["email"] != nil {
		details.ImpersonatorEmail = claims["impersonator"].(map[string]interface{})["email"].(string)
	}

	return &details, nil
}
