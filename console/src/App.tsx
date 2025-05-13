import React, { FC, useMemo } from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { createConnectTransport } from '@connectrpc/connect-web';
import { type Transport } from '@connectrpc/connect';
import { TransportProvider } from '@connectrpc/connect-query';
import { BrowserRouter, Route, Routes } from 'react-router-dom';
import NotFoundPage from './pages/NotFound';
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
import { Toaster } from '@/components/ui/sonner';
import { EditSAMLConnectionPage } from '@/pages/saml-connections/EditSAMLConnectionPage';
import { PageShell } from '@/components/page';
import { ViewSCIMAPIKeyPage } from '@/pages/scim-api-keys/ViewSCIMAPIKeyPage';
import { HomePage } from '@/pages/home/HomePage';
import { ProjectDetailsTab } from '@/pages/project/ProjectDetailsTab';
import { ViewPasskeyPage } from '@/pages/passkeys/ViewPasskeyPage';
import { OrganizationUserInvitesTab } from '@/pages/organizations/OrganizationUserInvitesTab';
import { ViewUserInvitePage } from '@/pages/user-invites/ViewUserInvitePage';
import { API_URL } from './config';
import { ViewPublishableKeyPage } from '@/pages/api-keys/ViewPublishableKeyPage';
import ProjectUISettingsPage from './pages/project/project-ui-settings/ProjectUISettings';
import { VaultDomainSettingsTab } from '@/pages/project/VaultDomainSettingsTab';
import { ViewBackendAPIKeyPage } from '@/pages/api-keys/ViewBackendAPIKeyPage';
import { LoginPage } from '@/pages/login/LoginPage';
import { SignupPage } from '@/pages/login/SignupPage';
import { LoginFlowLayout } from '@/pages/login/LoginFlowLayout';
import { VerifyEmailPage } from '@/pages/login/VerifyEmailPage';
import { GoogleOAuthCallbackPage } from '@/pages/login/GoogleOAuthCallbackPage';
import { MicrosoftOAuthCallbackPage } from '@/pages/login/MicrosoftOAuthCallbackPage';
import { ChooseOrganizationPage } from '@/pages/login/ChooseOrganizationPage';
import { CreateOrganizationPage } from '@/pages/login/CreateOrganizationPage';
import { OrganizationLoginPage } from '@/pages/login/OrganizationLoginPage';
import { AuthenticateAnotherWayPage } from '@/pages/login/AuthenticateAnotherWayPage';
import { VerifyPasswordPage } from '@/pages/login/VerifyPasswordPage';
import { ForgotPasswordPage } from '@/pages/login/ForgotPasswordPage';
import { VerifySecondaryFactorPage } from '@/pages/login/VerifySecondaryFactorPage';
import { VerifyAuthenticatorAppPage } from '@/pages/login/VerifyAuthenticatorAppPage';
import { VerifyAuthenticatorAppRecoveryCodePage } from '@/pages/login/VerifyAuthenticatorAppRecoveryCodePage';
import { VerifyPasskeyPage } from '@/pages/login/VerifyPasskeyPage';
import { RegisterPasswordPage } from '@/pages/login/RegisterPasswordPage';
import { RegisterSecondaryFactorPage } from '@/pages/login/RegisterSecondaryFactorPage';
import { RegisterPasskeyPage } from '@/pages/login/RegisterPasskeyPage';
import { RegisterAuthenticatorAppPage } from '@/pages/login/RegisterAuthenticatorAppPage';
import { FinishLoginPage } from '@/pages/login/FinishLoginPage';
import { ImpersonatePage } from '@/pages/login/ImpersonatePage';
import { SwitchOrganizationsPage } from '@/pages/login/SwitchOrganizationsPage';
import { LogoutPage } from '@/pages/login/LogoutPage';
import { useAccessToken } from '@/lib/AccessTokenProvider';
import { StripeCheckoutSuccessPage } from './pages/stripe/StripeCheckoutSuccessPage';
import { RBACSettingsTab } from '@/pages/project/RBACSettingsTab';
import { EditRBACPolicyPage } from '@/pages/project/EditRBACPolicyPage';
import { ViewRolePage } from '@/pages/roles/ViewRolePage';
import { OrganizationRolesTab } from '@/pages/organizations/OrganizationRolesTab';
import { CreateRolePage } from '@/pages/roles/CreateRolePage';
import { EditRolePage } from '@/pages/roles/EditRolePage';
import { GithubOAuthCallbackPage } from './pages/login/GithubOAuthCallbackPage';

const queryClient = new QueryClient();

const transport = createConnectTransport({
  baseUrl: `${API_URL}/api/internal/connect`,
  fetch: (input, init) =>
    fetch(input, {
      ...init,
      credentials: 'include',
    }),
});

const AppWithinQueryClient = () => {
  return (
    <TransportProvider transport={transport}>
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<PageShell />}>
            <Route path="" element={<HomePage />} />
          </Route>

          <Route path="login" element={<LoginPage />} />
          <Route path="signup" element={<SignupPage />} />

          <Route path="" element={<LoginFlowLayout />}>
            <Route path="verify-email" element={<VerifyEmailPage />} />
            <Route
              path="github-oauth-callback"
              element={<GithubOAuthCallbackPage />}
            />
            <Route
              path="google-oauth-callback"
              element={<GoogleOAuthCallbackPage />}
            />
            <Route
              path="microsoft-oauth-callback"
              element={<MicrosoftOAuthCallbackPage />}
            />
            <Route
              path="choose-organization"
              element={<ChooseOrganizationPage />}
            />
            <Route
              path="create-organization"
              element={<CreateOrganizationPage />}
            />
            <Route
              path="organizations/:organizationId/login"
              element={<OrganizationLoginPage />}
            />
            <Route
              path="authenticate-another-way"
              element={<AuthenticateAnotherWayPage />}
            />
            <Route path="verify-password" element={<VerifyPasswordPage />} />
            <Route path="forgot-password" element={<ForgotPasswordPage />} />
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

          <Route path="" element={<PageShell />}>
            <Route
              path="project-settings"
              element={<ViewProjectSettingsPage />}
            >
              <Route path="" element={<ProjectDetailsTab />} />
              <Route
                path="vault-ui-settings"
                element={<ProjectUISettingsPage />}
              />
              <Route
                path="vault-domain-settings"
                element={<VaultDomainSettingsTab />}
              />
              <Route path="rbac-settings" element={<RBACSettingsTab />} />
            </Route>

            <Route
              path="project-settings/rbac-settings/rbac-policy/edit"
              element={<EditRBACPolicyPage />}
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
              <Route path="roles" element={<OrganizationRolesTab />} />
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

            <Route path="roles/new" element={<CreateRolePage />} />
            <Route path="roles/:roleId" element={<ViewRolePage />} />
            <Route path="roles/:roleId/edit" element={<EditRolePage />} />

            <Route
              path="stripe-checkout-success"
              element={<StripeCheckoutSuccessPage />}
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
