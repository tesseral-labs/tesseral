package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/custom_hostnames"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/tesseral-labs/tesseral/internal/backend/authn"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	"github.com/tesseral-labs/tesseral/internal/backend/store/queries"
	"github.com/tesseral-labs/tesseral/internal/cloudflaredoh"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

func (s *Store) GetVaultDomainSettings(ctx context.Context, req *backendv1.GetVaultDomainSettingsRequest) (*backendv1.GetVaultDomainSettingsResponse, error) {
	if err := validateIsDogfoodSession(ctx); err != nil {
		return nil, fmt.Errorf("validate is dogfood session: %w", err)
	}

	qVaultDomainSettings, err := s.q.GetVaultDomainSettings(ctx, authn.ProjectID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get vault domain settings: %w", err)
	}

	emailIdentity, err := s.ses.GetEmailIdentity(ctx, &sesv2.GetEmailIdentityInput{
		EmailIdentity: &qVaultDomainSettings.PendingDomain,
	})
	if err != nil {
		return nil, fmt.Errorf("get email identity: %w", err)
	}

	customHostname, err := s.getCloudflareCustomHostname(ctx, qVaultDomainSettings.PendingDomain)
	if err != nil {
		return nil, fmt.Errorf("get cloudflare custom hostname: %w", err)
	}

	projectVerificationRecordConfigured, err := s.getProjectVerificationRecordConfigured(ctx, qVaultDomainSettings.PendingDomain, authn.ProjectID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get project verification record configured: %w", err)
	}

	return &backendv1.GetVaultDomainSettingsResponse{
		VaultDomainSettings: s.parseVaultDomainSettings(qVaultDomainSettings, emailIdentity, customHostname, projectVerificationRecordConfigured),
	}, nil
}

func (s *Store) UpdateVaultDomainSettings(ctx context.Context, req *backendv1.UpdateVaultDomainSettingsRequest) (*backendv1.UpdateVaultDomainSettingsResponse, error) {
	if err := validateIsDogfoodSession(ctx); err != nil {
		return nil, fmt.Errorf("validate is dogfood session: %w", err)
	}

	previousPendingDomain, err := s.getCurrentPendingDomain(ctx)
	if err != nil {
		return nil, fmt.Errorf("get current pending domain: %w", err)
	}

	emailIdentity, err := s.upsertSESEmailIdentity(ctx, req.VaultDomainSettings.PendingDomain)
	if err != nil {
		return nil, fmt.Errorf("upsert ses email identity: %w", err)
	}

	customHostname, err := s.upsertCloudflareCustomHostname(ctx, req.VaultDomainSettings.PendingDomain)
	if err != nil {
		return nil, fmt.Errorf("upsert cloudflare custom hostname: %w", err)
	}

	projectVerificationRecordConfigured, err := s.getProjectVerificationRecordConfigured(ctx, req.VaultDomainSettings.PendingDomain, authn.ProjectID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get project verification record configured: %w", err)
	}

	qVaultDomainSettings, err := s.q.UpsertVaultDomainSettings(ctx, queries.UpsertVaultDomainSettingsParams{
		ProjectID:     authn.ProjectID(ctx),
		PendingDomain: req.VaultDomainSettings.PendingDomain,
	})
	if err != nil {
		return nil, fmt.Errorf("upsert vault domain settings: %w", err)
	}

	// delete resources associated with previous pending domain, if not in use
	previousDomainInUse, err := s.q.GetVaultDomainInActiveOrPendingUse(ctx, &previousPendingDomain)
	if err != nil {
		return nil, fmt.Errorf("get vault domain in active or pending use: %w", err)
	}

	if !previousDomainInUse.Valid {
		panic("null from GetVaultDomainInActiveOrPendingUse")
	}

	if !previousDomainInUse.Bool {
		if _, err := s.ses.DeleteEmailIdentity(ctx, &sesv2.DeleteEmailIdentityInput{
			EmailIdentity: &previousPendingDomain,
		}); err != nil {
			return nil, fmt.Errorf("delete email identity: %w", err)
		}

		previousCustomHostname, err := s.getCloudflareCustomHostname(ctx, previousPendingDomain)
		if err != nil {
			return nil, fmt.Errorf("get cloudflare custom hostname: %w", err)
		}

		if _, err := s.cloudflare.CustomHostnames.Delete(ctx, previousCustomHostname.ID, custom_hostnames.CustomHostnameDeleteParams{
			ZoneID: cloudflare.F(s.tesseralDNSCloudflareZoneID),
		}); err != nil {
			return nil, fmt.Errorf("delete cloudflare custom hostname: %w", err)
		}
	}

	return &backendv1.UpdateVaultDomainSettingsResponse{
		VaultDomainSettings: s.parseVaultDomainSettings(qVaultDomainSettings, emailIdentity, customHostname, projectVerificationRecordConfigured),
	}, nil
}

func (s *Store) getCurrentPendingDomain(ctx context.Context) (string, error) {
	qVaultDomainSettings, err := s.q.GetVaultDomainSettings(ctx, authn.ProjectID(ctx))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", nil
		}
		return "", fmt.Errorf("get vault domain settings: %w", err)
	}

	return qVaultDomainSettings.PendingDomain, nil
}

func (s *Store) upsertSESEmailIdentity(ctx context.Context, emailIdentity string) (*sesv2.GetEmailIdentityOutput, error) {
	_, err := s.ses.CreateEmailIdentity(ctx, &sesv2.CreateEmailIdentityInput{
		EmailIdentity: &emailIdentity,
	})
	if err != nil {
		var alreadyExists *types.AlreadyExistsException
		if !errors.As(err, &alreadyExists) {
			return nil, fmt.Errorf("create email identity: %w", err)
		}
		// deliberate fallthrough if already exists
	}

	getEmailIdentityRes, err := s.ses.GetEmailIdentity(ctx, &sesv2.GetEmailIdentityInput{
		EmailIdentity: &emailIdentity,
	})
	if err != nil {
		return nil, fmt.Errorf("get email identity: %w", err)
	}
	return getEmailIdentityRes, nil
}

func (s *Store) getCloudflareCustomHostname(ctx context.Context, hostname string) (*cloudflareCustomHostname, error) {
	res, err := s.cloudflare.CustomHostnames.List(ctx, custom_hostnames.CustomHostnameListParams{
		ZoneID:   cloudflare.F(s.tesseralDNSCloudflareZoneID),
		Hostname: cloudflare.F(hostname),
	})
	if err != nil {
		return nil, fmt.Errorf("list cloudflare custom hostnames: %w", err)
	}

	if len(res.Result) != 0 {
		panic(fmt.Errorf("exactly one custom hostname expected"))
	}

	return &cloudflareCustomHostname{
		ID:     res.Result[0].ID,
		Status: string(res.Result[0].Status),
	}, nil
}

func (s *Store) upsertCloudflareCustomHostname(ctx context.Context, hostname string) (*cloudflareCustomHostname, error) {
	listRes, err := s.cloudflare.CustomHostnames.List(ctx, custom_hostnames.CustomHostnameListParams{
		ZoneID:   cloudflare.F(s.tesseralDNSCloudflareZoneID),
		Hostname: cloudflare.F(hostname),
	})
	if err != nil {
		return nil, fmt.Errorf("list cloudflare custom hostnames: %w", err)
	}

	if len(listRes.Result) != 0 {
		return &cloudflareCustomHostname{
			ID:     listRes.Result[0].ID,
			Status: string(listRes.Result[0].Status),
		}, nil
	}

	customHostname, err := s.cloudflare.CustomHostnames.New(ctx, custom_hostnames.CustomHostnameNewParams{
		ZoneID:   cloudflare.F(s.tesseralDNSCloudflareZoneID),
		Hostname: cloudflare.F(hostname),
		SSL: cloudflare.F(custom_hostnames.CustomHostnameNewParamsSSL{
			Method: cloudflare.F(custom_hostnames.DCVMethodHTTP),
			Type:   cloudflare.F(custom_hostnames.DomainValidationTypeDv),
		}),
	})
	if err != nil {
		return nil, fmt.Errorf("create cloudflare custom hostname: %w", err)
	}

	return &cloudflareCustomHostname{
		ID:     customHostname.ID,
		Status: string(customHostname.Status),
	}, nil
}

func (s *Store) getProjectVerificationRecordConfigured(ctx context.Context, pendingDomain string, projectID uuid.UUID) (bool, error) {
	res, err := s.cloudflareDOH.DNSQuery(ctx, &cloudflaredoh.DNSQueryRequest{
		Name: fmt.Sprintf("_tesseral_project_verification.%s", pendingDomain),
		Type: "TXT",
	})
	if err != nil {
		return false, fmt.Errorf("get project verification record configured: %w", err)
	}

	if len(res.Answer) != 1 {
		return false, nil
	}

	return res.Answer[0].Data == idformat.Project.Format(projectID), nil
}

// cloudflareCustomHostname exists to unify Cloudflare's SDK return types for
// List/New on custom hostnames.
type cloudflareCustomHostname struct {
	ID     string
	Status string
}

func (s *Store) parseVaultDomainSettings(qVaultDomainSettings queries.VaultDomainSetting, emailIdentity *sesv2.GetEmailIdentityOutput, customHostname *cloudflareCustomHostname, projectVerificationRecordConfigured bool) *backendv1.VaultDomainSettings {
	var dkimRecords []*backendv1.VaultDomainSettingsDNSRecord
	for _, token := range emailIdentity.DkimAttributes.Tokens {
		dkimRecords = append(dkimRecords, &backendv1.VaultDomainSettingsDNSRecord{
			Type:  "CNAME",
			Name:  fmt.Sprintf("%s._domainkey.%s", token, qVaultDomainSettings.PendingDomain),
			Value: fmt.Sprintf("%s.dkim.amazonses.com", token),
		})
	}

	return &backendv1.VaultDomainSettings{
		PendingDomain: qVaultDomainSettings.PendingDomain,
		MainRecord: &backendv1.VaultDomainSettingsDNSRecord{
			Type:  "CNAME",
			Name:  qVaultDomainSettings.PendingDomain,
			Value: s.tesseralDNSVaultCNAMEValue,
		},
		MainRecordConfigured: customHostname.Status == string(custom_hostnames.CustomHostnameListResponseStatusActive),
		DkimRecords:          dkimRecords,
		DkimConfigured:       emailIdentity.DkimAttributes.Status == types.DkimStatusSuccess,
		SpfRecords: []*backendv1.VaultDomainSettingsDNSRecord{
			{
				Type:  "MX",
				Name:  fmt.Sprintf("mail.%s", qVaultDomainSettings.PendingDomain),
				Value: s.sesSPFMXRecordValue,
			},
			{
				Type:  "TXT",
				Name:  fmt.Sprintf("mail.%s", qVaultDomainSettings.PendingDomain),
				Value: "v=spf1 include:amazonses.com ~all",
			},
		},
		SpfConfigured: emailIdentity.MailFromAttributes.MailFromDomainStatus == types.MailFromDomainStatusSuccess,
		ProjectVerificationRecord: &backendv1.VaultDomainSettingsDNSRecord{
			Type:  "TXT",
			Name:  fmt.Sprintf("_tesseral_project_verification.%s", qVaultDomainSettings.PendingDomain),
			Value: idformat.Project.Format(qVaultDomainSettings.ProjectID),
		},
		ProjectVerificationRecordConfigured: projectVerificationRecordConfigured,
	}
}
