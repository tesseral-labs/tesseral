import { Transport } from "@connectrpc/connect";
import { TransportProvider } from "@connectrpc/connect-query";
import { createConnectTransport } from "@connectrpc/connect-web";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import React from "react";
import { Route, Routes } from "react-router";
import { BrowserRouter } from "react-router-dom";

import { Toaster } from "@/components/ui/sonner";
import { NotFoundPage } from "@/pages/NotFoundPage";
import { LoggedInGate } from "@/pages/dashboard/LoggedInGate";
import { OrganizationSettingsPage } from "@/pages/dashboard/OrganizationSettingsPage";
import { ChooseOrganizationPage } from "@/pages/login/ChooseOrganizationPage";
import { CreateOrganizationPage } from "@/pages/login/CreateOrganizationPage";
import { FinishLoginPage } from "@/pages/login/FinishLoginPage";
import { ImpersonatePage } from "@/pages/login/ImpersonatePage";
import { LoginFlowLayout } from "@/pages/login/LoginFlowLayout";
import { LoginPage } from "@/pages/login/LoginPage";
import { OrganizationLoginPage } from "@/pages/login/OrganizationLoginPage";
import { RegisterAuthenticatorAppPage } from "@/pages/login/RegisterAuthenticatorAppPage";
import { RegisterPasskeyPage } from "@/pages/login/RegisterPasskeyPage";
import { RegisterSecondaryFactorPage } from "@/pages/login/RegisterSecondaryFactorPage";
import { SwitchOrganizationsPage } from "@/pages/login/SwitchOrganizationsPage";
import { VerifyAuthenticatorAppPage } from "@/pages/login/VerifyAuthenticatorAppPage";
import { VerifyAuthenticatorAppRecoveryCodePage } from "@/pages/login/VerifyAuthenticatorAppRecoveryCodePage";
import { VerifyEmailPage } from "@/pages/login/VerifyEmailPage";
import { VerifyPasskeyPage } from "@/pages/login/VerifyPasskeyPage";
import { VerifyPasswordPage } from "@/pages/login/VerifyPasswordPage";
import { VerifySecondaryFactorPage } from "@/pages/login/VerifySecondaryFactorPage";

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
            <Route path="login" element={<LoginPage />} />
            <Route path="" element={<LoginFlowLayout />}>
              <Route path="verify-email" element={<VerifyEmailPage />} />
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
            </Route>

            <Route path="" element={<LoggedInGate />}>
              <Route path="" element={<DashboardLayout />}>
                <Route path="user-settings" element={<UserSettingsPage />} />
                <Route
                  path="organization-settings"
                  element={<OrganizationSettingsPage />}
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
