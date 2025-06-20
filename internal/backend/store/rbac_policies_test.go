package store

import (
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/require"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
)

func TestRBACPolicy_GetAndUpdate(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)

	getResp, err := u.Store.GetRBACPolicy(ctx, &backendv1.GetRBACPolicyRequest{})
	require.NoError(t, err)
	require.NotNil(t, getResp)
	require.Empty(t, getResp.RbacPolicy.Actions)

	updateResp, err := u.Store.UpdateRBACPolicy(ctx, &backendv1.UpdateRBACPolicyRequest{
		RbacPolicy: &backendv1.RBACPolicy{
			Actions: []*backendv1.Action{
				{Name: "foo.bar.baz", Description: "desc1"},
				{Name: "foo.bar.qux", Description: "desc2"},
			},
		},
	})
	require.NoError(t, err)
	require.Len(t, updateResp.RbacPolicy.Actions, 2)

	getResp, err = u.Store.GetRBACPolicy(ctx, &backendv1.GetRBACPolicyRequest{})
	require.NoError(t, err)
	require.ElementsMatch(t, []string{"foo.bar.baz", "foo.bar.qux"},
		[]string{getResp.RbacPolicy.Actions[0].Name, getResp.RbacPolicy.Actions[1].Name})

	updateResp, err = u.Store.UpdateRBACPolicy(ctx, &backendv1.UpdateRBACPolicyRequest{
		RbacPolicy: &backendv1.RBACPolicy{
			Actions: []*backendv1.Action{
				{Name: "foo.bar.baz", Description: "desc1-updated"},
				{Name: "foo.bar.new", Description: "desc3"},
			},
		},
	})
	require.NoError(t, err)
	require.Len(t, updateResp.RbacPolicy.Actions, 2)
	var names []string
	for _, a := range updateResp.RbacPolicy.Actions {
		names = append(names, a.Name)
	}
	require.ElementsMatch(t, []string{"foo.bar.baz", "foo.bar.new"}, names)
	for _, a := range updateResp.RbacPolicy.Actions {
		if a.Name == "foo.bar.baz" {
			require.Equal(t, "desc1-updated", a.Description)
		}
	}

	updateResp, err = u.Store.UpdateRBACPolicy(ctx, &backendv1.UpdateRBACPolicyRequest{
		RbacPolicy: &backendv1.RBACPolicy{Actions: nil},
	})
	require.NoError(t, err)
	require.Empty(t, updateResp.RbacPolicy.Actions)

	getResp, err = u.Store.GetRBACPolicy(ctx, &backendv1.GetRBACPolicyRequest{})
	require.NoError(t, err)
	require.Empty(t, getResp.RbacPolicy.Actions)
}

func TestRBACPolicy_InvalidActionNames(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)

	cases := []struct {
		name   string
		action string
	}{
		{"missing dots", "foobar"},
		{"invalid chars", "foo.bar.BAZ"},
		{"reserved prefix", "tesseral.foo.bar"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := u.Store.UpdateRBACPolicy(ctx, &backendv1.UpdateRBACPolicyRequest{
				RbacPolicy: &backendv1.RBACPolicy{
					Actions: []*backendv1.Action{{Name: tc.action}},
				},
			})
			require.Error(t, err)
			var connectErr *connect.Error
			require.ErrorAs(t, err, &connectErr)
			require.Equal(t, connect.CodeInvalidArgument, connectErr.Code())
		})
	}
}

func TestRBACPolicy_UpsertIdempotency(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)

	_, err := u.Store.UpdateRBACPolicy(ctx, &backendv1.UpdateRBACPolicyRequest{
		RbacPolicy: &backendv1.RBACPolicy{
			Actions: []*backendv1.Action{{Name: "foo.bar.baz", Description: "desc"}},
		},
	})
	require.NoError(t, err)

	_, err = u.Store.UpdateRBACPolicy(ctx, &backendv1.UpdateRBACPolicyRequest{
		RbacPolicy: &backendv1.RBACPolicy{
			Actions: []*backendv1.Action{{Name: "foo.bar.baz", Description: "desc"}},
		},
	})
	require.NoError(t, err)
}
