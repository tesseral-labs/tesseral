import React, { FC } from 'react';
import EditProjectGoogleSettingsPage from './pages/project/edit/EditProjectGoogleSettingsPage';
import EditProjectMicrosoftSettingsPage from './pages/project/edit/EditProjectMicrosoftSettingsPage';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { createConnectTransport } from '@connectrpc/connect-web';
import { type Transport } from '@connectrpc/connect';
import { TransportProvider } from '@connectrpc/connect-query';
import { BrowserRouter, Navigate, Route, Routes } from 'react-router-dom';
import NotFoundPage from './pages/NotFound';

import { useAccessToken } from '@/lib/use-access-token';
import { Toaster } from '@/components/ui/sonner';
import { PageShell } from '@/components/page';
import { API_URL } from './config';

import { ListOrganizationsPage } from '@/pages/organizations/ListOrganizationsPage';
import { ViewOrganizationPage } from '@/pages/organizations/ViewOrganizationPage';
import { ViewUserPage } from '@/pages/users/ViewUserPage';
import { ListAPIKeysPage } from '@/pages/api-keys/ListAPIKeysPage';
import { ViewProjectSettingsPage } from '@/pages/project/ViewProjectSettingsPage';
import { OrganizationUsersTab } from '@/pages/organizations/OrganizationUsersTab';
import { OrganizationSAMLConnectionsTab } from '@/pages/organizations/OrganizationSAMLConnectionsTab';
import { OrganizationSCIMAPIKeysTab } from '@/pages/organizations/OrganizationSCIMAPIKeysTab';
import { OrganizationDetailsTab } from '@/pages/organizations/OrganizationDetailsTab';
import { EditOrganizationPage } from '@/pages/organizations/EditOrganizationPage';
import { ViewSAMLConnectionPage } from '@/pages/saml-connections/ViewSAMLConnectionPage';
import { EditSAMLConnectionPage } from '@/pages/saml-connections/EditSAMLConnectionPage';
import { ViewSCIMAPIKeyPage } from '@/pages/scim-api-keys/ViewSCIMAPIKeyPage';
import { HomePage } from '@/pages/home/HomePage';
import { ProjectDetailsTab } from '@/pages/project/ProjectDetailsTab';
import { ViewPasskeyPage } from '@/pages/passkeys/ViewPasskeyPage';
import { OrganizationUserInvitesTab } from '@/pages/organizations/OrganizationUserInvitesTab';
import { ViewUserInvitePage } from '@/pages/user-invites/ViewUserInvitePage';
import ProjectUISettingsPage from './pages/project/project-ui-settings/ProjectUISettings';
import { ViewPublishableKeyPage } from '@/pages/api-keys/ViewPublishableKeyPage';
import { VaultDomainSettingsTab } from '@/pages/project/VaultDomainSettingsTab';
import { ViewBackendAPIKeyPage } from '@/pages/api-keys/ViewBackendAPIKeyPage';

import {
  AuthenticateAnotherWayPage,
  FinishLoginPage,
  GoogleOAuthCallbackPage,
  ImpersonatePage,
  LoginFlowLayout,
  LoginPage,
  LogoutPage,
  MicrosoftOAuthCallbackPage,
  OrganizationLoginPage,
  RegisterAuthenticatorAppPage,
  RegisterPasskeyPage,
  RegisterPasswordPage,
  RegisterSecondaryFactorPage,
  SwitchOrganizationsPage,
  VerifyAuthenticatorAppPage,
  VerifyAuthenticatorAppRecoveryCodePage,
  VerifyEmailPage,
  VerifyPasskeyPage,
  VerifyPasswordPage,
  VerifySecondaryFactorPage,
} from '@tesseral/common-ui';
import ChooseProjectPage from '@/pages/login/ChooseProjectPage';
import CreateProjectPage from '@/pages/login/CreateProjectPage';
import { LoginPageWrapper } from './components/login-page';

const queryClient = new QueryClient();

const useTransport = (): Transport => {
  const accessToken = useAccessToken();

  return createConnectTransport({
    baseUrl: `${API_URL}/api/internal/connect`,
    fetch: (input, init) =>
      fetch(input, {
        ...init,
        headers: {
          Authorization: `Bearer ${accessToken}`,
          'Content-Type': 'application/json',
        },
        credentials: 'include',
      }),
    interceptors: [
      (next) => async (req) => {
        return next(req);
      },
    ],
  });
};

const AppWithinQueryClient = () => {
  const transport = useTransport();
  return (
    <TransportProvider transport={transport}>
      <BrowserRouter>
        <Routes>
          <Route path="" element={<LoginFlowLayout />}>
            <Route path="login" element={<LoginPage />} />
            <Route
              path="choose-organization"
              element={<Navigate to="/choose-project" />}
            />
            <Route path="verify-email" element={<VerifyEmailPage />} />
            <Route
              path="google-oauth-callback"
              element={<GoogleOAuthCallbackPage />}
            />
            <Route
              path="microsoft-oauth-callback"
              element={<MicrosoftOAuthCallbackPage />}
            />
            <Route path="choose-project" element={<ChooseProjectPage />} />
            <Route path="create-project" element={<CreateProjectPage />} />
            <Route
              path="organizations/:organizationId/login"
              element={<OrganizationLoginPage />}
            />
            <Route
              path="authenticate-another-way"
              element={<AuthenticateAnotherWayPage />}
            />
            <Route path="verify-password" element={<VerifyPasswordPage />} />
            <Route
              path="verify-secondary-factor"
              element={<VerifySecondaryFactorPage />}
            />
            <Route
              path="verify-authenticator-app"
              element={<VerifyAuthenticatorAppPage />}
            />
            <Route
              path="verify-authenticator-app-recovery-code"
              element={<VerifyAuthenticatorAppRecoveryCodePage />}
            />
            <Route path="verify-passkey" element={<VerifyPasskeyPage />} />
            <Route
              path="register-password"
              element={<RegisterPasswordPage />}
            />
            <Route
              path="register-secondary-factor"
              element={<RegisterSecondaryFactorPage />}
            />
            <Route path="register-passkey" element={<RegisterPasskeyPage />} />
            <Route
              path="register-authenticator-app"
              element={<RegisterAuthenticatorAppPage />}
            />
            <Route path="finish-login" element={<FinishLoginPage />} />

            <Route path="impersonate" element={<ImpersonatePage />} />
            <Route
              path="switch-organizations/:organizationId"
              element={<SwitchOrganizationsPage />}
            />

            <Route path="logout" element={<LogoutPage />} />
          </Route>

          <Route path="/" element={<PageShell />}>
            <Route path="" element={<HomePage />} />
            <Route
              path="project-settings"
              element={<ViewProjectSettingsPage />}
            >
              <Route path="" element={<ProjectDetailsTab />} />
              <Route
                path="vault-domain-settings"
                element={<VaultDomainSettingsTab />}
              />
            </Route>
            <Route
              path="project-settings/log-in-with-google/edit"
              element={<EditProjectGoogleSettingsPage />}
            />
            <Route
              path="project-settings/log-in-with-microsoft/edit"
              element={<EditProjectMicrosoftSettingsPage />}
            />

            <Route
              path="project-settings/api-keys"
              element={<ListAPIKeysPage />}
            />

            <Route
              path="project-settings/api-keys/publishable-keys/:publishableKeyId"
              element={<ViewPublishableKeyPage />}
            />

            <Route
              path="project-settings/api-keys/backend-api-keys/:backendApiKeyId"
              element={<ViewBackendAPIKeyPage />}
            />

            <Route
              path="project-settings/vault-ui-settings"
              element={<ProjectUISettingsPage />}
            />

            <Route path="organizations" element={<ListOrganizationsPage />} />

            <Route
              path="organizations/:organizationId"
              element={<ViewOrganizationPage />}
            >
              <Route path="" element={<OrganizationDetailsTab />} />
              <Route path="users" element={<OrganizationUsersTab />} />
              <Route
                path="user-invites"
                element={<OrganizationUserInvitesTab />}
              />
              <Route
                path="saml-connections"
                element={<OrganizationSAMLConnectionsTab />}
              />
              <Route
                path="scim-api-keys"
                element={<OrganizationSCIMAPIKeysTab />}
              />
            </Route>

            <Route
              path="organizations/:organizationId/edit"
              element={<EditOrganizationPage />}
            />
            <Route
              path="organizations/:organizationId/saml-connections/:samlConnectionId"
              element={<ViewSAMLConnectionPage />}
            />
            <Route
              path="organizations/:organizationId/saml-connections/:samlConnectionId/edit"
              element={<EditSAMLConnectionPage />}
            />
            <Route
              path="organizations/:organizationId/users/:userId"
              element={<ViewUserPage />}
            />
            <Route
              path="organizations/:organizationId/users/:userId/passkeys/:passkeyId"
              element={<ViewPasskeyPage />}
            />
            <Route
              path="organizations/:organizationId/user-invites/:userInviteId"
              element={<ViewUserInvitePage />}
            />
            <Route
              path="organizations/:organizationId/scim-api-keys/:scimApiKeyId"
              element={<ViewSCIMAPIKeyPage />}
            />
          </Route>

          <Route path="*" element={<NotFoundPage />} />
        </Routes>
      </BrowserRouter>
    </TransportProvider>
  );
};

const App: FC = () => {
  return (
    <QueryClientProvider client={queryClient}>
      <AppWithinQueryClient />
      <Toaster />
    </QueryClientProvider>
  );
};

export default App;
