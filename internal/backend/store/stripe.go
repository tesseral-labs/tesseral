package store

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/stripe/stripe-go/v82"
	"github.com/tesseral-labs/tesseral/internal/backend/authn"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
)

func (s *Store) GetProjectEntitlements(ctx context.Context, req *backendv1.GetProjectEntitlementsRequest) (*backendv1.GetProjectEntitlementsResponse, error) {
	qProject, err := s.q.GetProjectByID(ctx, authn.ProjectID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get project by id: %w", err)
	}

	return &backendv1.GetProjectEntitlementsResponse{
		EntitledCustomVaultDomains: qProject.EntitledCustomVaultDomains,
		EntitledBackendApiKeys:     qProject.EntitledBackendApiKeys,
	}, nil
}

func (s *Store) CreateStripeCheckoutLink(ctx context.Context, req *backendv1.CreateStripeCheckoutLinkRequest) (*backendv1.CreateStripeCheckoutLinkResponse, error) {
	if err := validateIsDogfoodSession(ctx); err != nil {
		return nil, fmt.Errorf("validate is dogfood session: %w", err)
	}

	qProject, err := s.q.GetProjectByID(ctx, authn.ProjectID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get project by id: %w", err)
	}

	if qProject.StripeCustomerID == nil {
		return nil, fmt.Errorf("project has no stripe customer id")
	}

	checkoutSession, err := s.stripe.CheckoutSessions.New(&stripe.CheckoutSessionParams{
		SuccessURL: stripe.String(fmt.Sprintf("https://%s/stripe-checkout-success", s.consoleDomain)),
		Mode:       stripe.String("subscription"),
		Customer:   qProject.StripeCustomerID,
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
