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
import { OwnerGate } from "./components/page/OwnerGate";
import { SetUpSamlConnectionDialog } from "./components/saml-connections/SetUpSamlConnectionDialog";
import { TestSamlConnectionDialog } from "./components/saml-connections/TestSamlConnectionDialog";
import { AssignEntraSamlUsers } from "./components/saml-connections/entra/AssignEntraSamlUsers";
import { ConfigureEntraSamlIdentifier } from "./components/saml-connections/entra/ConfigureEntraSamlIdentifier";
import { ConfigureEntraSamlReplyUrl } from "./components/saml-connections/entra/ConfigureEntraSamlReplyUrl";
import { CreateEntraSamlApplication } from "./components/saml-connections/entra/CreateEntraSamlApplication";
import { DownloadEntraSamlMetadata } from "./components/saml-connections/entra/DownloadEntraSamlMetadata";
import { EntraSamlConnectionFlow } from "./components/saml-connections/entra/EntraSamlConnectionFlow";
import { AssignGoogleSamlUsers } from "./components/saml-connections/google/AssignGoogleSamlUsers";
import { ConfigureGoogleSamlApplication } from "./components/saml-connections/google/ConfigureGoogleSamlApplication";
import { CreateGoogleSamlApplication } from "./components/saml-connections/google/CreateGoogleSamlApplication";
import { DownloadGoogleSamlMetadata } from "./components/saml-connections/google/DownloadGoogleSamlMetadata";
import { GoogleSamlConnectionFlow } from "./components/saml-connections/google/GoogleSamlConnectionFlow";
import { NameGoogleSamlApplication } from "./components/saml-connections/google/NameGoogleSamlApplication";
import { AssignOktaSamlUsers } from "./components/saml-connections/okta/AssignOktaSamlUsers";
import { ConfigureOktaSamlApplication } from "./components/saml-connections/okta/ConfigureOktaSamlApplication";
import { CreateOktaSamlApplication } from "./components/saml-connections/okta/CreateOktaSamlApplication";
import { NameOktaSamlApplication } from "./components/saml-connections/okta/NameOktaSamlApplication";
import { OktaSamlConnectionFlow } from "./components/saml-connections/okta/OktaSamlConnectionFlow";
import { SyncOktaSamlMetadata } from "./components/saml-connections/okta/SyncOktaSamlMetadata";
import { AssignOtherSamlUsers } from "./components/saml-connections/other/AssignOtherSamlUsers";
import { ConfigureOtherSamlApplication } from "./components/saml-connections/other/ConfigureOtherSamlApplication";
import { CreateOtherSamlApplication } from "./components/saml-connections/other/CreateOtherSamlApplication";
import { DownloadOtherSamlMetadata } from "./components/saml-connections/other/DownloadOtherSamlMetadata";
import { OtherSamlConnectionFlow } from "./components/saml-connections/other/OtherSamlConnectionFlow";
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
import { OidcConnectionPage } from "./pages/vault/organization/oidc-connections/OidcConnectionPage";
import { OidcConnectionsPage } from "./pages/vault/organization/oidc-connections/OidcConnectionsPage";
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
                <Route path="logs" element={<OwnerGate />}>
                  <Route index element={<AuditLogsPage />} />
                </Route>
                <Route path="organization" element={<OwnerGate />}>
                  <Route path="" element={<OrganizationPage />}>
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
                    <Route
                      path="api-keys"
                      element={<OrganizationApiKeysTab />}
                    />
                  </Route>
                  <Route path="api-keys/:apiKeyId" element={<ApiKeyPage />}>
                    <Route index element={<ApiKeyDetailsTab />} />
                    <Route path="roles" element={<ApiKeyRolesTab />} />
                  </Route>
                  <Route
                    path="saml-connections"
                    element={<SamlConnectionsPage />}
                  />
                  <Route
                    path="saml-connections/:samlConnectionId"
                    element={<SamlConnectionPage />}
                  >
                    <Route
                      path="setup"
                      element={<SetUpSamlConnectionDialog />}
                    />
                    <Route
                      path="setup/google"
                      element={<GoogleSamlConnectionFlow />}
                    >
                      <Route index element={<CreateGoogleSamlApplication />} />
                      <Route
                        path="name"
                        element={<NameGoogleSamlApplication />}
                      />
                      <Route
                        path="metadata"
                        element={<DownloadGoogleSamlMetadata />}
                      />
                      <Route
                        path="configure"
                        element={<ConfigureGoogleSamlApplication />}
                      />
                      <Route path="users" element={<AssignGoogleSamlUsers />} />
                    </Route>
                    <Route
                      path="setup/entra"
                      element={<EntraSamlConnectionFlow />}
                    >
                      <Route index element={<CreateEntraSamlApplication />} />
                      <Route
                        path="identifier"
                        element={<ConfigureEntraSamlIdentifier />}
                      />
                      <Route
                        path="reply-url"
                        element={<ConfigureEntraSamlReplyUrl />}
                      />
                      <Route
                        path="metadata"
                        element={<DownloadEntraSamlMetadata />}
                      />
                      <Route path="users" element={<AssignEntraSamlUsers />} />
                    </Route>
                    <Route
                      path="setup/okta"
                      element={<OktaSamlConnectionFlow />}
                    >
                      <Route index element={<CreateOktaSamlApplication />} />
                      <Route
                        path="name"
                        element={<NameOktaSamlApplication />}
                      />
                      <Route
                        path="configure"
                        element={<ConfigureOktaSamlApplication />}
                      />
                      <Route
                        path="metadata"
                        element={<SyncOktaSamlMetadata />}
                      />
                      <Route path="users" element={<AssignOktaSamlUsers />} />
                    </Route>
                    <Route
                      path="setup/other"
                      element={<OtherSamlConnectionFlow />}
                    >
                      <Route index element={<CreateOtherSamlApplication />} />
                      <Route
                        path="configure"
                        element={<ConfigureOtherSamlApplication />}
                      />
                      <Route
                        path="metadata"
                        element={<DownloadOtherSamlMetadata />}
                      />
                      <Route path="users" element={<AssignOtherSamlUsers />} />
                    </Route>
                    <Route path="test" element={<TestSamlConnectionDialog />} />
                  </Route>
                  <Route
                    path="oidc-connections"
                    element={<OidcConnectionsPage />}
                  />
                  <Route
                    path="oidc-connections/:oidcConnectionId"
                    element={<OidcConnectionPage />}
                  />
                  <Route
                    path="scim-api-keys/:scimApiKeyId"
                    element={<ScimApiKeyPage />}
                  />
                  <Route
                    path="users/:userId"
                    element={<OrganizationUserPage />}
                  >
                    <Route index element={<OrganizationUserDetailsTab />} />
                    <Route
                      path="roles"
                      element={<OrganizationUserRolesTab />}
                    />
                  </Route>
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
