package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/tesseral-labs/tesseral/internal/backend/authn"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	"github.com/tesseral-labs/tesseral/internal/backend/store/queries"
	"github.com/tesseral-labs/tesseral/internal/common/apierror"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
	"github.com/tesseral-labs/tesseral/internal/uuidv7"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

var (
	titleCaser = cases.Title(language.English)
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
	}

	// Resolve the actor type/ID from the given inputs.
	var (
		orgID       = req.AuditLogEvent.GetOrganizationId()
		orgUUID     *uuid.UUID
		userID      = req.AuditLogEvent.GetUserId()
		userUUID    *uuid.UUID
		sessionID   = req.AuditLogEvent.GetSessionId()
		sessionUUID *uuid.UUID
		apiKeyID    = req.AuditLogEvent.GetApiKeyId()
		apiKeyUUID  *uuid.UUID
	)
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

type AuditLogEventData struct {
	OrganizationID *uuid.UUID

	ResourceType     queries.AuditLogEventResourceType
	ResourceID       uuid.UUID
	Resource         proto.Message
	PreviousResource proto.Message

	// For example, `create`, `update`, `delete`, etc.
	EventType string
}

func (s *Store) CreateTesseralAuditLogEvent(ctx context.Context, data AuditLogEventData) (queries.AuditLogEvent, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return queries.AuditLogEvent{}, err
	}
	defer rollback()

	resourceType := string(data.ResourceType)
	jsonResourceType := snakeToCamelCase(resourceType)

	pluralResourceType := resourceType
	if !strings.HasSuffix(resourceType, "s") {
		pluralResourceType += "s"
	}
	// e.g. "tesseral.projects.create"
	eventName := fmt.Sprintf("tesseral.%s.%s", pluralResourceType, data.EventType)

	details := make(map[string]any)
	if data.PreviousResource != nil {
		previousResourceBytes, err := protojson.Marshal(data.PreviousResource)
		if err != nil {
			return queries.AuditLogEvent{}, err
		}
		details[fmt.Sprintf("previous%s", capitalizeFirstLetter(jsonResourceType))] = json.RawMessage(previousResourceBytes)
	}
	if data.Resource != nil {
		resourceBytes, err := protojson.Marshal(data.Resource)
		if err != nil {
			return queries.AuditLogEvent{}, err
		}
		details[jsonResourceType] = json.RawMessage(resourceBytes)
	}
	detailsBytes, err := json.Marshal(details)
	if err != nil {
		return queries.AuditLogEvent{}, err
	}

	// Generate the UUIDv7 based on the event time.
	eventTime := time.Now().UTC()
	eventID, err := uuidv7.NewWithTime(eventTime)
	if err != nil {
		return queries.AuditLogEvent{}, fmt.Errorf("failed to create UUID: %w", err)
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
		ResourceID:       &data.ResourceID,
		EventName:        eventName,
		EventTime:        &eventTime,
		EventDetails:     detailsBytes,
	}

	qEvent, err := q.CreateAuditLogEvent(ctx, qEventParams)
	if err != nil {
		return queries.AuditLogEvent{}, err
	}

	if err := commit(); err != nil {
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
		OrganizationId:        organizationID,
		UserId:                userID,
		SessionId:             sessionID,
		ApiKeyId:              apiKeyID,
		DogfoodUserId:         dogfoodUserID,
		DogfoodSessionId:      dogfoodSessionID,
		BackendApiKeyId:       backendApiKeyID,
		IntermediateSessionId: intermediateSessionID,
		ResourceType:          resourceType,
		ResourceId:            resourceID,
		EventName:             qEvent.EventName,
		EventTime:             timestampOrNil(qEvent.EventTime),
		EventDetails:          &eventDetails,
	}, nil
}

func capitalizeFirstLetter(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func snakeToCamelCase(s string) string {
	parts := strings.Split(s, "_")
	for i, part := range parts {
		if i > 0 {
			parts[i] = titleCaser.String(part)
		} else {
			parts[i] = strings.ToLower(part)
		}
	}
	return strings.Join(parts, "")
}
