package store

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/sha256"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/tesseral-labs/tesseral/internal/backend/authn"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	"github.com/tesseral-labs/tesseral/internal/backend/store/queries"
	"github.com/tesseral-labs/tesseral/internal/common/apierror"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
	"github.com/tesseral-labs/tesseral/internal/ujwt"
	"google.golang.org/protobuf/types/known/structpb"
)

type ActorType string

const (
	ActorTypeUser   ActorType = "user"
	ActorTypeApiKey ActorType = "api_key"
)

func (s *Store) CreateAuditLogEvent(ctx context.Context, req *backendv1.CreateAuditLogEventRequest) (*backendv1.CreateAuditLogEventResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, fmt.Errorf("create audit log event: %w", err)
	}
	defer rollback()

	orgID, err := idformat.Organization.Parse(req.AuditLogEvent.OrganizationId)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid organization id", fmt.Errorf("parse organization id: %w", err))
	}

	projectID := authn.ProjectID(ctx)

	if _, err := q.GetOrganizationByProjectIDAndID(ctx, queries.GetOrganizationByProjectIDAndIDParams{
		ID:        orgID,
		ProjectID: projectID,
	}); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("organization not found", fmt.Errorf("get organization: %w", err))
		}
		return nil, fmt.Errorf("create audit log event: get organization: %w", err)
	}

	// TODO: Feature flag check?

	id, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("create audit log event: failed to create UUID: %w", err)
	}
	eventTime := pgtype.Timestamptz{
		Time:  time.Now(),
		Valid: true,
	}
	if req.AuditLogEvent.Timestamp != nil {
		eventTime.Time = req.AuditLogEvent.Timestamp.AsTime()
	}
	eventName := req.AuditLogEvent.Name
	if eventName == "" {
		return nil, apierror.NewInvalidArgumentError("", errors.New("missing event name"))
	}

	// Resolve the actor type/ID from the given inputs.
	var (
		actorType   ActorType
		actorID     uuid.UUID
		credentials = req.AuditLogEvent.ActorCredentials
		userID      = req.AuditLogEvent.GetUserId()
		apiKeyID    = req.AuditLogEvent.GetApiKeyId()
	)
	switch {
	case credentials != "":
		actorType, actorID, err = s.parseActorCredentials(ctx, credentials, projectID, orgID, q)
		if err != nil {
			return nil, err
		}
	case userID != "":
		actorType = ActorTypeUser
		actorID, err = idformat.User.Parse(userID)
		if err != nil {
			return nil, apierror.NewInvalidArgumentError("invalid user_id", err)
		}
	case apiKeyID != "":
		actorType = ActorTypeApiKey
		actorID, err = idformat.APIKey.Parse(apiKeyID)
		if err != nil {
			return nil, apierror.NewInvalidArgumentError("invalid api_key_id", err)
		}
	default:
		return nil, apierror.NewInvalidArgumentError("", errors.New("either actor_credentials, user_id, or api_key_id must be provided"))
	}

	// Marshal the details to JSON if provided.
	var detailsJSON []byte
	if details := req.AuditLogEvent.Details; details != nil {
		detailsJSON, err = details.MarshalJSON()
		if err != nil {
			return nil, fmt.Errorf("create audit log event: failed to marshal event details JSON: %w", err)
		}
	}
	qEventParams := queries.CreateAuditLogEventParams{
		ID:             id,
		OrganizationID: orgID,
		EventTime:      eventTime,
		EventName:      eventName,
		ActorType:      string(actorType),
		ActorID:        actorID,
		EventDetails:   detailsJSON,
	}
	qEvent, err := q.CreateAuditLogEvent(ctx, qEventParams)
	if err != nil {
		return nil, fmt.Errorf("create audit log event: %w", err)
	}

	event, err := parseAuditLogEvent(qEvent)
	if err != nil {
		return nil, fmt.Errorf("create audit log event: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("create audit log event: %w", err)
	}

	return &backendv1.CreateAuditLogEventResponse{
		AuditLogEvent: event,
	}, nil
}

// parseActorCredentials decodes the given access token or API key and extracts its identifying details.
func (s *Store) parseActorCredentials(
	ctx context.Context,
	credentials string,
	requestProjectID uuid.UUID,
	requestOrgID uuid.UUID,
	q *queries.Queries,
) (ActorType, uuid.UUID, error) {
	projectID := idformat.Project.Format(requestProjectID)
	orgID := idformat.Organization.Format(requestOrgID)

	// Check if it's an access token
	kid, err := ujwt.KeyID(credentials)
	if err == nil {
		publicKeys, err := s.GetSessionPublicKeysByProjectID(ctx, projectID)
		if err != nil {
			return "", uuid.UUID{}, fmt.Errorf("failed to get signing keys for project %q: %w", projectID, err)
		}
		var publicKey *ecdsa.PublicKey
		for _, key := range publicKeys {
			if kid == key.ID {
				publicKey = key.PublicKey
			}
		}
		if publicKey == nil {
			return "", uuid.UUID{}, apierror.NewInvalidArgumentError("invalid actor_credentials", fmt.Errorf("invalid access token: no signing key found for key ID %q in project ID %q", kid, projectID))
		}

		aud := fmt.Sprintf("https://%s.tesseral.app", strings.ReplaceAll(projectID, "_", "-"))
		var claims map[string]any
		if err := ujwt.Claims(publicKey, aud, time.Now(), &claims, credentials); err != nil {
			return "", uuid.UUID{}, apierror.NewInvalidArgumentError("invalid actor_credentials", fmt.Errorf("failed to get claims from access token: %w", err))
		}

		userID, err := idformat.User.Parse(claims["user"].(map[string]any)["id"].(string))
		if err != nil {
			return "", uuid.UUID{}, err
		}
		userOrgID, err := idformat.Organization.Parse(claims["organization"].(map[string]any)["id"].(string))
		if err != nil {
			return "", uuid.UUID{}, err
		}

		if !bytes.Equal(userOrgID[:], requestOrgID[:]) {
			return "", uuid.UUID{}, apierror.NewInvalidArgumentError("invalid actor_credentials", fmt.Errorf("user with ID %q does not belong to organization %q", idformat.User.Format(userID), orgID))
		}
		return ActorTypeUser, userID, nil
	}

	// Ensure it's a valid API key
	secretToken, err := idformat.APIKey.Parse(credentials)
	if err != nil {
		return "", uuid.UUID{}, apierror.NewInvalidArgumentError("invalid actor_credentials", errors.New("actor_credentials must be either a valid API key or user access token"))
	}

	secretTokenSHA := sha256.Sum256(secretToken[:])
	qApiKey, err := q.GetAPIKeyDetailsBySecretTokenSHA256(ctx, queries.GetAPIKeyDetailsBySecretTokenSHA256Params{
		ProjectID:         requestProjectID,
		SecretTokenSha256: secretTokenSHA[:],
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", uuid.UUID{}, apierror.NewInvalidArgumentError("invalid actor_credentials", fmt.Errorf("api key not found"))
		}
		return "", uuid.UUID{}, fmt.Errorf("get api key details by secret token sha256: %w", err)
	}

	if !bytes.Equal(requestOrgID[:], qApiKey.OrganizationID[:]) {
		apiKeyID := idformat.APIKey.Format(qApiKey.ID)
		return "", uuid.UUID{}, apierror.NewInvalidArgumentError("invalid actor_credentials", fmt.Errorf("API Key with ID %q does not belong to organization %q", apiKeyID, orgID))
	}

	return ActorTypeApiKey, qApiKey.ID, nil
}

func parseAuditLogEvent(qEvent queries.OrganizationAuditLogEvent) (*backendv1.AuditLogEvent, error) {
	detailsJSON := qEvent.EventDetails
	var details structpb.Struct
	if err := details.UnmarshalJSON(detailsJSON); err != nil {
		return nil, err
	}

	event := &backendv1.AuditLogEvent{
		Id:             idformat.AuditLogEvent.Format(qEvent.ID),
		OrganizationId: idformat.Organization.Format(qEvent.OrganizationID),
		Name:           qEvent.EventName,
		Timestamp:      timestampOrNil(&qEvent.EventTime.Time),
		Details:        &details,

		ActorCredentials: "",  // input only
		Actor:            nil, // oneof, defined below
	}
	switch ActorType(qEvent.ActorType) {
	case ActorTypeUser:
		event.Actor = &backendv1.AuditLogEvent_UserId{
			UserId: idformat.User.Format(qEvent.ActorID),
		}
	case ActorTypeApiKey:
		event.Actor = &backendv1.AuditLogEvent_ApiKeyId{
			ApiKeyId: idformat.APIKey.Format(qEvent.ActorID),
		}
	default:
		return nil, fmt.Errorf("invalid actor_type: %q", qEvent.ActorType)
	}

	return event, nil
}
