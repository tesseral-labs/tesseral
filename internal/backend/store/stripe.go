package store

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/stripe/stripe-go/v82"
	"github.com/tesseral-labs/tesseral/internal/backend/authn"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
)

func (s *Store) CreateStripeCheckoutLink(ctx context.Context, req *backendv1.CreateStripeCheckoutLinkRequest) (*backendv1.CreateStripeCheckoutLinkResponse, error) {
	qProject, err := s.q.GetProjectByID(ctx, authn.ProjectID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get project by id: %w", err)
	}

	if qProject.StripeCustomerID == nil {
		return nil, fmt.Errorf("project has no stripe customer id")
	}

	checkoutSession, err := s.stripe.CheckoutSessions.New(&stripe.CheckoutSessionParams{
		Customer: qProject.StripeCustomerID,
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    &s.stripePriceIDGrowthTier,
				Quantity: stripe.Int64(1),
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("create stripe checkout session: %w", err)
	}

	slog.InfoContext(ctx, "create_stripe_checkout_session", "customer_id", qProject.StripeCustomerID, "id", checkoutSession.ID)

	return &backendv1.CreateStripeCheckoutLinkResponse{
		Url: checkoutSession.URL,
	}, nil
}
