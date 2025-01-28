package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/openauth/openauth/internal/backend/authn"
	backendv1 "github.com/openauth/openauth/internal/backend/gen/openauth/backend/v1"
	"github.com/openauth/openauth/internal/backend/store/queries"
	"github.com/openauth/openauth/internal/common/apierror"
	"github.com/openauth/openauth/internal/store/idformat"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Store) CreateOrganization(ctx context.Context, req *backendv1.CreateOrganizationRequest) (*backendv1.CreateOrganizationResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	// TODO these are a list now
	//var googleHostedDomain *string
	//if req.Organization.GoogleHostedDomain != "" {
	//	googleHostedDomain = &req.Organization.GoogleHostedDomain
	//}
	//
	//var microsoftTenantId *string
	//if req.Organization.MicrosoftTenantId != "" {
	//	microsoftTenantId = &req.Organization.MicrosoftTenantId
	//}

	var (
		disableLogInWithGoogle,
		disableLogInWithMicrosoft,
		disableLogInWithPassword *bool
	)

	if req.Organization.OverrideLogInMethods != nil {
		disableLogInWithGoogle = aws.Bool(!req.Organization.LogInWithGoogleEnabled)
		disableLogInWithMicrosoft = aws.Bool(!req.Organization.LogInWithMicrosoftEnabled)
		disableLogInWithPassword = aws.Bool(!req.Organization.LogInWithPasswordEnabled)
	}

	qProject, err := q.GetProjectByID(ctx, authn.ProjectID(ctx))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("project not found", fmt.Errorf("get project by id: %w", err))
		}

		return nil, fmt.Errorf("get project by id: %w", err)
	}

	var samlEnabled bool
	if req.Organization.SamlEnabled != nil {
		samlEnabled = *req.Organization.SamlEnabled
	}

	var scimEnabled bool
	if req.Organization.ScimEnabled != nil {
		scimEnabled = *req.Organization.ScimEnabled
	}

	qOrg, err := q.CreateOrganization(ctx, queries.CreateOrganizationParams{
		ID:          uuid.New(),
		ProjectID:   authn.ProjectID(ctx),
		DisplayName: req.Organization.DisplayName,
		// TODO these are a list now
		//GoogleHostedDomain: googleHostedDomain,
		//MicrosoftTenantID:  microsoftTenantId,

		OverrideLogInMethods:      derefOrEmpty(req.Organization.OverrideLogInMethods),
		DisableLogInWithGoogle:    disableLogInWithGoogle,
		DisableLogInWithMicrosoft: disableLogInWithMicrosoft,
		DisableLogInWithPassword:  disableLogInWithPassword,

		SamlEnabled: samlEnabled,
		ScimEnabled: scimEnabled,
	})
	if err != nil {
		return nil, fmt.Errorf("create organization: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &backendv1.CreateOrganizationResponse{
		Organization: parseOrganization(parseOrganizationArgs{
			qProject:             qProject,
			qOrg:                 qOrg,
			qGoogleHostedDomains: nil, // TODO
			qMicrosoftTenantIDs:  nil, // TODO
		}),
	}, nil
}

func (s *Store) ListOrganizations(ctx context.Context, req *backendv1.ListOrganizationsRequest) (*backendv1.ListOrganizationsResponse, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	var startID uuid.UUID
	if err := s.pageEncoder.Unmarshal(req.PageToken, &startID); err != nil {
		return nil, err
	}

	qProject, err := q.GetProjectByID(ctx, authn.ProjectID(ctx))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("project not found", fmt.Errorf("get project by id: %w", err))
		}

		return nil, fmt.Errorf("get project by id: %w", err)
	}

	limit := 10
	qOrgs, err := q.ListOrganizationsByProjectId(ctx, queries.ListOrganizationsByProjectIdParams{
		ProjectID: authn.ProjectID(ctx),
		Limit:     int32(limit + 1),
	})
	if err != nil {
		return nil, fmt.Errorf("list organizations: %w", err)
	}

	var organizations []*backendv1.Organization
	for _, qOrg := range qOrgs {
		organizations = append(organizations, parseOrganization(parseOrganizationArgs{
			qProject:             qProject,
			qOrg:                 qOrg,
			qGoogleHostedDomains: nil, // TODO
			qMicrosoftTenantIDs:  nil, // TODO
		}))
	}

	var nextPageToken string
	if len(organizations) == limit+1 {
		nextPageToken = s.pageEncoder.Marshal(organizations[limit].Id)
		organizations = organizations[:limit]
	}

	return &backendv1.ListOrganizationsResponse{
		Organizations: organizations,
		NextPageToken: nextPageToken,
	}, nil
}

func (s *Store) GetOrganization(ctx context.Context, req *backendv1.GetOrganizationRequest) (*backendv1.GetOrganizationResponse, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	organizationId, err := idformat.Organization.Parse(req.Id)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid organization id", fmt.Errorf("parse organization id: %w", err))
	}

	qProject, err := q.GetProjectByID(ctx, authn.ProjectID(ctx))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("project not found", fmt.Errorf("get project by id: %w", err))
		}

		return nil, fmt.Errorf("get project by id: %w", err)
	}

	qOrg, err := q.GetOrganizationByProjectIDAndID(ctx, queries.GetOrganizationByProjectIDAndIDParams{
		ProjectID: authn.ProjectID(ctx),
		ID:        organizationId,
	})
	if err != nil {
		return nil, fmt.Errorf("get organization: %w", err)
	}

	return &backendv1.GetOrganizationResponse{Organization: parseOrganization(parseOrganizationArgs{
		qProject:             qProject,
		qOrg:                 qOrg,
		qGoogleHostedDomains: nil, // TODO
		qMicrosoftTenantIDs:  nil, // TODO
	})}, nil
}

func (s *Store) UpdateOrganization(ctx context.Context, req *backendv1.UpdateOrganizationRequest) (*backendv1.UpdateOrganizationResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	orgID, err := idformat.Organization.Parse(req.Id)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid organization id", fmt.Errorf("parse organization id: %w", err))
	}

	// fetch existing org; this also acts as a permission check
	qOrg, err := q.GetOrganizationByProjectIDAndID(ctx, queries.GetOrganizationByProjectIDAndIDParams{
		ProjectID: authn.ProjectID(ctx),
		ID:        orgID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("organization not found", fmt.Errorf("get organization by id: %w", err))
		}

		return nil, fmt.Errorf("get organization: %w", err)
	}

	updates := queries.UpdateOrganizationParams{
		ID: orgID,
	}

	updates.DisplayName = qOrg.DisplayName
	if req.Organization.DisplayName != "" {
		updates.DisplayName = req.Organization.DisplayName
	}

	updates.OverrideLogInMethods = qOrg.OverrideLogInMethods
	if req.Organization.OverrideLogInMethods != nil {
		updates.OverrideLogInMethods = *req.Organization.OverrideLogInMethods
	}

	// update the override_log_in_with_..._enabled columns to null unless the
	// organization is overriding those columns.
	if req.Organization.GetOverrideLogInMethods() {
		updates.DisableLogInWithGoogle = aws.Bool(!req.Organization.LogInWithGoogleEnabled)
		updates.DisableLogInWithMicrosoft = aws.Bool(!req.Organization.LogInWithMicrosoftEnabled)
		updates.DisableLogInWithPassword = aws.Bool(!req.Organization.LogInWithPasswordEnabled)
	}

	updates.SamlEnabled = qOrg.SamlEnabled
	if req.Organization.SamlEnabled != nil {
		updates.SamlEnabled = *req.Organization.SamlEnabled
	}

	updates.ScimEnabled = qOrg.ScimEnabled
	if req.Organization.ScimEnabled != nil {
		updates.ScimEnabled = *req.Organization.ScimEnabled
	}

	qUpdatedOrg, err := q.UpdateOrganization(ctx, updates)
	if err != nil {
		return nil, fmt.Errorf("update organization: %w", err)
	}

	fmt.Println("update org", req.Organization.GoogleHostedDomains, len(req.Organization.GoogleHostedDomains), req.Organization.GoogleHostedDomains == nil)

	qProject, err := q.GetProjectByID(ctx, authn.ProjectID(ctx))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("project not found", fmt.Errorf("get project by id: %w", err))
		}

		return nil, fmt.Errorf("get project by id: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &backendv1.UpdateOrganizationResponse{Organization: parseOrganization(parseOrganizationArgs{
		qProject:             qProject,
		qOrg:                 qUpdatedOrg,
		qGoogleHostedDomains: nil, // TODO
		qMicrosoftTenantIDs:  nil, // TODO
	})}, nil
}

func (s *Store) DeleteOrganization(ctx context.Context, req *backendv1.DeleteOrganizationRequest) (*backendv1.DeleteOrganizationResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	orgID, err := idformat.Organization.Parse(req.Id)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid organization id", fmt.Errorf("parse organization id: %w", err))
	}

	// authz check
	if _, err := q.GetOrganizationByProjectIDAndID(ctx, queries.GetOrganizationByProjectIDAndIDParams{
		ProjectID: authn.ProjectID(ctx),
		ID:        orgID,
	}); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("organization not found", fmt.Errorf("get organization by id: %w", err))
		}

		return nil, fmt.Errorf("get organization: %w", err)
	}

	if err := q.DeleteOrganization(ctx, orgID); err != nil {
		return nil, fmt.Errorf("delete organization: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &backendv1.DeleteOrganizationResponse{}, nil
}

func (s *Store) DisableOrganizationLogins(ctx context.Context, req *backendv1.DisableOrganizationLoginsRequest) (*backendv1.DisableOrganizationLoginsResponse, error) {
	if err := validateIsDogfoodSession(ctx); err != nil {
		return nil, fmt.Errorf("validate is dogfood session: %w", err)
	}

	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	if err := q.DisableOrganizationLogins(ctx, authn.ProjectID(ctx)); err != nil {
		return nil, fmt.Errorf("lockout organization: %w", err)
	}

	if err := q.RevokeAllOrganizationSessions(ctx, authn.ProjectID(ctx)); err != nil {
		return nil, fmt.Errorf("revoke all organization sessions: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &backendv1.DisableOrganizationLoginsResponse{}, nil
}

func (s *Store) EnableOrganizationLogins(ctx context.Context, req *backendv1.EnableOrganizationLoginsRequest) (*backendv1.EnableOrganizationLoginsResponse, error) {
	if err := validateIsDogfoodSession(ctx); err != nil {
		return nil, fmt.Errorf("validate is dogfood session: %w", err)
	}

	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	if err := q.EnableOrganizationLogins(ctx, authn.ProjectID(ctx)); err != nil {
		return nil, fmt.Errorf("unlock organization: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &backendv1.EnableOrganizationLoginsResponse{}, nil
}

type parseOrganizationArgs struct {
	qProject             queries.Project
	qOrg                 queries.Organization
	qGoogleHostedDomains []queries.OrganizationGoogleHostedDomain
	qMicrosoftTenantIDs  []queries.OrganizationMicrosoftTenantID
}

func parseOrganization(args parseOrganizationArgs) *backendv1.Organization {
	// sanity-check consistency of args
	if args.qProject.ID != args.qOrg.ProjectID {
		panic("project id mismatch")
	}
	for _, qGoogleHostedDomain := range args.qGoogleHostedDomains {
		if qGoogleHostedDomain.OrganizationID != args.qOrg.ID {
			panic("google hosted domain organization id mismatch")
		}
	}
	for _, qMicrosoftTenantID := range args.qMicrosoftTenantIDs {
		if qMicrosoftTenantID.OrganizationID != args.qOrg.ID {
			panic("microsoft tenant id organization id mismatch")
		}
	}

	logInWithGoogleEnabled := args.qProject.LogInWithGoogleEnabled
	logInWithMicrosoftEnabled := args.qProject.LogInWithMicrosoftEnabled
	logInWithPasswordEnabled := args.qProject.LogInWithPasswordEnabled

	// allow orgs to disable login methods
	if derefOrEmpty(args.qOrg.DisableLogInWithGoogle) {
		logInWithGoogleEnabled = false
	}
	if derefOrEmpty(args.qOrg.DisableLogInWithMicrosoft) {
		logInWithMicrosoftEnabled = false
	}
	if derefOrEmpty(args.qOrg.DisableLogInWithPassword) {
		logInWithPasswordEnabled = false
	}

	var googleHostedDomains []string
	for _, qGoogleHostedDomain := range args.qGoogleHostedDomains {
		googleHostedDomains = append(googleHostedDomains, qGoogleHostedDomain.GoogleHostedDomain)
	}

	var microsoftTenantIDs []string
	for _, qMicrosoftTenantID := range args.qMicrosoftTenantIDs {
		microsoftTenantIDs = append(microsoftTenantIDs, qMicrosoftTenantID.MicrosoftTenantID)
	}

	return &backendv1.Organization{
		Id:                        idformat.Organization.Format(args.qOrg.ID),
		DisplayName:               args.qOrg.DisplayName,
		CreateTime:                timestamppb.New(*args.qOrg.CreateTime),
		UpdateTime:                timestamppb.New(*args.qOrg.UpdateTime),
		OverrideLogInMethods:      &args.qOrg.OverrideLogInMethods,
		LogInWithPasswordEnabled:  logInWithPasswordEnabled,
		LogInWithGoogleEnabled:    logInWithGoogleEnabled,
		LogInWithMicrosoftEnabled: logInWithMicrosoftEnabled,
		GoogleHostedDomains:       googleHostedDomains,
		MicrosoftTenantIds:        microsoftTenantIDs,
		SamlEnabled:               &args.qOrg.SamlEnabled,
		ScimEnabled:               &args.qOrg.ScimEnabled,
	}
}
