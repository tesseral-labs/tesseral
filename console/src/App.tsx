import { TransportProvider } from "@connectrpc/connect-query";
import { createConnectTransport } from "@connectrpc/connect-web";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import React from "react";
import { BrowserRouter, Navigate, Route, Routes } from "react-router";
import { Toaster } from "sonner";

import { PageShell } from "@/components/page";
import { NotFoundPage } from "@/pages/NotFoundPage";
import { HomePage } from "@/pages/console/home/HomePage";
import { ListOrganizationsPage } from "@/pages/console/organizations/ListOrganizationsPage";
import { OrganizationDetailsTab } from "@/pages/console/organizations/OrganizationDetailsTab";
import { OrganizationPage } from "@/pages/console/organizations/OrganizationPage";
import { OrganizationUsersTab } from "@/pages/console/organizations/OrganizationUsersTab";
import { UserDetailsTab } from "@/pages/console/organizations/users/UserDetailsTab";
import { UserPage } from "@/pages/console/organizations/users/UserPage";

import { API_URL } from "./config";
import { OrganizationApiKeysTab } from "./pages/console/organizations/OrganizationApiKeysTab";
import { OrganizationAuthentication } from "./pages/console/organizations/OrganizationAuthenticationTab";
import { OrganizationLogs } from "./pages/console/organizations/OrganizationLogsTab";
import { OrganizationApiKeyRolesTab } from "./pages/console/organizations/api-keys/OrganizationAPIKeyRolesTab";
import { OrganizationApiKeyDetailsTab } from "./pages/console/organizations/api-keys/OrganizationApiKeyDetailsTab";
import { OrganizationApiKeyLogsTab } from "./pages/console/organizations/api-keys/OrganizationApiKeyLogsTab";
import { OrganizationApiKeyPage } from "./pages/console/organizations/api-keys/OrganizationApiKeyPage";
import { OrganizationSamlConnectionPage } from "./pages/console/organizations/saml-connections/OrganizationSamlConnectionPage";
import { OrganizationScimApiKeyPage } from "./pages/console/organizations/scim-api-keys/OrganizationScimApiKeyPage";
import { UserActivityTab } from "./pages/console/organizations/users/UserActivityTab";
import { UserHistoryTab } from "./pages/console/organizations/users/UserHistoryTab";
import { UserPasskeysTab } from "./pages/console/organizations/users/UserPasskeysTab";
import { UserRolesTab } from "./pages/console/organizations/users/UserRolesTab";
import { UserSessionsTab } from "./pages/console/organizations/users/UserSessionsTab";
import { PasskeyPage } from "./pages/console/organizations/users/passkeys/PasskeyPage";
import { SessionPage } from "./pages/console/organizations/users/sessions/SessionPage";
import { AccessSettingsPage } from "./pages/console/settings/AccessSettingsPage";
import { ApiKeySettingsPage } from "./pages/console/settings/ApiKeySettingsPage";
import { SettingsOverviewPage } from "./pages/console/settings/SettingsOverviewPage";
import { VaultCustomizationPage } from "./pages/console/settings/VaultCustomizationPage";
import { BackendApiKeysPage } from "./pages/console/settings/api-keys/BackendApiKeysPage";
import { BackendApiKeyPage } from "./pages/console/settings/api-keys/backend-api-keys/BackendApiKeyPage";
import { AuthenticationSettingsPage } from "./pages/console/settings/authentication/AuthenticationSettingsPage";
import { VaultBrandingSettingsTab } from "./pages/console/settings/vault/VaultBrandingSettingsTab";
import { VaultDetailsTab } from "./pages/console/settings/vault/VaultDetailsTab";
import { VaultDomainSettingsTab } from "./pages/console/settings/vault/VaultDomainSettingsTab";
import { StripeCheckoutSuccessPage } from "./pages/console/stripe/StripeCheckoutSuccessPage";
import { AuthenticateAnotherWayPage } from "./pages/login/AuthenticateAnotherWayPage";
import { ChooseProjectPage } from "./pages/login/ChooseProjectPage";
import { CreateProjectPage } from "./pages/login/CreateProjectPage";
import { FinishLoginPage } from "./pages/login/FinishLoginPage";
import { ForgotPasswordPage } from "./pages/login/ForgotPasswordPage";
import { GithubOAuthCallbackPage } from "./pages/login/GithubOAuthCallbackPage";
import { GoogleOAuthCallbackPage } from "./pages/login/GoogleOAuthCallbackPage";
import { ImpersonatePage } from "./pages/login/ImpersonatePage";
import { LoginFlowLayout } from "./pages/login/LoginFlowLayout";
import { LoginPage } from "./pages/login/LoginPage";
import { LogoutPage } from "./pages/login/LogoutPage";
import { MicrosoftOAuthCallbackPage } from "./pages/login/MicrosoftOAuthCallbackPage";
import { OrganizationLoginPage } from "./pages/login/OrganizationLoginPage";
import { RegisterAuthenticatorAppPage } from "./pages/login/RegisterAuthenticatorAppPage";
import { RegisterPasskeyPage } from "./pages/login/RegisterPasskeyPage";
import { RegisterPasswordPage } from "./pages/login/RegisterPasswordPage";
import { RegisterSecondaryFactorPage } from "./pages/login/RegisterSecondaryFactorPage";
import { SignupPage } from "./pages/login/SignupPage";
import { SwitchOrganizationsPage } from "./pages/login/SwitchOrganizationsPage";
import { VerifyAuthenticatorAppPage } from "./pages/login/VerifyAuthenticatorAppPage";
import { VerifyAuthenticatorAppRecoveryCodePage } from "./pages/login/VerifyAuthenticatorAppRecoveryCodePage";
import { VerifyEmailPage } from "./pages/login/VerifyEmailPage";
import { VerifyPasskeyPage } from "./pages/login/VerifyPasskeyPage";
import { VerifyPasswordPage } from "./pages/login/VerifyPasswordPage";
import { VerifySecondaryFactorPage } from "./pages/login/VerifySecondaryFactorPage";

const queryClient = new QueryClient();

const transport = createConnectTransport({
  baseUrl: `${API_URL}/api/internal/connect`,
  fetch: (input, init) =>
    fetch(input, {
      ...init,
      credentials: "include",
    }),
});

function AppWithinQueryClient() {
  return (
    <TransportProvider transport={transport}>
      <BrowserRouter>
        <Routes>
          <Route
            path="/project-settings/publishable-keys"
            element={<Navigate to="/settings/api-keys" replace />}
          />

          {/* Console Routes */}
          <Route path="/" element={<PageShell />}>
            <Route path="" element={<HomePage />} />

            <Route path="organizations">
              <Route path="" element={<ListOrganizationsPage />} />
              <Route path=":organizationId" element={<OrganizationPage />}>
                <Route path="" element={<OrganizationDetailsTab />} />
                <Route
                  path="authentication"
                  element={<OrganizationAuthentication />}
                />
                <Route path="api-keys" element={<OrganizationApiKeysTab />} />
                <Route path="users" element={<OrganizationUsersTab />} />

                <Route path="logs" element={<OrganizationLogs />} />
              </Route>
            </Route>

            <Route
              path="organizations/:organizationId/users/:userId"
              element={<UserPage />}
            >
              <Route path="" element={<UserDetailsTab />} />
              <Route path="sessions" element={<UserSessionsTab />} />
              <Route path="roles" element={<UserRolesTab />} />
              <Route path="passkeys" element={<UserPasskeysTab />} />
              <Route path="history" element={<UserHistoryTab />} />
              <Route path="activity" element={<UserActivityTab />} />
            </Route>

            <Route
              path="organizations/:organizationId/users/:userId/passkeys/:passkeyId"
              element={<PasskeyPage />}
            />

            <Route
              path="organizations/:organizationId/users/:userId/sessions/:sessionId"
              element={<SessionPage />}
            />

            <Route
              path="organizations/:organizationId/api-keys/:apiKeyId"
              element={<OrganizationApiKeyPage />}
            >
              <Route path="" element={<OrganizationApiKeyDetailsTab />} />
              <Route path="roles" element={<OrganizationApiKeyRolesTab />} />
              <Route path="logs" element={<OrganizationApiKeyLogsTab />} />
            </Route>

            <Route
              path="organizations/:organizationId/saml-connections/:samlConnectionId"
              element={<OrganizationSamlConnectionPage />}
            />

            <Route
              path="organizations/:organizationId/scim-api-keys/:scimApiKeyId"
              element={<OrganizationScimApiKeyPage />}
            />

            <Route path="settings">
              <Route path="" element={<SettingsOverviewPage />} />
              <Route
                path="authentication"
                element={<AuthenticationSettingsPage />}
              />
              <Route path="api-keys" element={<ApiKeySettingsPage />} />
              <Route path="api-keys/backend-api-keys">
                <Route path="" element={<BackendApiKeysPage />} />
                <Route
                  path=":backendApiKeyId"
                  element={<BackendApiKeyPage />}
                />
              </Route>

              <Route path="access" element={<AccessSettingsPage />} />
              <Route path="vault" element={<VaultCustomizationPage />}>
                <Route path="" element={<VaultDetailsTab />} />
                <Route path="domains" element={<VaultDomainSettingsTab />} />
                <Route path="branding" element={<VaultBrandingSettingsTab />} />
              </Route>
            </Route>

            <Route
              path="stripe-checkout-success"
              element={<StripeCheckoutSuccessPage />}
            />
          </Route>

          {/* Login and Signup Routes */}
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
            <Route path="choose-organization" element={<ChooseProjectPage />} />
            <Route path="create-organization" element={<CreateProjectPage />} />
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

          <Route path="*" element={<NotFoundPage />} />
        </Routes>
      </BrowserRouter>
    </TransportProvider>
  );
}

export function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <AppWithinQueryClient />
      <Toaster />
    </QueryClientProvider>
  );
}
