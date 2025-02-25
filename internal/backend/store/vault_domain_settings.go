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

	vaultDomainSettings, err := s.getVaultDomainSettings(ctx)
	if err != nil {
		return nil, fmt.Errorf("get vault domain settings: %w", err)
	}

	return &backendv1.GetVaultDomainSettingsResponse{
		VaultDomainSettings: vaultDomainSettings,
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

	qVaultDomainSettings, err := s.q.UpsertVaultDomainSettings(ctx, queries.UpsertVaultDomainSettingsParams{
		ProjectID:     authn.ProjectID(ctx),
		PendingDomain: req.VaultDomainSettings.PendingDomain,
	})
	if err != nil {
		return nil, fmt.Errorf("upsert vault domain settings: %w", err)
	}

	if previousPendingDomain != "" {
		// delete resources associated with previous pending domain, if not in use
		previousDomainInUse, err := s.q.GetVaultDomainInActiveOrPendingUse(ctx, previousPendingDomain)
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

			previousCustomHostname, err := s.getCloudflareCustomHostnameByHostname(ctx, previousPendingDomain)
			if err != nil {
				return nil, fmt.Errorf("get cloudflare custom hostname: %w", err)
			}

			if _, err := s.cloudflare.CustomHostnames.Delete(ctx, previousCustomHostname.ID, custom_hostnames.CustomHostnameDeleteParams{
				ZoneID: cloudflare.F(s.tesseralDNSCloudflareZoneID),
			}); err != nil {
				return nil, fmt.Errorf("delete cloudflare custom hostname: %w", err)
			}
		}
	}

	vaultDomainSettings, err := s.parseVaultDomainSettings(ctx, qVaultDomainSettings, emailIdentity, customHostname)
	if err != nil {
		return nil, fmt.Errorf("parse vault domain settings: %w", err)
	}

	return &backendv1.UpdateVaultDomainSettingsResponse{
		VaultDomainSettings: vaultDomainSettings,
	}, nil
}

func (s *Store) EnableCustomVaultDomain(ctx context.Context, req *backendv1.EnableCustomVaultDomainRequest) (*backendv1.EnableCustomVaultDomainResponse, error) {
	if err := validateIsDogfoodSession(ctx); err != nil {
		return nil, fmt.Errorf("validate is dogfood session: %w", err)
	}

	vaultDomainSettings, err := s.getVaultDomainSettings(ctx)
	if err != nil {
		return nil, fmt.Errorf("get vault domain settings: %w", err)
	}

	if !vaultDomainSettings.PendingVaultDomainReady {
		return nil, fmt.Errorf("vault domain not ready")
	}

	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	if _, err := q.UpdateProjectVaultDomain(ctx, queries.UpdateProjectVaultDomainParams{
		ID:          authn.ProjectID(ctx),
		VaultDomain: vaultDomainSettings.PendingDomain,
	}); err != nil {
		return nil, fmt.Errorf("update project vault domain: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &backendv1.EnableCustomVaultDomainResponse{}, nil
}

func (s *Store) EnableEmailSendFromDomain(ctx context.Context, req *backendv1.EnableEmailSendFromDomainRequest) (*backendv1.EnableEmailSendFromDomainResponse, error) {
	if err := validateIsDogfoodSession(ctx); err != nil {
		return nil, fmt.Errorf("validate is dogfood session: %w", err)
	}

	vaultDomainSettings, err := s.getVaultDomainSettings(ctx)
	if err != nil {
		return nil, fmt.Errorf("get vault domain settings: %w", err)
	}

	if !vaultDomainSettings.PendingSendFromDomainReady {
		return nil, fmt.Errorf("email send-from domain not ready")
	}

	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	if _, err := q.UpdateProjectEmailSendFromDomain(ctx, queries.UpdateProjectEmailSendFromDomainParams{
		ID:                  authn.ProjectID(ctx),
		EmailSendFromDomain: fmt.Sprintf("mail.%s", vaultDomainSettings.PendingDomain),
	}); err != nil {
		return nil, fmt.Errorf("update project email send-from domain: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &backendv1.EnableEmailSendFromDomainResponse{}, nil
}

func (s *Store) getVaultDomainSettings(ctx context.Context) (*backendv1.VaultDomainSettings, error) {
	qVaultDomainSettings, err := s.q.GetVaultDomainSettings(ctx, authn.ProjectID(ctx))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("get vault domain settings: %w", err)
	}

	emailIdentity, err := s.ses.GetEmailIdentity(ctx, &sesv2.GetEmailIdentityInput{
		EmailIdentity: &qVaultDomainSettings.PendingDomain,
	})
	if err != nil {
		return nil, fmt.Errorf("get email identity: %w", err)
	}

	currentCustomHostname, err := s.getCloudflareCustomHostnameByHostname(ctx, qVaultDomainSettings.PendingDomain)
	if err != nil {
		return nil, fmt.Errorf("get cloudflare custom hostname: %w", err)
	}

	// issue a no-op edit on the custom hostname, to get cloudflare to check if
	// the customer's CNAME is in place
	refreshedCustomHostname, err := s.cloudflare.CustomHostnames.Edit(ctx, currentCustomHostname.ID, custom_hostnames.CustomHostnameEditParams{
		ZoneID: cloudflare.F(s.tesseralDNSCloudflareZoneID),
	})
	if err != nil {
		return nil, fmt.Errorf("edit cloudflare custom hostname: %w", err)
	}

	vaultDomainSettings, err := s.parseVaultDomainSettings(ctx, qVaultDomainSettings, emailIdentity, &cloudflareCustomHostname{
		ID:     refreshedCustomHostname.ID,
		Status: string(refreshedCustomHostname.Status),
	})
	if err != nil {
		return nil, fmt.Errorf("parse vault domain settings: %w", err)
	}

	return vaultDomainSettings, nil
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

func (s *Store) getCloudflareCustomHostnameByHostname(ctx context.Context, hostname string) (*custom_hostnames.CustomHostnameListResponse, error) {
	res, err := s.cloudflare.CustomHostnames.List(ctx, custom_hostnames.CustomHostnameListParams{
		ZoneID:   cloudflare.F(s.tesseralDNSCloudflareZoneID),
		Hostname: cloudflare.F(hostname),
	})
	if err != nil {
		return nil, fmt.Errorf("list cloudflare custom hostnames: %w", err)
	}

	if len(res.Result) != 1 {
		panic(fmt.Errorf("exactly one custom hostname expected"))
	}

	return &res.Result[0], nil
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

func (s *Store) addActualValuesToDNSRecord(ctx context.Context, dnsRecord *backendv1.VaultDomainSettingsDNSRecord) (*backendv1.VaultDomainSettingsDNSRecord, error) {
	res, err := s.cloudflareDOH.DNSQuery(ctx, &cloudflaredoh.DNSQueryRequest{
		Name: dnsRecord.Name,
		Type: dnsRecord.Type,
	})
	if err != nil {
		return nil, fmt.Errorf("dns query: %w", err)
	}

	var recordType int32
	switch dnsRecord.Type {
	case "CNAME":
		recordType = 5
	case "MX":
		recordType = 15
	case "TXT":
		recordType = 16
	}

	var values []string
	var ttl uint32
	for _, answer := range res.Answer {
		if answer.Name != dnsRecord.Name || answer.Type != recordType {
			continue // this is just a related record
		}

		values = append(values, answer.Data)
		if answer.TTL > ttl {
			ttl = answer.TTL
		}
	}

	return &backendv1.VaultDomainSettingsDNSRecord{
		Type:             dnsRecord.Type,
		Name:             dnsRecord.Name,
		WantValue:        dnsRecord.WantValue,
		ActualValues:     values,
		ActualTtlSeconds: ttl,
		Correct:          len(values) == 1 && values[0] == dnsRecord.WantValue,
	}, nil
}

func (s *Store) parseVaultDomainSettings(ctx context.Context, qVaultDomainSettings queries.VaultDomainSetting, emailIdentity *sesv2.GetEmailIdentityOutput, customHostname *cloudflareCustomHostname) (*backendv1.VaultDomainSettings, error) {
	vaultDomainRecords := []*backendv1.VaultDomainSettingsDNSRecord{
		{
			Type:      "CNAME",
			Name:      qVaultDomainSettings.PendingDomain,
			WantValue: s.tesseralDNSVaultCNAMEValue,
		},
		{
			Type:      "TXT",
			Name:      fmt.Sprintf("_tesseral_project_verification.%s", qVaultDomainSettings.PendingDomain),
			WantValue: fmt.Sprintf("\"%s\"", idformat.Project.Format(qVaultDomainSettings.ProjectID)),
		},
	}

	emailSendFromRecords := []*backendv1.VaultDomainSettingsDNSRecord{
		{
			Type:      "MX",
			Name:      fmt.Sprintf("mail.%s", qVaultDomainSettings.PendingDomain),
			WantValue: s.sesSPFMXRecordValue,
		},
		{
			Type:      "TXT",
			Name:      fmt.Sprintf("mail.%s", qVaultDomainSettings.PendingDomain),
			WantValue: "\"v=spf1 include:amazonses.com ~all\"",
		},
	}

	for _, token := range emailIdentity.DkimAttributes.Tokens {
		emailSendFromRecords = append(emailSendFromRecords, &backendv1.VaultDomainSettingsDNSRecord{
			Type:      "CNAME",
			Name:      fmt.Sprintf("%s._domainkey.%s", token, qVaultDomainSettings.PendingDomain),
			WantValue: fmt.Sprintf("%s.dkim.amazonses.com.", token),
		})
	}

	for i := range vaultDomainRecords {
		var err error
		vaultDomainRecords[i], err = s.addActualValuesToDNSRecord(ctx, vaultDomainRecords[i])
		if err != nil {
			return nil, fmt.Errorf("add actual values to dns record: %w", err)
		}
	}

	for i := range emailSendFromRecords {
		var err error
		emailSendFromRecords[i], err = s.addActualValuesToDNSRecord(ctx, emailSendFromRecords[i])
		if err != nil {
			return nil, fmt.Errorf("add actual values to dns record: %w", err)
		}
	}

	cloudflareOK := customHostname.Status == string(custom_hostnames.CustomHostnameListResponseStatusActive)
	emailIdentityOK := emailIdentity.VerificationStatus == types.VerificationStatusSuccess
	dkimOK := emailIdentity.DkimAttributes.Status == types.DkimStatusSuccess

	return &backendv1.VaultDomainSettings{
		PendingDomain:              qVaultDomainSettings.PendingDomain,
		PendingVaultDomainReady:    cloudflareOK,
		PendingSendFromDomainReady: emailIdentityOK && dkimOK,
		VaultDomainRecords:         vaultDomainRecords,
		EmailSendFromRecords:       emailSendFromRecords,
	}, nil
}
