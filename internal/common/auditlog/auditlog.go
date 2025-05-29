package auditlog

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/tesseral-labs/tesseral/internal/common/store/queries"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type Event queries.AuditLogEvent

type TesseralEventData struct {
	ProjectID        uuid.UUID
	OrganizationID   *uuid.UUID
	UserID           *uuid.UUID
	SessionID        *uuid.UUID
	ApiKeyID         *uuid.UUID
	EventName        string
	ResourceName     string
	Resource         proto.Message
	PreviousResource proto.Message
}

var (
	titleCaser = cases.Title(language.English)
)

func NewTesseralEvent(data TesseralEventData) (Event, error) {
	details := make(map[string]any)
	if data.PreviousResource != nil {
		previousResourceBytes, err := protojson.Marshal(data.PreviousResource)
		if err != nil {
			return Event{}, err
		}
		details[fmt.Sprintf("previous%s", titleCaser.String(data.ResourceName))] = json.RawMessage(previousResourceBytes)
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
