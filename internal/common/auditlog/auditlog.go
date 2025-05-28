package auditlog

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/tesseral-labs/tesseral/internal/common/store/queries"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type Event queries.AuditLogEvent
type EventName string

const (
	CreateSAMLConnectionEventName EventName = "tesseral.saml_connection.create"
	UpdateSAMLConnectionEventName EventName = "tesseral.saml_connection.update"
	DeleteSAMLConnectionEventName EventName = "tesseral.saml_connection.delete"
)

type EventData struct {
	ProjectID      uuid.UUID
	OrganizationID *uuid.UUID
	UserID         *uuid.UUID
	SessionID      *uuid.UUID
	ApiKeyID       *uuid.UUID
	EventName      EventName
	ResourceBefore proto.Message
	ResourceAfter  proto.Message
}

func NewEvent(data EventData) (event Event, err error) {
	var (
		beforeBytes []byte
		afterBytes  []byte
	)
	if data.ResourceBefore != nil {
		beforeBytes, err = protojson.Marshal(data.ResourceBefore)
		if err != nil {
			return
		}
	}
	if data.ResourceAfter != nil {
		afterBytes, err = protojson.Marshal(data.ResourceAfter)
		if err != nil {
			return
		}
	}
	details := struct {
		Before json.RawMessage `json:"before,omitempty"`
		After  json.RawMessage `json:"after,omitempty"`
	}{
		Before: beforeBytes,
		After:  afterBytes,
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
