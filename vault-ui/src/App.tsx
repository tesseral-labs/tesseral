import { Transport } from "@connectrpc/connect";
import { TransportProvider } from "@connectrpc/connect-query";
import { createConnectTransport } from "@connectrpc/connect-web";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import React from "react";
import { Navigate, Route, Routes } from "react-router";
import { BrowserRouter } from "react-router";

import { LoginFlowLayout } from "@/components/login/LoginFlowLayout";
import { Toaster } from "@/components/ui/sonner";
import { NotFoundPage } from "@/pages/NotFoundPage";
import { AuthenticateAnotherWayPage } from "@/pages/login/AuthenticateAnotherWayPage";
import { ChooseOrganizationPage } from "@/pages/login/ChooseOrganizationPage";
import { CreateOrganizationPage } from "@/pages/login/CreateOrganizationPage";
import { FinishLoginPage } from "@/pages/login/FinishLoginPage";
import { ForgotPasswordPage } from "@/pages/login/ForgotPasswordPage";
import { GoogleOAuthCallbackPage } from "@/pages/login/GoogleOAuthCallbackPage";
import { ImpersonatePage } from "@/pages/login/ImpersonatePage";
import { LoginPage } from "@/pages/login/LoginPage";
import { LogoutPage } from "@/pages/login/LogoutPage";
import { MicrosoftOAuthCallbackPage } from "@/pages/login/MicrosoftOAuthCallbackPage";
import { OrganizationLoginPage } from "@/pages/login/OrganizationLoginPage";
import { RegisterAuthenticatorAppPage } from "@/pages/login/RegisterAuthenticatorAppPage";
import { RegisterPasskeyPage } from "@/pages/login/RegisterPasskeyPage";
import { RegisterPasswordPage } from "@/pages/login/RegisterPasswordPage";
import { RegisterSecondaryFactorPage } from "@/pages/login/RegisterSecondaryFactorPage";
import { SignupPage } from "@/pages/login/SignupPage";
import { SwitchOrganizationsPage } from "@/pages/login/SwitchOrganizationsPage";
import { VerifyAuthenticatorAppPage } from "@/pages/login/VerifyAuthenticatorAppPage";
import { VerifyAuthenticatorAppRecoveryCodePage } from "@/pages/login/VerifyAuthenticatorAppRecoveryCodePage";
import { VerifyEmailPage } from "@/pages/login/VerifyEmailPage";
import { VerifyPasskeyPage } from "@/pages/login/VerifyPasskeyPage";
import { VerifyPasswordPage } from "@/pages/login/VerifyPasswordPage";
import { VerifySecondaryFactorPage } from "@/pages/login/VerifySecondaryFactorPage";

import { Page } from "./components/page";
import { LoggedInGate } from "./components/page/LoggedInGate";
import { GithubOAuthCallbackPage } from "./pages/login/GithubOAuthCallbackPage";
import { AuditLogsPage } from "./pages/vault/AuditLogsPage";
import { OrganizationPage } from "./pages/vault/OrganizationPage";
import { UserPage } from "./pages/vault/UserPage";
import { OrganizationApiKeysTab } from "./pages/vault/organization/OrganizationApiKeysTab";
import { OrganizationAuthenticationTab } from "./pages/vault/organization/OrganizationAuthenticationTab";
import { OrganizationDetailsTab } from "./pages/vault/organization/OrganizationDetailsTab";
import { OrganizationUserInvitesTab } from "./pages/vault/organization/OrganizationUserInvitesTab";
import { OrganizationUsersTab } from "./pages/vault/organization/OrganizationUsersTab";
import { ApiKeyDetailsTab } from "./pages/vault/organization/api-keys/ApiKeyDetailsTab";
import { ApiKeyPage } from "./pages/vault/organization/api-keys/ApiKeyPage";
import { ApiKeyRolesTab } from "./pages/vault/organization/api-keys/ApiKeyRolesTab";
import { SamlConnectionPage } from "./pages/vault/organization/saml-connections/SamlConnectionPage";
import { SamlConnectionsPage } from "./pages/vault/organization/saml-connections/SamlConnectionsPage";
import { ScimApiKeyPage } from "./pages/vault/organization/scim-api-keys/ScimApiKeyPage";
import { OrganizationUserDetailsTab } from "./pages/vault/organization/users/OrganizationUserDetailsTab";
import { OrganizationUserPage } from "./pages/vault/organization/users/OrganizationUserPage";
import { OrganizationUserRolesTab } from "./pages/vault/organization/users/OrganizationUserRolesTab";
import { UserAuthenticationTab } from "./pages/vault/user/UserAuthenticationTab";
import { UserDetailsTab } from "./pages/vault/user/UserDetailsTab";

const queryClient = new QueryClient();

function useTransport(): Transport {
  return createConnectTransport({
    baseUrl: `/api/internal/connect`,
    fetch: (input, init) => fetch(input, { ...init, credentials: "include" }),
  });
}

function AppWithRoutes() {
  const transport = useTransport();

  return (
    <TransportProvider transport={transport}>
      <QueryClientProvider client={queryClient}>
        <BrowserRouter>
          <Routes>
            <Route
              path="/audit-logs"
              element={<Navigate to="/logs" replace />}
            />
            <Route
              path="/user-settings"
              element={<Navigate to="/user" replace />}
            />
            <Route
              path="/organization-settings"
              element={<Navigate to="/organization" replace />}
            />
            <Route path="/" element={<Navigate to="/user" replace />} />

            <Route path="/" element={<LoggedInGate />}>
              <Route path="" element={<Page />}>
                <Route path="logs" element={<AuditLogsPage />} />
                <Route path="organization" element={<OrganizationPage />}>
                  <Route index element={<OrganizationDetailsTab />} />
                  <Route path="users" element={<OrganizationUsersTab />} />
                  <Route
                    path="user-invites"
                    element={<OrganizationUserInvitesTab />}
                  />
                  <Route
                    path="authentication"
                    element={<OrganizationAuthenticationTab />}
                  />
                  <Route path="api-keys" element={<OrganizationApiKeysTab />} />
                </Route>
                <Route
                  path="organization/api-keys/:apiKeyId"
                  element={<ApiKeyPage />}
                >
                  <Route index element={<ApiKeyDetailsTab />} />
                  <Route path="roles" element={<ApiKeyRolesTab />} />
                </Route>
                <Route
                  path="organization/saml-connections"
                  element={<SamlConnectionsPage />}
                />
                <Route
                  path="organization/saml-connections/:samlConnectionId"
                  element={<SamlConnectionPage />}
                />
                <Route
                  path="organization/scim-api-keys/:scimApiKeyId"
                  element={<ScimApiKeyPage />}
                />
                <Route
                  path="organization/users/:userId"
                  element={<OrganizationUserPage />}
                >
                  <Route index element={<OrganizationUserDetailsTab />} />
                  <Route path="roles" element={<OrganizationUserRolesTab />} />
                </Route>

                <Route path="user" element={<UserPage />}>
                  <Route index element={<UserDetailsTab />} />
                  <Route
                    path="authentication"
                    element={<UserAuthenticationTab />}
                  />
                </Route>
              </Route>
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
              <Route
                path="register-passkey"
                element={<RegisterPasskeyPage />}
              />
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

            <Route path="" element={<LoggedInGate />}>
              <Route path="" element={<Page />} />
            </Route>

            <Route path="*" element={<NotFoundPage />} />
          </Routes>
        </BrowserRouter>

        <Toaster />
      </QueryClientProvider>
    </TransportProvider>
  );
}

export function App() {
  return <AppWithRoutes />;
}
