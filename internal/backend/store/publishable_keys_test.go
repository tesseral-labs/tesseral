package store

import (
	"testing"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

func TestCreatePublishableKey_Success(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)

	resp, err := u.Store.CreatePublishableKey(ctx, &backendv1.CreatePublishableKeyRequest{
		PublishableKey: &backendv1.PublishableKey{
			DisplayName:     "key1",
			CrossDomainMode: refOrNil(true),
		},
	})
	require.NoError(t, err)
	require.NotNil(t, resp.PublishableKey)
	require.Equal(t, "key1", resp.PublishableKey.DisplayName)
	require.NotEmpty(t, resp.PublishableKey.Id)
	require.NotEmpty(t, resp.PublishableKey.CreateTime)
	require.NotEmpty(t, resp.PublishableKey.UpdateTime)
	require.True(t, *resp.PublishableKey.CrossDomainMode)
}

func TestGetPublishableKey_Exists(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)

	createResp, err := u.Store.CreatePublishableKey(ctx, &backendv1.CreatePublishableKeyRequest{
		PublishableKey: &backendv1.PublishableKey{
			DisplayName: "key1",
		},
	})
	require.NoError(t, err)
	keyID := createResp.PublishableKey.Id

	getResp, err := u.Store.GetPublishableKey(ctx, &backendv1.GetPublishableKeyRequest{Id: keyID})
	require.NoError(t, err)
	require.NotNil(t, getResp.PublishableKey)
	require.Equal(t, keyID, getResp.PublishableKey.Id)
}

func TestGetPublishableKey_DoesNotExist(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)

	_, err := u.Store.GetPublishableKey(ctx, &backendv1.GetPublishableKeyRequest{
		Id: idformat.PublishableKey.Format(uuid.New()),
	})
	var connectErr *connect.Error
	require.ErrorAs(t, err, &connectErr)
	require.Equal(t, connect.CodeNotFound, connectErr.Code())
}

func TestUpdatePublishableKey_UpdatesFields(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)

	createResp, err := u.Store.CreatePublishableKey(ctx, &backendv1.CreatePublishableKeyRequest{
		PublishableKey: &backendv1.PublishableKey{
			DisplayName: "key1",
		},
	})
	require.NoError(t, err)
	keyID := createResp.PublishableKey.Id

	updateResp, err := u.Store.UpdatePublishableKey(ctx, &backendv1.UpdatePublishableKeyRequest{
		Id: keyID,
		PublishableKey: &backendv1.PublishableKey{
			DisplayName:     "key2",
			CrossDomainMode: refOrNil(false),
		},
	})
	require.NoError(t, err)
	require.Equal(t, "key2", updateResp.PublishableKey.DisplayName)
	require.False(t, *updateResp.PublishableKey.CrossDomainMode)
}

func TestDeletePublishableKey_RemovesKey(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)

	createResp, err := u.Store.CreatePublishableKey(ctx, &backendv1.CreatePublishableKeyRequest{
		PublishableKey: &backendv1.PublishableKey{
			DisplayName: "key1",
		},
	})
	require.NoError(t, err)
	keyID := createResp.PublishableKey.Id

	_, err = u.Store.DeletePublishableKey(ctx, &backendv1.DeletePublishableKeyRequest{Id: keyID})
	require.NoError(t, err)

	_, err = u.Store.GetPublishableKey(ctx, &backendv1.GetPublishableKeyRequest{Id: keyID})
	var connectErr *connect.Error
	require.ErrorAs(t, err, &connectErr)
	require.Equal(t, connect.CodeNotFound, connectErr.Code())
}

func TestListPublishableKeys_ReturnsAll(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)

	var ids []string
	for range 3 {
		resp, err := u.Store.CreatePublishableKey(ctx, &backendv1.CreatePublishableKeyRequest{
			PublishableKey: &backendv1.PublishableKey{
				DisplayName: "test-key",
			},
		})
		require.NoError(t, err)
		ids = append(ids, resp.PublishableKey.Id)
	}

	listResp, err := u.Store.ListPublishableKeys(ctx, &backendv1.ListPublishableKeysRequest{})
	require.NoError(t, err)
	require.NotNil(t, listResp)
	require.Len(t, listResp.PublishableKeys, 3)

	var respIds []string
	for _, key := range listResp.PublishableKeys {
		respIds = append(respIds, key.Id)
	}
	require.ElementsMatch(t, ids, respIds)
}

func TestListPublishableKeys_Pagination(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)

	var createdIDs []string
	for range 15 {
		resp, err := u.Store.CreatePublishableKey(ctx, &backendv1.CreatePublishableKeyRequest{
			PublishableKey: &backendv1.PublishableKey{
				DisplayName: "test-key",
			},
		})
		require.NoError(t, err)
		createdIDs = append(createdIDs, resp.PublishableKey.Id)
	}

	resp1, err := u.Store.ListPublishableKeys(ctx, &backendv1.ListPublishableKeysRequest{})
	require.NoError(t, err)
	require.NotNil(t, resp1)
	require.Len(t, resp1.PublishableKeys, 10)
	require.NotEmpty(t, resp1.NextPageToken)

	resp2, err := u.Store.ListPublishableKeys(ctx, &backendv1.ListPublishableKeysRequest{PageToken: resp1.NextPageToken})
	require.NoError(t, err)
	require.NotNil(t, resp2)
	require.Len(t, resp2.PublishableKeys, 5)
	require.Empty(t, resp2.NextPageToken)

	var allIDs []string
	for _, k := range resp1.PublishableKeys {
		allIDs = append(allIDs, k.Id)
	}
	for _, k := range resp2.PublishableKeys {
		allIDs = append(allIDs, k.Id)
	}
	require.ElementsMatch(t, createdIDs, allIDs)
}
