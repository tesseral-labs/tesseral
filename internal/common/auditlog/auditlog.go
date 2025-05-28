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
	ID                uuid.UUID
	Email             string
	GoogleUserID      *string
	MicrosoftUserID   *string
	GithubUserID      *string
	IsOwner           bool
	DisplayName       *string
	ProfilePictureURL *string
}

type userDetails struct {
	ID                string  `json:"id"`
	Email             string  `json:"email"`
	GoogleUserID      *string `json:"google_user_id"`
	MicrosoftUserID   *string `json:"microsoft_user_id"`
	GithubUserID      *string `json:"github_user_id"`
	IsOwner           bool    `json:"is_owner"`
	DisplayName       *string `json:"display_name"`
	ProfilePictureURL *string `json:"profile_picture_url"`
}

func (data UserData) details() userDetails {
	return userDetails{
		ID:                idformat.User.Format(data.ID),
		Email:             data.Email,
		GoogleUserID:      data.GoogleUserID,
		MicrosoftUserID:   data.MicrosoftUserID,
		GithubUserID:      data.GithubUserID,
		IsOwner:           data.IsOwner,
		DisplayName:       data.DisplayName,
		ProfilePictureURL: data.ProfilePictureURL,
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

type UpdateUserEventData struct {
	ProjectID      uuid.UUID
	OrganizationID uuid.UUID

	// The user after the update.
	User UserData

	// The user before the update.
	PreviousUser UserData
}

type updateUserEventDetails struct {
	User         userDetails `json:"user"`
	PreviousUser userDetails `json:"previous_user"`
}

func (data UpdateUserEventData) details() updateUserEventDetails {
	return updateUserEventDetails{
		User:         data.User.details(),
		PreviousUser: data.PreviousUser.details(),
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
	PrimaryAuthFactor     string
	Success               bool
}

type authLoginEventDetails struct {
	IntermediateSessionID string      `json:"intermediate_session_id"`
	User                  userDetails `json:"user"`
	PrimaryAuthFactor     string      `json:"primary_auth_factor"`
	Success               bool        `json:"success"`
}

func (data AuthLoginEventData) details() authLoginEventDetails {
	return authLoginEventDetails{
		IntermediateSessionID: idformat.IntermediateSession.Format(data.IntermediateSessionID),
		User:                  data.User.details(),
		PrimaryAuthFactor:     data.PrimaryAuthFactor,
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
