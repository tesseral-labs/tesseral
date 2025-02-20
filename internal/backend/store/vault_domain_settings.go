package store

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	"github.com/cloudflare/cloudflare-go"
	"github.com/tesseral-labs/tesseral/internal/backend/authn"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
)

func (s *Store) GetVaultDomainSettings(ctx context.Context, req *backendv1.GetVaultDomainSettingsRequest) (*backendv1.GetVaultDomainSettingsResponse, error) {
	if err := validateIsDogfoodSession(ctx); err != nil {
		return nil, fmt.Errorf("validate is dogfood session: %w", err)
	}

	qVaultDomainSettings, err := s.q.GetVaultDomainSettings(ctx, authn.ProjectID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get vault domain settings: %w", err)
	}

	getEmailIdentityRes, err := s.ses.GetEmailIdentity(ctx, &sesv2.GetEmailIdentityInput{
		EmailIdentity: &qVaultDomainSettings.PendingDomain,
	})
	if err != nil {
		return nil, fmt.Errorf("get email identity: %w", err)
	}

	var dkimRecords []*backendv1.VaultDomainSettingsDNSRecord
	for _, token := range getEmailIdentityRes.DkimAttributes.Tokens {
		dkimRecords = append(dkimRecords, &backendv1.VaultDomainSettingsDNSRecord{
			Type:  "CNAME",
			Name:  fmt.Sprintf("%s._domainkey.%s", token, qVaultDomainSettings.PendingDomain),
			Value: fmt.Sprintf("%s.dkim.amazonses.com", token),
		})
	}

	cloudflareHostnameID, err := s.cloudflare.CustomHostnameIDByName(ctx, s.tesseralDNSCloudflareZoneID, qVaultDomainSettings.PendingDomain)
	if err != nil {
		return nil, fmt.Errorf("get cloudflare hostname id: %w", err)
	}

	customHostname, err := s.cloudflare.CustomHostname(ctx, s.tesseralDNSCloudflareZoneID, cloudflareHostnameID)
	if err != nil {
		return nil, fmt.Errorf("get cloudflare hostname: %w", err)
	}

	return &backendv1.GetVaultDomainSettingsResponse{
		VaultDomainSettings: &backendv1.VaultDomainSettings{
			PendingDomain: qVaultDomainSettings.PendingDomain,

			// todo validate this more thoroughly -- what does cloudflare return if they've provisioned the cert, but the CNAME isn't yet in place?
			TesseralCnameRecordConfigured: customHostname.Status == cloudflare.ACTIVE,
			TesseralCnameRecord: &backendv1.VaultDomainSettingsDNSRecord{
				Type:  "CNAME",
				Name:  qVaultDomainSettings.PendingDomain,
				Value: s.tesseralDNSVaultCNAMEValue,
			},

			DkimRecords:    dkimRecords,
			DkimConfigured: getEmailIdentityRes.DkimAttributes.Status == types.DkimStatusSuccess,

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
			SpfConfigured: getEmailIdentityRes.MailFromAttributes.MailFromDomainStatus == types.MailFromDomainStatusSuccess,
		},
	}, nil
}
