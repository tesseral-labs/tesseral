import React from 'react';
import {
  ConsoleCard,
  ConsoleCardContent,
  ConsoleCardDescription,
  ConsoleCardDetails,
  ConsoleCardHeader,
  ConsoleCardTitle,
} from '@/components/ui/console-card';
import { useParams } from 'react-router';
import { useQuery } from '@connectrpc/connect-query';
import {
  getOrganization,
  getOrganizationDomains,
  getOrganizationGoogleHostedDomains,
  getOrganizationMicrosoftTenantIDs,
  getProject,
} from '@/gen/tesseral/backend/v1/backend-BackendService_connectquery';
import { Button } from '@/components/ui/button';
import { Link } from 'react-router-dom';
import {
  DetailsGrid,
  DetailsGridColumn,
  DetailsGridEntry,
  DetailsGridKey,
  DetailsGridValue,
} from '@/components/details-grid';
import { EditOrganizationGoogleConfigurationButton } from '@/pages/organizations/EditOrganizationGoogleConfigurationButton';
import { EditOrganizationMicrosoftConfigurationButton } from '@/pages/organizations/EditOrganizationMicrosoftConfigurationButton';

export const OrganizationDetailsTab = () => {
  const { organizationId } = useParams();
  const { data: getOrganizationResponse } = useQuery(getOrganization, {
    id: organizationId,
  });
  const { data: getProjectResponse } = useQuery(getProject, {});
  const { data: getOrganizationDomainsResponse } = useQuery(
    getOrganizationDomains,
    {
      organizationId,
    },
  );
  const { data: getOrganizationGoogleHostedDomainsResponse } = useQuery(
    getOrganizationGoogleHostedDomains,
    {
      organizationId,
    },
  );
  const { data: getOrganizationMicrosoftTenantIdsResponse } = useQuery(
    getOrganizationMicrosoftTenantIDs,
    {
      organizationId,
    },
  );

  return (
    <div className="space-y-8">
      <ConsoleCard>
        <ConsoleCardHeader>
          <ConsoleCardDetails>
            <ConsoleCardTitle>Details</ConsoleCardTitle>
            <ConsoleCardDescription>
              Additional details about this Organization.
            </ConsoleCardDescription>
          </ConsoleCardDetails>
          <Button variant="outline" asChild>
            <Link to={`/organizations/${organizationId}/edit`}>Edit</Link>
          </Button>
        </ConsoleCardHeader>
        <ConsoleCardContent>
          <DetailsGrid>
            <DetailsGridColumn>
              {getProjectResponse?.project?.logInWithGoogle && (
                <DetailsGridEntry>
                  <DetailsGridKey>Log in with Google</DetailsGridKey>
                  <DetailsGridValue>
                    {getOrganizationResponse?.organization?.logInWithGoogle
                      ? 'Enabled'
                      : 'Disabled'}
                  </DetailsGridValue>
                </DetailsGridEntry>
              )}

              {getProjectResponse?.project?.logInWithMicrosoft && (
                <DetailsGridEntry>
                  <DetailsGridKey>Log in with Microsoft</DetailsGridKey>
                  <DetailsGridValue>
                    {getOrganizationResponse?.organization?.logInWithMicrosoft
                      ? 'Enabled'
                      : 'Disabled'}
                  </DetailsGridValue>
                </DetailsGridEntry>
              )}

              {getProjectResponse?.project?.logInWithGithub && (
                <DetailsGridEntry>
                  <DetailsGridKey>Log in with GitHub</DetailsGridKey>
                  <DetailsGridValue>
                    {getOrganizationResponse?.organization?.logInWithGithub
                      ? 'Enabled'
                      : 'Disabled'}
                  </DetailsGridValue>
                </DetailsGridEntry>
              )}

              {getProjectResponse?.project?.logInWithEmail && (
                <DetailsGridEntry>
                  <DetailsGridKey>Log in with Email</DetailsGridKey>
                  <DetailsGridValue>
                    {getOrganizationResponse?.organization?.logInWithEmail
                      ? 'Enabled'
                      : 'Disabled'}
                  </DetailsGridValue>
                </DetailsGridEntry>
              )}

              {getProjectResponse?.project?.logInWithPassword && (
                <DetailsGridEntry>
                  <DetailsGridKey>Log in with Password</DetailsGridKey>
                  <DetailsGridValue>
                    {getOrganizationResponse?.organization?.logInWithPassword
                      ? 'Enabled'
                      : 'Disabled'}
                  </DetailsGridValue>
                </DetailsGridEntry>
              )}
            </DetailsGridColumn>
            <DetailsGridColumn>
              {getProjectResponse?.project?.logInWithAuthenticatorApp && (
                <DetailsGridEntry>
                  <DetailsGridKey>Log in with Authenticator App</DetailsGridKey>
                  <DetailsGridValue>
                    {getOrganizationResponse?.organization
                      ?.logInWithAuthenticatorApp
                      ? 'Enabled'
                      : 'Disabled'}
                  </DetailsGridValue>
                </DetailsGridEntry>
              )}
              {getProjectResponse?.project?.logInWithPasskey && (
                <DetailsGridEntry>
                  <DetailsGridKey>Log in with Passkey</DetailsGridKey>
                  <DetailsGridValue>
                    {getOrganizationResponse?.organization?.logInWithPasskey
                      ? 'Enabled'
                      : 'Disabled'}
                  </DetailsGridValue>
                </DetailsGridEntry>
              )}
              <DetailsGridEntry>
                <DetailsGridKey>Require MFA</DetailsGridKey>
                <DetailsGridValue>
                  {getOrganizationResponse?.organization?.requireMfa
                    ? 'Enabled'
                    : 'Disabled'}
                </DetailsGridValue>
              </DetailsGridEntry>
            </DetailsGridColumn>
            <DetailsGridColumn>
              <DetailsGridEntry>
                <DetailsGridKey>Log in with SAML</DetailsGridKey>
                <DetailsGridValue>
                  {getOrganizationResponse?.organization?.logInWithSaml
                    ? 'Enabled'
                    : 'Disabled'}
                </DetailsGridValue>
              </DetailsGridEntry>
              <DetailsGridEntry>
                <DetailsGridKey>SCIM Directory Syncing</DetailsGridKey>
                <DetailsGridValue>
                  {getOrganizationResponse?.organization?.scimEnabled
                    ? 'Enabled'
                    : 'Disabled'}
                </DetailsGridValue>
              </DetailsGridEntry>
              <DetailsGridEntry>
                <DetailsGridKey>SAML / SCIM Domains</DetailsGridKey>
                <DetailsGridValue>
                  {getOrganizationDomainsResponse?.organizationDomains?.domains
                    ? getOrganizationDomainsResponse.organizationDomains.domains.map(
                        (s) => <div key={s}>{s}</div>,
                      )
                    : '-'}
                </DetailsGridValue>
              </DetailsGridEntry>
            </DetailsGridColumn>
          </DetailsGrid>
        </ConsoleCardContent>
      </ConsoleCard>

      {getOrganizationResponse?.organization?.logInWithGoogle && (
        <ConsoleCard>
          <ConsoleCardHeader>
            <ConsoleCardDetails>
              <ConsoleCardTitle>Google configuration</ConsoleCardTitle>
              <ConsoleCardDescription>
                Settings related to logging into this organization with Google.
              </ConsoleCardDescription>
            </ConsoleCardDetails>
            <EditOrganizationGoogleConfigurationButton />
          </ConsoleCardHeader>
          <ConsoleCardContent>
            <DetailsGrid>
              <DetailsGridColumn>
                <DetailsGridEntry>
                  <DetailsGridKey>Log in with Google</DetailsGridKey>
                  <DetailsGridValue>Enabled</DetailsGridValue>
                </DetailsGridEntry>
              </DetailsGridColumn>
              <DetailsGridColumn>
                <DetailsGridEntry>
                  <DetailsGridKey>Google Hosted Domains</DetailsGridKey>
                  <DetailsGridValue>
                    {getOrganizationGoogleHostedDomainsResponse
                      ?.organizationGoogleHostedDomains?.googleHostedDomains
                      ? getOrganizationGoogleHostedDomainsResponse.organizationGoogleHostedDomains.googleHostedDomains.map(
                          (s) => <div key={s}>{s}</div>,
                        )
                      : '-'}
                  </DetailsGridValue>
                </DetailsGridEntry>
              </DetailsGridColumn>
            </DetailsGrid>
          </ConsoleCardContent>
        </ConsoleCard>
      )}

      {getOrganizationResponse?.organization?.logInWithMicrosoft && (
        <ConsoleCard>
          <ConsoleCardHeader>
            <ConsoleCardDetails>
              <ConsoleCardTitle>Microsoft configuration</ConsoleCardTitle>
              <ConsoleCardDescription>
                Settings related to logging into this organization with
                Microsoft.
              </ConsoleCardDescription>
            </ConsoleCardDetails>
            <EditOrganizationMicrosoftConfigurationButton />
          </ConsoleCardHeader>
          <ConsoleCardContent>
            <DetailsGrid>
              <DetailsGridColumn>
                <DetailsGridEntry>
                  <DetailsGridKey>Log in with Microsoft</DetailsGridKey>
                  <DetailsGridValue>Enabled</DetailsGridValue>
                </DetailsGridEntry>
              </DetailsGridColumn>
              <DetailsGridColumn>
                <DetailsGridEntry>
                  <DetailsGridKey>Microsoft Tenant IDs</DetailsGridKey>
                  <DetailsGridValue>
                    {getOrganizationMicrosoftTenantIdsResponse
                      ?.organizationMicrosoftTenantIds?.microsoftTenantIds
                      ? getOrganizationMicrosoftTenantIdsResponse.organizationMicrosoftTenantIds.microsoftTenantIds.map(
                          (s) => <div key={s}>{s}</div>,
                        )
                      : '-'}
                  </DetailsGridValue>
                </DetailsGridEntry>
              </DetailsGridColumn>
            </DetailsGrid>
          </ConsoleCardContent>
        </ConsoleCard>
      )}
    </div>
  );
};
