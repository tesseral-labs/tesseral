package auditlog

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/tesseral-labs/tesseral/internal/common/store/queries"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

type Event queries.AuditLogEvent
type EventName string

const (
	CreateUserEventName   EventName = "tesseral:create_user"
	UpdateUserEventName   EventName = "tesseral:update_user"
	LoginAttemptEventName EventName = "tesseral:login_attempt"
)

type UserEventData struct {
	ProjectID uuid.UUID
	User      UserData
}

type userEventDetails struct {
	User userDetails `json:"user"`
}

func (data UserEventData) details() userEventDetails {
	return userEventDetails{
		User: data.User.details(),
	}
}

type UserData struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	Email          string
}

type userDetails struct {
	ID             string `json:"id"`
	OrganizationID string `json:"organization_id"`
	Email          string `json:"email"`
}

func (data UserData) details() userDetails {
	return userDetails{
		ID:             idformat.User.Format(data.ID),
		OrganizationID: idformat.Organization.Format(data.OrganizationID),
		Email:          data.Email,
	}
}

func NewCreateUserEvent(data UserEventData) (Event, error) {
	details := data.details()
	detailsJSON, err := json.Marshal(details)
	if err != nil {
		return Event{}, err
	}
	return Event{
		ProjectID:      data.ProjectID,
		OrganizationID: &data.User.OrganizationID,
		UserID:         &data.User.ID,
		EventName:      string(CreateUserEventName),
		EventDetails:   detailsJSON,
	}, nil
}

func NewUpdateUserEvent(data UserEventData) (Event, error) {
	details := data.details()
	detailsJSON, err := json.Marshal(details)
	if err != nil {
		return Event{}, err
	}
	return Event{
		ProjectID:      data.ProjectID,
		OrganizationID: &data.User.OrganizationID,
		UserID:         &data.User.ID,
		EventName:      string(UpdateUserEventName),
		EventDetails:   detailsJSON,
	}, nil
}

type LoginAttemptEventData struct {
	ProjectID             uuid.UUID
	User                  *UserData
	OrganizationID        uuid.UUID
	IntermediateSessionID uuid.UUID
	SessionID             *uuid.UUID
	Success               bool
}

type loginAttemptEventDetails struct {
	User                  *userDetails `json:"user,omitempty"`
	IntermediateSessionID string       `json:"intermediate_session_id"`
	Success               bool         `json:"success"`
}

func (data LoginAttemptEventData) details() loginAttemptEventDetails {
	details := loginAttemptEventDetails{
		IntermediateSessionID: idformat.IntermediateSession.Format(data.IntermediateSessionID),
		Success:               data.Success,
	}
	if user := data.User; user != nil {
		userDetails := user.details()
		details.User = &userDetails
	}
	return details
}

func NewLoginAttemptEvent(data LoginAttemptEventData) (Event, error) {
	details := data.details()
	detailsJSON, err := json.Marshal(details)
	if err != nil {
		return Event{}, err
	}
	return Event{
		ProjectID:      data.ProjectID,
		OrganizationID: &data.OrganizationID,
		UserID:         &data.User.ID,
		SessionID:      data.SessionID,
		EventName:      string(LoginAttemptEventName),
		EventDetails:   detailsJSON,
	}, nil
}
