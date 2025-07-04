package store

import (
	"fmt"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/tesseral-labs/tesseral/internal/backend/authn"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	"github.com/tesseral-labs/tesseral/internal/backend/store/queries"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
	"github.com/tesseral-labs/tesseral/internal/uuidv7"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestCreateCustomAuditLogEvent_Success(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)
	orgID := u.Environment.NewOrganization(t, u.ProjectID, &backendv1.Organization{
		DisplayName: "test",
	})

	eventTime := time.Now()
	eventDetailsMap := map[string]any{
		"key1": "value1",
		"key2": 42,
		"key3": true,
		"nested": map[string]any{
			"nestedKey1": "nestedValue1",
			"nestedKey2": 3.14,
		},
		"list":  []any{"item1", "item2", 100},
		"empty": nil,
	}
	eventDetails, err := structpb.NewStruct(eventDetailsMap)
	require.NoError(t, err)

	resp, err := u.Store.CreateCustomAuditLogEvent(ctx, &backendv1.CreateAuditLogEventRequest{
		AuditLogEvent: &backendv1.AuditLogEvent{
			OrganizationId: orgID,
			EventName:      "custom.event.created",
			EventTime:      timestamppb.New(eventTime),
			EventDetails:   eventDetails,
		},
	})
	require.NoError(t, err)
	require.NotNil(t, resp.AuditLogEvent)
	require.Equal(t, orgID, resp.AuditLogEvent.OrganizationId)
	require.Equal(t, "custom.event.created", resp.AuditLogEvent.EventName)
	require.Equal(t, eventTime.UnixMilli(), resp.AuditLogEvent.EventTime.AsTime().UnixMilli())

	respEventDetails := resp.AuditLogEvent.GetEventDetails().AsMap()
	require.EqualValues(t, eventDetailsMap["key1"], respEventDetails["key1"])
	require.EqualValues(t, eventDetailsMap["key2"], respEventDetails["key2"])
	require.EqualValues(t, eventDetailsMap["key3"], respEventDetails["key3"])
	require.Equal(t, eventDetailsMap["nested"], respEventDetails["nested"])
	require.Len(t, respEventDetails["list"], 3)
	require.EqualValues(t, "item1", respEventDetails["list"].([]any)[0])
	require.EqualValues(t, "item2", respEventDetails["list"].([]any)[1])
	require.EqualValues(t, 100, respEventDetails["list"].([]any)[2])
	require.Contains(t, respEventDetails, "empty")
}

func TestCreateCustomAuditLogEvent_Actor(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)
	orgID := u.Environment.NewOrganization(t, u.ProjectID, &backendv1.Organization{
		DisplayName:    "test",
		ApiKeysEnabled: refOrNil(true),
	})

	userID := u.Environment.NewUser(t, orgID, &backendv1.User{
		Email: "test@example.com",
	})
	sessionID, refreshToken := u.Environment.NewSession(t, userID)

	accessToken, err := u.Common.IssueAccessToken(ctx, authn.ProjectID(ctx), refreshToken)
	require.NoError(t, err)

	apiKey, err := u.Store.CreateAPIKey(ctx, &backendv1.CreateAPIKeyRequest{
		ApiKey: &backendv1.APIKey{
			OrganizationId: orgID,
			DisplayName:    "test-key",
		},
	})
	require.NoError(t, err)
	require.NotNil(t, apiKey)
	apiKeyID := apiKey.ApiKey.Id

	backendAPIKeyID := idformat.BackendAPIKey.Format(uuid.New())
	intermediateSessionID := idformat.IntermediateSession.Format(uuid.New())

	tt := []struct {
		Name                       string
		OrganizationID             string
		ActorUserID                string
		ActorAPIKeyID              string
		ActorSessionID             string
		ActorBackendAPIKeyID       string
		ActorIntermediateSessionID string
		ActorCredentials           string
		wantError                  connect.Code
	}{
		{
			Name:      "NoActor",
			wantError: connect.CodeInvalidArgument,
		},
		{
			Name:           "WithOrganizationID",
			OrganizationID: orgID,
		},
		{
			Name:        "WithActorUserID",
			ActorUserID: userID,
		},
		{
			Name:           "WithSessionID",
			ActorSessionID: sessionID,
		},
		{
			Name:          "WithAPIKeyID",
			ActorAPIKeyID: apiKeyID,
		},
		{
			Name:                 "WithActorBackendAPIKeyID",
			ActorBackendAPIKeyID: backendAPIKeyID,
			wantError:            connect.CodeInvalidArgument,
		},
		{
			Name:                       "WithActorIntermediateSessionID",
			ActorIntermediateSessionID: intermediateSessionID,
			wantError:                  connect.CodeInvalidArgument,
		},
		{
			Name:           "WithMultipleActors",
			OrganizationID: orgID,
			ActorUserID:    userID,
			wantError:      connect.CodeInvalidArgument,
		},
		{
			Name:             "WithUserCredentials",
			ActorCredentials: accessToken,
		},
		{
			Name:             "WithAPIKeyCredentials",
			ActorCredentials: apiKey.ApiKey.SecretToken,
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			resp, err := u.Store.CreateCustomAuditLogEvent(ctx, &backendv1.CreateAuditLogEventRequest{
				AuditLogEvent: &backendv1.AuditLogEvent{
					OrganizationId:             tc.OrganizationID,
					ActorUserId:                tc.ActorUserID,
					ActorApiKeyId:              tc.ActorAPIKeyID,
					ActorSessionId:             tc.ActorSessionID,
					ActorBackendApiKeyId:       tc.ActorBackendAPIKeyID,
					ActorIntermediateSessionId: tc.ActorIntermediateSessionID,
					ActorCredentials:           tc.ActorCredentials,
					EventName:                  "custom.event.created",
				},
			})
			if tc.wantError != 0 {
				var connectErr *connect.Error
				require.ErrorAs(t, err, &connectErr)
				require.Equal(t, tc.wantError, connectErr.Code())
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp.AuditLogEvent)
				require.Equal(t, orgID, resp.AuditLogEvent.OrganizationId)
				require.Equal(t, "custom.event.created", resp.AuditLogEvent.EventName)
			}
		})
	}
}

func TestCreateAuditLogEvent_InvalidEventName(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)
	orgID := u.Environment.NewOrganization(t, u.ProjectID, &backendv1.Organization{
		DisplayName: "test",
	})

	_, err := u.Store.CreateCustomAuditLogEvent(ctx, &backendv1.CreateAuditLogEventRequest{
		AuditLogEvent: &backendv1.AuditLogEvent{
			OrganizationId: orgID,
			EventName:      "invalidName",
		},
	})
	var connectErr *connect.Error
	require.ErrorAs(t, err, &connectErr)
	require.Equal(t, connect.CodeInvalidArgument, connectErr.Code())
}

func TestConsoleListCustomAuditLogEvents_ReturnsAll(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)
	orgID := u.Environment.NewOrganization(t, u.ProjectID, &backendv1.Organization{
		DisplayName: "test",
	})

	for range 3 {
		_, err := u.Store.CreateCustomAuditLogEvent(ctx, &backendv1.CreateAuditLogEventRequest{
			AuditLogEvent: &backendv1.AuditLogEvent{
				OrganizationId: orgID,
				EventName:      "custom.event.created",
			},
		})
		require.NoError(t, err)
	}

	listResp, err := u.Store.ConsoleListCustomAuditLogEvents(ctx, &backendv1.ConsoleListAuditLogEventsRequest{
		OrganizationId: orgID,
	})
	require.NoError(t, err)
	require.NotNil(t, listResp)
	require.Len(t, listResp.AuditLogEvents, 3)
}

func TestConsoleListCustomAuditLogEvents_PaginateByActor(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)
	orgID := u.Environment.NewOrganization(t, u.ProjectID, &backendv1.Organization{
		DisplayName:    "test",
		ApiKeysEnabled: refOrNil(true),
		ScimEnabled:    refOrNil(true),
	})
	orgUUID, err := idformat.Organization.Parse(orgID)
	require.NoError(t, err)

	userID := u.Environment.NewUser(t, orgID, &backendv1.User{
		Email: "test@example.com",
	})
	sessionID, refreshToken := u.Environment.NewSession(t, userID)

	accessToken, err := u.Common.IssueAccessToken(ctx, authn.ProjectID(ctx), refreshToken)
	require.NoError(t, err)

	apiKey, err := u.Store.CreateAPIKey(ctx, &backendv1.CreateAPIKeyRequest{
		ApiKey: &backendv1.APIKey{
			OrganizationId: orgID,
			DisplayName:    "test-key",
		},
	})
	require.NoError(t, err)
	require.NotNil(t, apiKey)
	apiKeyID := apiKey.ApiKey.Id

	scimApiKey, err := u.Store.CreateSCIMAPIKey(ctx, &backendv1.CreateSCIMAPIKeyRequest{
		ScimApiKey: &backendv1.SCIMAPIKey{
			OrganizationId: orgID,
			DisplayName:    "test-scim-key",
		},
	})
	require.NoError(t, err)
	require.NotNil(t, scimApiKey)
	scimApiKeyID, err := idformat.SCIMAPIKey.Parse(scimApiKey.ScimApiKey.Id)
	require.NoError(t, err)

	// Org-wide events
	for i := range 5 {
		_, err := u.Store.CreateCustomAuditLogEvent(ctx, &backendv1.CreateAuditLogEventRequest{
			AuditLogEvent: &backendv1.AuditLogEvent{
				OrganizationId: orgID,
				EventName:      fmt.Sprintf("custom.organization.%d", i),
			},
		})
		require.NoError(t, err)
	}
	// User-specific events
	for i := range 5 {
		_, err := u.Store.CreateCustomAuditLogEvent(ctx, &backendv1.CreateAuditLogEventRequest{
			AuditLogEvent: &backendv1.AuditLogEvent{
				ActorCredentials: accessToken,
				EventName:        fmt.Sprintf("custom.user.%d", i),
			},
		})
		require.NoError(t, err)
	}
	// API Key-specific events
	for i := range 5 {
		_, err := u.Store.CreateCustomAuditLogEvent(ctx, &backendv1.CreateAuditLogEventRequest{
			AuditLogEvent: &backendv1.AuditLogEvent{
				ActorCredentials: apiKey.ApiKey.SecretToken,
				EventName:        fmt.Sprintf("custom.api_key.%d", i),
			},
		})
		require.NoError(t, err)
	}
	// SCIM API Key-specific events
	for i := range 5 {
		eventTime := time.Now()
		eventID := uuidv7.NewWithTime(eventTime)
		_, err := u.Store.q.CreateAuditLogEvent(ctx, queries.CreateAuditLogEventParams{
			ID:                eventID,
			ProjectID:         authn.ProjectID(ctx),
			OrganizationID:    (*uuid.UUID)(&orgUUID),
			EventName:         fmt.Sprintf("custom.scim_api_key.%d", i),
			EventTime:         &eventTime,
			ActorScimApiKeyID: (*uuid.UUID)(&scimApiKeyID),
		})
		require.NoError(t, err)
	}

	t.Run("AllEvents", func(t *testing.T) {
		t.Parallel()

		resp1, err := u.Store.ConsoleListCustomAuditLogEvents(ctx, &backendv1.ConsoleListAuditLogEventsRequest{
			OrganizationId: orgID,
		})
		require.NoError(t, err)
		require.NotNil(t, resp1)
		require.Len(t, resp1.AuditLogEvents, 10)
		require.NotEmpty(t, resp1.NextPageToken)

		resp2, err := u.Store.ConsoleListCustomAuditLogEvents(ctx, &backendv1.ConsoleListAuditLogEventsRequest{
			OrganizationId: orgID,
			PageToken:      resp1.NextPageToken,
		})
		require.NoError(t, err)
		require.NotNil(t, resp2)
		require.Len(t, resp2.AuditLogEvents, 10)
		require.NotEmpty(t, resp2.NextPageToken)

		resp3, err := u.Store.ConsoleListCustomAuditLogEvents(ctx, &backendv1.ConsoleListAuditLogEventsRequest{
			OrganizationId: orgID,
			PageToken:      resp2.NextPageToken,
		})
		require.NoError(t, err)
		require.NotNil(t, resp3)
		require.Len(t, resp3.AuditLogEvents, 2 /* CreateAPIKey + CreateSCIMAPIKey create an event */)
		require.Empty(t, resp3.NextPageToken)
	})

	t.Run("ByActorUserID", func(t *testing.T) {
		t.Parallel()

		resp, err := u.Store.ConsoleListCustomAuditLogEvents(ctx, &backendv1.ConsoleListAuditLogEventsRequest{
			OrganizationId: orgID,
			ActorUserId:    userID,
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Len(t, resp.AuditLogEvents, 5)
		require.Empty(t, resp.NextPageToken)

		respEventNames, err := u.Store.ConsoleListAuditLogEventNames(ctx, &backendv1.ConsoleListAuditLogEventNamesRequest{
			ActorUserId: userID,
		})
		require.NoError(t, err)
		require.NotNil(t, respEventNames)
		require.ElementsMatch(t, []string{
			"custom.user.0",
			"custom.user.1",
			"custom.user.2",
			"custom.user.3",
			"custom.user.4",
		}, respEventNames.EventNames)
	})

	t.Run("ByActorSessionID", func(t *testing.T) {
		t.Parallel()

		resp, err := u.Store.ConsoleListCustomAuditLogEvents(ctx, &backendv1.ConsoleListAuditLogEventsRequest{
			OrganizationId: orgID,
			ActorSessionId: sessionID,
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Len(t, resp.AuditLogEvents, 5)
		require.Empty(t, resp.NextPageToken)

		respEventNames, err := u.Store.ConsoleListAuditLogEventNames(ctx, &backendv1.ConsoleListAuditLogEventNamesRequest{
			ActorSessionId: sessionID,
		})
		require.NoError(t, err)
		require.NotNil(t, respEventNames)
		require.ElementsMatch(t, []string{
			"custom.user.0",
			"custom.user.1",
			"custom.user.2",
			"custom.user.3",
			"custom.user.4",
		}, respEventNames.EventNames)
	})

	t.Run("ByActorAPIKeyID", func(t *testing.T) {
		t.Parallel()

		resp, err := u.Store.ConsoleListCustomAuditLogEvents(ctx, &backendv1.ConsoleListAuditLogEventsRequest{
			OrganizationId: orgID,
			ActorApiKeyId:  apiKeyID,
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Len(t, resp.AuditLogEvents, 5)
		require.Empty(t, resp.NextPageToken)

		respEventNames, err := u.Store.ConsoleListAuditLogEventNames(ctx, &backendv1.ConsoleListAuditLogEventNamesRequest{
			ActorApiKeyId: apiKeyID,
		})
		require.NoError(t, err)
		require.NotNil(t, respEventNames)
		require.ElementsMatch(t, []string{
			"custom.api_key.0",
			"custom.api_key.1",
			"custom.api_key.2",
			"custom.api_key.3",
			"custom.api_key.4",
		}, respEventNames.EventNames)
	})

	t.Run("ByActorSCIMAPIKeyID", func(t *testing.T) {
		t.Parallel()

		resp, err := u.Store.ConsoleListCustomAuditLogEvents(ctx, &backendv1.ConsoleListAuditLogEventsRequest{
			OrganizationId:    orgID,
			ActorScimApiKeyId: idformat.SCIMAPIKey.Format(scimApiKeyID),
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Len(t, resp.AuditLogEvents, 5)
		require.Empty(t, resp.NextPageToken)

		respEventNames, err := u.Store.ConsoleListAuditLogEventNames(ctx, &backendv1.ConsoleListAuditLogEventNamesRequest{
			ActorScimApiKeyId: idformat.SCIMAPIKey.Format(scimApiKeyID),
		})
		require.NoError(t, err)
		require.NotNil(t, respEventNames)
		require.ElementsMatch(t, []string{
			"custom.scim_api_key.0",
			"custom.scim_api_key.1",
			"custom.scim_api_key.2",
			"custom.scim_api_key.3",
			"custom.scim_api_key.4",
		}, respEventNames.EventNames)
	})
}

func TestConsoleListCustomAuditLogEvents_ResourceTypeUnspecified(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)

	// Create some events for the project
	var organizationIDs []string
	for range 3 {
		resp, err := u.Store.CreateOrganization(ctx, &backendv1.CreateOrganizationRequest{
			Organization: &backendv1.Organization{
				DisplayName: "Test Org",
			},
		})
		require.NoError(t, err)
		require.NotNil(t, resp.Organization)
		organizationIDs = append(organizationIDs, resp.Organization.Id)
	}

	resp, err := u.Store.ConsoleListCustomAuditLogEvents(ctx, &backendv1.ConsoleListAuditLogEventsRequest{
		ResourceType: backendv1.AuditLogEventResourceType_AUDIT_LOG_EVENT_RESOURCE_TYPE_UNSPECIFIED,
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Len(t, resp.AuditLogEvents, 3)
	require.Empty(t, resp.NextPageToken)

	for _, orgID := range organizationIDs {
		resp, err := u.Store.ConsoleListCustomAuditLogEvents(ctx, &backendv1.ConsoleListAuditLogEventsRequest{
			ResourceType: backendv1.AuditLogEventResourceType_AUDIT_LOG_EVENT_RESOURCE_TYPE_ORGANIZATION,
			ResourceId:   orgID,
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Len(t, resp.AuditLogEvents, 1)
		require.Empty(t, resp.NextPageToken)
	}
	for _, orgID := range organizationIDs {
		resp, err := u.Store.ConsoleListCustomAuditLogEvents(ctx, &backendv1.ConsoleListAuditLogEventsRequest{
			OrganizationId: orgID,
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Len(t, resp.AuditLogEvents, 1)
		require.Empty(t, resp.NextPageToken)
	}
}

func TestConsoleListCustomAuditLogEvents_ResourceTypeAPIKey(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)
	orgID := u.Environment.NewOrganization(t, u.ProjectID, &backendv1.Organization{
		DisplayName:    "test",
		ApiKeysEnabled: refOrNil(true),
	})

	var apiKeyIDs []string
	for range 3 {
		resp, err := u.Store.CreateAPIKey(ctx, &backendv1.CreateAPIKeyRequest{
			ApiKey: &backendv1.APIKey{
				OrganizationId: orgID,
				DisplayName:    "Test API Key",
			},
		})
		require.NoError(t, err)
		apiKeyIDs = append(apiKeyIDs, resp.ApiKey.Id)
	}

	resp, err := u.Store.ConsoleListCustomAuditLogEvents(ctx, &backendv1.ConsoleListAuditLogEventsRequest{
		OrganizationId: orgID,
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Len(t, resp.AuditLogEvents, 3)
	require.Empty(t, resp.NextPageToken)

	for _, apiKeyID := range apiKeyIDs {
		resp, err := u.Store.ConsoleListCustomAuditLogEvents(ctx, &backendv1.ConsoleListAuditLogEventsRequest{
			ResourceType: backendv1.AuditLogEventResourceType_AUDIT_LOG_EVENT_RESOURCE_TYPE_API_KEY,
			ResourceId:   apiKeyID,
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Len(t, resp.AuditLogEvents, 1)
		require.Empty(t, resp.NextPageToken)
	}
}

func TestConsoleListCustomAuditLogEvents_ResourceTypeOrganization(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)
	orgID := u.Environment.NewOrganization(t, u.ProjectID, &backendv1.Organization{
		DisplayName:    "test",
		ApiKeysEnabled: refOrNil(true),
	})

	for range 3 {
		_, err := u.Store.UpdateOrganization(ctx, &backendv1.UpdateOrganizationRequest{
			Id: orgID,
			Organization: &backendv1.Organization{
				DisplayName: "Test Org",
			},
		})
		require.NoError(t, err)
	}

	// List by organization
	resp, err := u.Store.ConsoleListCustomAuditLogEvents(ctx, &backendv1.ConsoleListAuditLogEventsRequest{
		OrganizationId: orgID,
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Len(t, resp.AuditLogEvents, 3)
	require.Empty(t, resp.NextPageToken)

	// List by resource type = organization
	resp2, err := u.Store.ConsoleListCustomAuditLogEvents(ctx, &backendv1.ConsoleListAuditLogEventsRequest{
		ResourceType: backendv1.AuditLogEventResourceType_AUDIT_LOG_EVENT_RESOURCE_TYPE_ORGANIZATION,
		ResourceId:   orgID,
	})
	require.NoError(t, err)
	require.NotNil(t, resp2)
	require.Len(t, resp2.AuditLogEvents, 3)
	require.Empty(t, resp2.NextPageToken)
}

func TestConsoleListCustomAuditLogEvents_ResourceTypePasskey(t *testing.T) {
	t.Skip("Passkey audit log events are only recorded in frontend")
}

func TestConsoleListCustomAuditLogEvents_ResourceTypeRole(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)
	orgID := u.Environment.NewOrganization(t, u.ProjectID, &backendv1.Organization{
		DisplayName: "test",
	})

	var roleIDs []string
	for range 3 {
		resp, err := u.Store.CreateRole(ctx, &backendv1.CreateRoleRequest{
			Role: &backendv1.Role{
				OrganizationId: orgID,
				DisplayName:    "Test Role",
			},
		})
		require.NoError(t, err)
		roleIDs = append(roleIDs, resp.Role.Id)
	}

	resp, err := u.Store.ConsoleListCustomAuditLogEvents(ctx, &backendv1.ConsoleListAuditLogEventsRequest{
		OrganizationId: orgID,
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Len(t, resp.AuditLogEvents, 3)
	require.Empty(t, resp.NextPageToken)

	for _, roleID := range roleIDs {
		resp, err := u.Store.ConsoleListCustomAuditLogEvents(ctx, &backendv1.ConsoleListAuditLogEventsRequest{
			ResourceType: backendv1.AuditLogEventResourceType_AUDIT_LOG_EVENT_RESOURCE_TYPE_ROLE,
			ResourceId:   roleID,
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Len(t, resp.AuditLogEvents, 1)
		require.Empty(t, resp.NextPageToken)
	}
}

func TestConsoleListCustomAuditLogEvents_ResourceTypeSAMLConnection(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)
	orgID := u.Environment.NewOrganization(t, u.ProjectID, &backendv1.Organization{
		DisplayName:   "test",
		LogInWithSaml: refOrNil(true),
	})

	var samlConnectionIDs []string
	for range 3 {
		resp, err := u.Store.CreateSAMLConnection(ctx, &backendv1.CreateSAMLConnectionRequest{
			SamlConnection: &backendv1.SAMLConnection{
				OrganizationId: orgID,
				IdpRedirectUrl: "https://idp.example.com/saml/redirect",
				IdpEntityId:    "https://idp.example.com/saml/idp",
			},
		})
		require.NoError(t, err)
		samlConnectionIDs = append(samlConnectionIDs, resp.SamlConnection.Id)
	}

	resp, err := u.Store.ConsoleListCustomAuditLogEvents(ctx, &backendv1.ConsoleListAuditLogEventsRequest{
		OrganizationId: orgID,
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Len(t, resp.AuditLogEvents, 3)
	require.Empty(t, resp.NextPageToken)

	for _, samlConnectionID := range samlConnectionIDs {
		resp, err := u.Store.ConsoleListCustomAuditLogEvents(ctx, &backendv1.ConsoleListAuditLogEventsRequest{
			ResourceType: backendv1.AuditLogEventResourceType_AUDIT_LOG_EVENT_RESOURCE_TYPE_SAML_CONNECTION,
			ResourceId:   samlConnectionID,
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Len(t, resp.AuditLogEvents, 1)
		require.Empty(t, resp.NextPageToken)
	}
}

func TestConsoleListCustomAuditLogEvents_ResourceTypeOIDCConnection(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)
	orgID := u.Environment.NewOrganization(t, u.ProjectID, &backendv1.Organization{
		DisplayName:   "test",
		LogInWithOidc: refOrNil(true),
	})

	var oidcConnectionIDs []string
	for range 3 {
		resp, err := u.Store.CreateOIDCConnection(ctx, &backendv1.CreateOIDCConnectionRequest{
			OidcConnection: &backendv1.OIDCConnection{
				OrganizationId:   orgID,
				ConfigurationUrl: "https://accounts.google.com/.well-known/openid-configuration",
				ClientId:         "client-id",
				ClientSecret:     "client-secret",
			},
		})
		require.NoError(t, err)
		oidcConnectionIDs = append(oidcConnectionIDs, resp.OidcConnection.Id)
	}

	resp, err := u.Store.ConsoleListCustomAuditLogEvents(ctx, &backendv1.ConsoleListAuditLogEventsRequest{
		OrganizationId: orgID,
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Len(t, resp.AuditLogEvents, 3)
	require.Empty(t, resp.NextPageToken)

	for _, oidcConnectionID := range oidcConnectionIDs {
		resp, err := u.Store.ConsoleListCustomAuditLogEvents(ctx, &backendv1.ConsoleListAuditLogEventsRequest{
			ResourceType: backendv1.AuditLogEventResourceType_AUDIT_LOG_EVENT_RESOURCE_TYPE_OIDC_CONNECTION,
			ResourceId:   oidcConnectionID,
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Len(t, resp.AuditLogEvents, 1)
		require.Empty(t, resp.NextPageToken)
	}
}

func TestConsoleListCustomAuditLogEvents_ResourceTypeSCIMAPIKey(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)
	orgID := u.Environment.NewOrganization(t, u.ProjectID, &backendv1.Organization{
		DisplayName: "test",
		ScimEnabled: refOrNil(true),
	})

	var scimAPIKeyIDs []string
	for range 3 {
		resp, err := u.Store.CreateSCIMAPIKey(ctx, &backendv1.CreateSCIMAPIKeyRequest{
			ScimApiKey: &backendv1.SCIMAPIKey{
				OrganizationId: orgID,
				DisplayName:    "Test SCIM Key",
			},
		})
		require.NoError(t, err)
		scimAPIKeyIDs = append(scimAPIKeyIDs, resp.ScimApiKey.Id)
	}

	resp, err := u.Store.ConsoleListCustomAuditLogEvents(ctx, &backendv1.ConsoleListAuditLogEventsRequest{
		OrganizationId: orgID,
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Len(t, resp.AuditLogEvents, 3)
	require.Empty(t, resp.NextPageToken)

	for _, scimAPIKeyID := range scimAPIKeyIDs {
		resp, err := u.Store.ConsoleListCustomAuditLogEvents(ctx, &backendv1.ConsoleListAuditLogEventsRequest{
			ResourceType: backendv1.AuditLogEventResourceType_AUDIT_LOG_EVENT_RESOURCE_TYPE_SCIM_API_KEY,
			ResourceId:   scimAPIKeyID,
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Len(t, resp.AuditLogEvents, 1)
		require.Empty(t, resp.NextPageToken)
	}
}

func TestConsoleListCustomAuditLogEvents_ResourceTypeSession(t *testing.T) {
	t.Skip("Session audit log events are only recorded in intermediate")
}

func TestConsoleListCustomAuditLogEvents_ResourceTypeUserInvite(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)
	orgID := u.Environment.NewOrganization(t, u.ProjectID, &backendv1.Organization{
		DisplayName: "test",
	})

	var inviteIDs []string
	for i := range 3 {
		resp, err := u.Store.CreateUserInvite(ctx, &backendv1.CreateUserInviteRequest{
			UserInvite: &backendv1.UserInvite{
				OrganizationId: orgID,
				Email:          fmt.Sprintf("invite-%d@example.com", i),
			},
		})
		require.NoError(t, err)
		inviteIDs = append(inviteIDs, resp.UserInvite.Id)
	}

	resp, err := u.Store.ConsoleListCustomAuditLogEvents(ctx, &backendv1.ConsoleListAuditLogEventsRequest{
		OrganizationId: orgID,
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Len(t, resp.AuditLogEvents, 3)
	require.Empty(t, resp.NextPageToken)

	for _, inviteID := range inviteIDs {
		resp, err := u.Store.ConsoleListCustomAuditLogEvents(ctx, &backendv1.ConsoleListAuditLogEventsRequest{
			ResourceType: backendv1.AuditLogEventResourceType_AUDIT_LOG_EVENT_RESOURCE_TYPE_USER_INVITE,
			ResourceId:   inviteID,
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Len(t, resp.AuditLogEvents, 1)
		require.Empty(t, resp.NextPageToken)
	}
}

func TestConsoleListCustomAuditLogEvents_ResourceTypeUser(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)
	orgID := u.Environment.NewOrganization(t, u.ProjectID, &backendv1.Organization{
		DisplayName: "test",
	})

	var userIDs []string
	for i := range 3 {
		resp, err := u.Store.CreateUser(ctx, &backendv1.CreateUserRequest{
			User: &backendv1.User{
				OrganizationId: orgID,
				Email:          fmt.Sprintf("test-%d@example.com", i),
			},
		})
		require.NoError(t, err)
		require.NotNil(t, resp.User)
		userIDs = append(userIDs, resp.User.Id)
	}

	resp, err := u.Store.ConsoleListCustomAuditLogEvents(ctx, &backendv1.ConsoleListAuditLogEventsRequest{
		OrganizationId: orgID,
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Len(t, resp.AuditLogEvents, 3)
	require.Empty(t, resp.NextPageToken)

	for _, userID := range userIDs {
		resp, err := u.Store.ConsoleListCustomAuditLogEvents(ctx, &backendv1.ConsoleListAuditLogEventsRequest{
			ResourceType: backendv1.AuditLogEventResourceType_AUDIT_LOG_EVENT_RESOURCE_TYPE_USER,
			ResourceId:   userID,
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Len(t, resp.AuditLogEvents, 1)
		require.Empty(t, resp.NextPageToken)
	}
}

func TestConsoleListAuditLogEventNames_ResourceTypeUnspecified(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)
	orgID := u.Environment.NewOrganization(t, u.ProjectID, &backendv1.Organization{
		DisplayName: "test",
	})

	_, err := u.Store.ConsoleListAuditLogEventNames(ctx, &backendv1.ConsoleListAuditLogEventNamesRequest{
		ResourceType: backendv1.AuditLogEventResourceType_AUDIT_LOG_EVENT_RESOURCE_TYPE_UNSPECIFIED,
	})
	require.Error(t, err)

	_, err = u.Store.ConsoleListAuditLogEventNames(ctx, &backendv1.ConsoleListAuditLogEventNamesRequest{
		OrganizationId: orgID,
		ResourceType:   backendv1.AuditLogEventResourceType_AUDIT_LOG_EVENT_RESOURCE_TYPE_UNSPECIFIED,
	})
	require.Error(t, err)
}
