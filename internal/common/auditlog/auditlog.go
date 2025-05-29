package auditlog

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/tesseral-labs/tesseral/internal/common/store/queries"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type Event queries.AuditLogEvent

// TesseralEventName represents the name of a first-party event type.
type TesseralEventName string

const (
	CreateAPIKeyEventName               TesseralEventName = "tesseral.api_keys.create"
	UpdateAPIKeyEventName               TesseralEventName = "tesseral.api_keys.update"
	DeleteAPIKeyEventName               TesseralEventName = "tesseral.api_keys.delete"
	RevokeAPIKeyEventName               TesseralEventName = "tesseral.api_keys.revoke"
	CreateAPIKeyRoleAssignmentEventName TesseralEventName = "tesseral.api_key_role_assignments.create"
	DeleteAPIKeyRoleAssignmentEventName TesseralEventName = "tesseral.api_key_role_assignments.delete"
	UpdateGoogleHostedDomainsEventName  TesseralEventName = "tesseral.google_hosted_domains.update"
	UpdateMicrosoftTenantIDsEventName   TesseralEventName = "tesseral.microsoft_tenant_ids.update"
	UpdateOrganizationEventName         TesseralEventName = "tesseral.organizations.update"
	CreatePasskeyEventName              TesseralEventName = "tesseral.passkeys.create"
	DeletePasskeyEventName              TesseralEventName = "tesseral.passkeys.delete"
	CreateRoleEventName                 TesseralEventName = "tesseral.roles.create"
	UpdateRoleEventName                 TesseralEventName = "tesseral.roles.update"
	DeleteRoleEventName                 TesseralEventName = "tesseral.roles.delete"
	CreateSAMLConnectionEventName       TesseralEventName = "tesseral.saml_connections.create"
	UpdateSAMLConnectionEventName       TesseralEventName = "tesseral.saml_connections.update"
	DeleteSAMLConnectionEventName       TesseralEventName = "tesseral.saml_connections.delete"
	CreateSCIMAPIKeyEventName           TesseralEventName = "tesseral.scim_api_keys.create"
	UpdateSCIMAPIKeyEventName           TesseralEventName = "tesseral.scim_api_keys.update"
	DeleteSCIMAPIKeyEventName           TesseralEventName = "tesseral.scim_api_keys.delete"
	RevokeSCIMAPIKeyEventName           TesseralEventName = "tesseral.scim_api_keys.revoke"
	CreateUserEventName                 TesseralEventName = "tesseral.users.create"
	UpdateUserEventName                 TesseralEventName = "tesseral.users.update"
	DeleteUserEventName                 TesseralEventName = "tesseral.users.delete"
	CreateUserInviteEventName           TesseralEventName = "tesseral.user_invites.create"
	DeleteUserInviteEventName           TesseralEventName = "tesseral.user_invites.delete"
	CreateUserRoleAssignmentEventName   TesseralEventName = "tesseral.user_role_assignments.create"
	DeleteUserRoleAssignmentEventName   TesseralEventName = "tesseral.user_role_assignments.delete"
)

type TesseralEventData struct {
	ProjectID        uuid.UUID
	OrganizationID   *uuid.UUID
	UserID           *uuid.UUID
	SessionID        *uuid.UUID
	ApiKeyID         *uuid.UUID
	EventName        TesseralEventName
	ResourceName     string
	Resource         proto.Message
	PreviousResource proto.Message
}

func NewTesseralEvent(data TesseralEventData) (Event, error) {
	details := make(map[string]any)
	if data.PreviousResource != nil {
		previousResourceBytes, err := protojson.Marshal(data.PreviousResource)
		if err != nil {
			return Event{}, err
		}
		details[fmt.Sprintf("previous_%s", data.ResourceName)] = json.RawMessage(previousResourceBytes)
	}
	if data.Resource != nil {
		resourceBytes, err := protojson.Marshal(data.Resource)
		if err != nil {
			return Event{}, err
		}
		details[data.ResourceName] = json.RawMessage(resourceBytes)
	}
	detailsBytes, err := json.Marshal(details)
	if err != nil {
		return Event{}, err
	}
	return Event{
		ProjectID:      data.ProjectID,
		OrganizationID: data.OrganizationID,
		UserID:         data.UserID,
		SessionID:      data.SessionID,
		ApiKeyID:       data.ApiKeyID,
		EventName:      string(data.EventName),
		EventDetails:   detailsBytes,
	}, nil
}
