package store

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	"github.com/tesseral-labs/tesseral/internal/frontend/authn"
	frontendv1 "github.com/tesseral-labs/tesseral/internal/frontend/gen/tesseral/frontend/v1"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestListAuditLogEvents(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName:    "test",
		LogInWithSaml:  refOrNil(true),
		ApiKeysEnabled: refOrNil(true),
	})

	// Custom events
	for i := range 5 {
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

		_, err = u.Store.logAuditEvent(ctx, u.Store.q, logAuditEventParams{
			EventName:    fmt.Sprintf("custom.event.%d", i),
			EventDetails: eventDetails,
		})
		require.NoError(t, err)
	}

	// Resource events
	for i := range 5 {
		_, err := u.Store.CreateAPIKey(ctx, &frontendv1.CreateAPIKeyRequest{
			ApiKey: &frontendv1.APIKey{
				DisplayName: fmt.Sprintf("API Key %d", i),
			},
		})
		require.NoError(t, err)
	}
	for i := range 5 {
		_, err := u.Store.CreateSAMLConnection(ctx, &frontendv1.CreateSAMLConnectionRequest{
			SamlConnection: &frontendv1.SAMLConnection{
				IdpRedirectUrl: fmt.Sprintf("https://idp.example.com/saml/redirect/%d", i),
				IdpEntityId:    fmt.Sprintf("https://idp.example.com/saml/idp/%d", i),
			},
		})
		require.NoError(t, err)
	}

	t.Run("AllEvents", func(t *testing.T) {
		t.Parallel()

		resp1, err := u.Store.ListAuditLogEvents(ctx, &frontendv1.ListAuditLogEventsRequest{})
		require.NoError(t, err)
		require.NotNil(t, resp1)
		require.Len(t, resp1.AuditLogEvents, 10)
		require.NotEmpty(t, resp1.NextPageToken)

		resp2, err := u.Store.ListAuditLogEvents(ctx, &frontendv1.ListAuditLogEventsRequest{PageToken: resp1.NextPageToken})
		require.NoError(t, err)
		require.NotNil(t, resp2)
		require.Len(t, resp2.AuditLogEvents, 5)
		require.Empty(t, resp2.NextPageToken)
	})

	t.Run("FilterByEventName", func(t *testing.T) {
		t.Parallel()

		for i := range 5 {
			eventName := fmt.Sprintf("custom.event.%d", i)
			respCustom, err := u.Store.ListAuditLogEvents(ctx, &frontendv1.ListAuditLogEventsRequest{
				FilterEventName: eventName,
			})
			require.NoError(t, err)
			require.NotNil(t, respCustom)
			require.Len(t, respCustom.AuditLogEvents, 1)
			require.Equal(t, eventName, respCustom.AuditLogEvents[0].EventName)
		}

		respApiKeys, err := u.Store.ListAuditLogEvents(ctx, &frontendv1.ListAuditLogEventsRequest{
			FilterEventName: "tesseral.api_keys.create",
		})
		require.NoError(t, err)
		require.NotNil(t, respApiKeys)
		require.Len(t, respApiKeys.AuditLogEvents, 5)

		respSaml, err := u.Store.ListAuditLogEvents(ctx, &frontendv1.ListAuditLogEventsRequest{
			FilterEventName: "tesseral.saml_connections.create",
		})
		require.NoError(t, err)
		require.NotNil(t, respSaml)
		require.Len(t, respSaml.AuditLogEvents, 5)
	})

	t.Run("FilterByUserId", func(t *testing.T) {
		t.Parallel()

		resp, err := u.Store.ListAuditLogEvents(ctx, &frontendv1.ListAuditLogEventsRequest{
			FilterUserId: idformat.User.Format(authn.UserID(ctx)),
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Len(t, resp.AuditLogEvents, 10)
		require.NotEmpty(t, resp.NextPageToken)

		resp2, err := u.Store.ListAuditLogEvents(ctx, &frontendv1.ListAuditLogEventsRequest{
			FilterUserId: idformat.User.Format(authn.UserID(ctx)),
			PageToken:    resp.NextPageToken,
		})
		require.NoError(t, err)
		require.NotNil(t, resp2)
		require.Len(t, resp2.AuditLogEvents, 5)
		require.Empty(t, resp2.NextPageToken)
	})

	t.Run("FilterByEventNameAndUserId", func(t *testing.T) {
		t.Parallel()

		eventName := "tesseral.api_keys.create"
		resp, err := u.Store.ListAuditLogEvents(ctx, &frontendv1.ListAuditLogEventsRequest{
			FilterEventName: eventName,
			FilterUserId:    idformat.User.Format(authn.UserID(ctx)),
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Len(t, resp.AuditLogEvents, 5)
		require.Empty(t, resp.NextPageToken)

		resp2, err := u.Store.ListAuditLogEvents(ctx, &frontendv1.ListAuditLogEventsRequest{
			FilterEventName: eventName,
			FilterUserId:    idformat.User.Format(uuid.New()),
		})
		require.NoError(t, err)
		require.NotNil(t, resp2)
		require.Len(t, resp2.AuditLogEvents, 0)
		require.Empty(t, resp2.NextPageToken)
	})

	t.Run("FilterByStartTime", func(t *testing.T) {
		t.Parallel()

		startTime := time.Now().Add(-1 * time.Hour)
		resp, err := u.Store.ListAuditLogEvents(ctx, &frontendv1.ListAuditLogEventsRequest{
			FilterStartTime: timestamppb.New(startTime),
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Len(t, resp.AuditLogEvents, 10)
		require.NotEmpty(t, resp.NextPageToken)

		resp2, err := u.Store.ListAuditLogEvents(ctx, &frontendv1.ListAuditLogEventsRequest{
			FilterStartTime: timestamppb.New(startTime),
			PageToken:       resp.NextPageToken,
		})
		require.NoError(t, err)
		require.NotNil(t, resp2)
		require.Len(t, resp2.AuditLogEvents, 5)
		require.Empty(t, resp2.NextPageToken)
	})

	t.Run("FilterByEndTime", func(t *testing.T) {
		t.Parallel()

		endTime := time.Now().Add(1 * time.Hour)
		resp, err := u.Store.ListAuditLogEvents(ctx, &frontendv1.ListAuditLogEventsRequest{
			FilterEndTime: timestamppb.New(endTime),
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Len(t, resp.AuditLogEvents, 10)
		require.NotEmpty(t, resp.NextPageToken)

		resp2, err := u.Store.ListAuditLogEvents(ctx, &frontendv1.ListAuditLogEventsRequest{
			FilterEndTime: timestamppb.New(endTime),
			PageToken:     resp.NextPageToken,
		})
		require.NoError(t, err)
		require.NotNil(t, resp2)
		require.Len(t, resp2.AuditLogEvents, 5)
		require.Empty(t, resp2.NextPageToken)
	})

	t.Run("FilterByStartAndEndTime", func(t *testing.T) {
		t.Parallel()

		startTime := time.Now().Add(-1 * time.Hour)
		endTime := time.Now().Add(1 * time.Hour)
		resp, err := u.Store.ListAuditLogEvents(ctx, &frontendv1.ListAuditLogEventsRequest{
			FilterStartTime: timestamppb.New(startTime),
			FilterEndTime:   timestamppb.New(endTime),
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Len(t, resp.AuditLogEvents, 10)
		require.NotEmpty(t, resp.NextPageToken)

		resp2, err := u.Store.ListAuditLogEvents(ctx, &frontendv1.ListAuditLogEventsRequest{
			FilterStartTime: timestamppb.New(startTime),
			FilterEndTime:   timestamppb.New(endTime),
			PageToken:       resp.NextPageToken,
		})
		require.NoError(t, err)
		require.NotNil(t, resp2)
		require.Len(t, resp2.AuditLogEvents, 5)
		require.Empty(t, resp2.NextPageToken)
	})
}
