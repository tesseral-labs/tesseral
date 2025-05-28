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
type EventName string

const (
	CreateSAMLConnectionEventName EventName = "tesseral.saml_connections.create"
	UpdateSAMLConnectionEventName EventName = "tesseral.saml_connections.update"
	DeleteSAMLConnectionEventName EventName = "tesseral.saml_connections.delete"
)

type EventData struct {
	ProjectID        uuid.UUID
	OrganizationID   *uuid.UUID
	UserID           *uuid.UUID
	SessionID        *uuid.UUID
	ApiKeyID         *uuid.UUID
	EventName        EventName
	ResourceName     string
	Resource         proto.Message
	PreviousResource proto.Message
}

func NewEvent(data EventData) (Event, error) {
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
