import { Transport } from "@connectrpc/connect";
import { TransportProvider } from "@connectrpc/connect-query";
import { createConnectTransport } from "@connectrpc/connect-web";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import React from "react";
import { Navigate, Route, Routes } from "react-router";
import { BrowserRouter } from "react-router-dom";

import { Toaster } from "@/components/ui/sonner";
import { NotFoundPage } from "@/pages/NotFoundPage";
import { LoggedInGate } from "@/pages/dashboard/LoggedInGate";
import { OrganizationSettingsPage } from "@/pages/dashboard/OrganizationSettingsPage";
import { OrganizationUsersTab } from "@/pages/dashboard/OrganizationUsersTab";
import { CreateRolePage } from "@/pages/dashboard/roles/CreateRolePage";
import { EditRolePage } from "@/pages/dashboard/roles/EditRolePage";
import { AuthenticateAnotherWayPage } from "@/pages/login/AuthenticateAnotherWayPage";
import { ChooseOrganizationPage } from "@/pages/login/ChooseOrganizationPage";
import { CreateOrganizationPage } from "@/pages/login/CreateOrganizationPage";
import { FinishLoginPage } from "@/pages/login/FinishLoginPage";
import { ForgotPasswordPage } from "@/pages/login/ForgotPasswordPage";
import { GoogleOAuthCallbackPage } from "@/pages/login/GoogleOAuthCallbackPage";
import { ImpersonatePage } from "@/pages/login/ImpersonatePage";
import { LoginFlowLayout } from "@/pages/login/LoginFlowLayout";
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

import { APIKeysTab } from "./pages/dashboard/APIKeysTab";
import { DashboardLayout } from "./pages/dashboard/DashboardLayout";
import { OrganizationLoginSettingsTab } from "./pages/dashboard/OrganizationLoginSettingsTab";
import { OrganizationSAMLConnectionsTab } from "./pages/dashboard/OrganizationSAMLConnectionsTab";
import { UserSettingsPage } from "./pages/dashboard/UserSettingsPage";
import { ViewAPIKeyPage } from "./pages/dashboard/api-keys/ViewAPIKeyPage";
import { ViewSAMLConnectionPage } from "./pages/dashboard/saml-connections/ViewSAMLConnectionPage";
import { GithubOAuthCallbackPage } from "./pages/login/GithubOAuthCallbackPage";

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
              path="/"
              element={<Navigate to="/user-settings" replace />}
            />
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
              <Route path="" element={<DashboardLayout />}>
                <Route path="user-settings" element={<UserSettingsPage />} />
                <Route
                  path="organization-settings"
                  element={<OrganizationSettingsPage />}
                >
                  <Route path="" element={<OrganizationUsersTab />} />
                  <Route path="api-keys" element={<APIKeysTab />} />
                  <Route
                    path="api-keys/:apiKeyId"
                    element={<ViewAPIKeyPage />}
                  />
                  <Route
                    path="login-settings"
                    element={<OrganizationLoginSettingsTab />}
                  />
                  <Route
                    path="saml-connections"
                    element={<OrganizationSAMLConnectionsTab />}
                  />
                  <Route
                    path="saml-connections/:samlConnectionId"
                    element={<ViewSAMLConnectionPage />}
                  />
                </Route>

                <Route
                  path="organization-settings/roles/new"
                  element={<CreateRolePage />}
                />
                <Route
                  path="organization-settings/roles/:roleId/edit"
                  element={<EditRolePage />}
                />
              </Route>
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
