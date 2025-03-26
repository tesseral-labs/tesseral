import { Transport } from "@connectrpc/connect";
import { TransportProvider } from "@connectrpc/connect-query";
import { createConnectTransport } from "@connectrpc/connect-web";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import {
  AuthenticateAnotherWayPage,
  ChooseOrganizationPage,
  CreateOrganizationPage,
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
} from "@tesseral/common-ui";
import React from "react";
import { Route, Routes } from "react-router";
import { BrowserRouter } from "react-router-dom";

import { Toaster } from "@/components/ui/sonner";
import { NotFoundPage } from "@/pages/NotFoundPage";
import { LoggedInGate } from "@/pages/dashboard/LoggedInGate";
import { OrganizationAdvancedTab } from "@/pages/dashboard/OrganizationAdvancedTab";
import { OrganizationSettingsPage } from "@/pages/dashboard/OrganizationSettingsPage";
import { OrganizationUsersTab } from "@/pages/dashboard/OrganizationUsersTab";

import { DashboardLayout } from "./pages/dashboard/DashboardLayout";
import { UserSettingsPage } from "./pages/dashboard/UserSettingsPage";

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
            <Route path="" element={<LoginFlowLayout />}>
              <Route path="login" element={<LoginPage />} />
              <Route path="verify-email" element={<VerifyEmailPage />} />
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
                  <Route
                    path="advanced"
                    element={<OrganizationAdvancedTab />}
                  />
                </Route>
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
