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
	CreateUserEventName EventName = "user.create"
	UpdateUserEventName EventName = "user.update"
	AuthLoginEventName  EventName = "auth.login"
)

type UserData struct {
	ID    uuid.UUID
	Email string
}

type userDetails struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

func (data UserData) details() userDetails {
	return userDetails{
		ID:    idformat.User.Format(data.ID),
		Email: data.Email,
	}
}

type CreateUserEventData struct {
	ProjectID      uuid.UUID
	OrganizationID uuid.UUID

	// The user created.
	User UserData
}

type createUserEventDetails struct {
	User userDetails `json:"user"`
}

func (data CreateUserEventData) details() createUserEventDetails {
	return createUserEventDetails{
		User: data.User.details(),
	}
}

func NewCreateUserEvent(data CreateUserEventData) (Event, error) {
	details := data.details()
	detailsJSON, err := json.Marshal(details)
	if err != nil {
		return Event{}, err
	}
	return Event{
		ProjectID:      data.ProjectID,
		OrganizationID: &data.OrganizationID,
		EventName:      string(CreateUserEventName),
		EventDetails:   detailsJSON,
	}, nil
}

type UserUpdate struct {
	Email             string `json:"email,omitempty"`
	GoogleUserID      string `json:"google_user_id,omitempty"`
	MicrosoftUserID   string `json:"microsoft_user_id,omitempty"`
	GithubUserID      string `json:"github_user_id,omitempty"`
	IsOwner           *bool  `json:"is_owner,omitempty"`
	DisplayName       string `json:"display_name,omitempty"`
	ProfilePictureURL string `json:"profile_picture_url,omitempty"`
}

type UpdateUserEventData struct {
	ProjectID      uuid.UUID
	OrganizationID uuid.UUID

	// The user updated.
	User UserData

	// The information updated for the user.
	Update UserUpdate
}

type updateUserEventDetails struct {
	User   userDetails `json:"user"`
	Update UserUpdate  `json:"update"`
}

func (data UpdateUserEventData) details() updateUserEventDetails {
	return updateUserEventDetails{
		User:   data.User.details(),
		Update: data.Update,
	}
}

func NewUpdateUserEvent(data UpdateUserEventData) (Event, error) {
	details := data.details()
	detailsJSON, err := json.Marshal(details)
	if err != nil {
		return Event{}, err
	}
	return Event{
		ProjectID:      data.ProjectID,
		OrganizationID: &data.OrganizationID,
		EventName:      string(UpdateUserEventName),
		EventDetails:   detailsJSON,
	}, nil
}

type AuthLoginEventData struct {
	ProjectID             uuid.UUID
	OrganizationID        uuid.UUID
	IntermediateSessionID uuid.UUID
	SessionID             uuid.UUID
	User                  UserData
	Factor                string
	Success               bool
}

type authLoginEventDetails struct {
	IntermediateSessionID string      `json:"intermediate_session_id"`
	User                  userDetails `json:"user"`
	Factor                string      `json:"factor"`
	Success               bool        `json:"success"`
}

func (data AuthLoginEventData) details() authLoginEventDetails {
	return authLoginEventDetails{
		IntermediateSessionID: idformat.IntermediateSession.Format(data.IntermediateSessionID),
		User:                  data.User.details(),
		Factor:                data.Factor,
		Success:               data.Success,
	}
}

func NewAuthLoginEvent(data AuthLoginEventData) (Event, error) {
	details := data.details()
	detailsJSON, err := json.Marshal(details)
	if err != nil {
		return Event{}, err
	}
	return Event{
		ProjectID:      data.ProjectID,
		OrganizationID: &data.OrganizationID,
		UserID:         &data.User.ID,
		SessionID:      &data.SessionID,
		EventName:      string(AuthLoginEventName),
		EventDetails:   detailsJSON,
	}, nil
}
