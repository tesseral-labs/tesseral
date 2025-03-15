import { Transport } from "@connectrpc/connect";
import { TransportProvider } from "@connectrpc/connect-query";
import { createConnectTransport } from "@connectrpc/connect-web";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import React from "react";
import { Route, Routes } from "react-router";
import { BrowserRouter } from "react-router-dom";

import { NotFoundPage } from "@/pages/NotFoundPage";
import { LoginFlowLayout } from "@/pages/login/LoginFlowLayout";
import { LoginPage } from "@/pages/login/LoginPage";
import { VerifyEmailPage } from "@/pages/login/VerifyEmailPage";

import { DashboardPage } from "./components/DashboardPage";
import { EditSAMLConnectionsPage } from "./pages/dashboard/EditSAMLConnectionsPage";
import { OrganizationSettingsPage } from "./pages/dashboard/OrganizationSettingsPage";
import { UserSettingsPage } from "./pages/dashboard/UserSettingsPage";
import { ChooseOrganizationPage } from "@/pages/login/ChooseOrganizationPage";
import { OrganizationLoginPage } from "@/pages/login/OrganizationLoginPage";
import { VerifyPasswordPage } from "@/pages/login/VerifyPasswordPage";
import { FinishLoginPage } from "@/pages/login/FinishLoginPage";
import { CreateOrganizationPage } from "@/pages/login/CreateOrganizationPage";
import {
  RegisterSecondaryFactorPage
} from "@/pages/login/RegisterSecondaryFactorPage";
import {
  RegisterAuthenticatorAppPage
} from "@/pages/login/RegisterAuthenticatorAppPage";
import { RegisterPasskeyPage } from "@/pages/login/RegisterPasskeyPage";
import { VerifyPasskeyPage } from "@/pages/login/VerifyPasskeyPage";

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
              <Route path="choose-organization" element={<ChooseOrganizationPage />} />
              <Route path="create-organization" element={<CreateOrganizationPage />} />
              <Route path="organizations/:organizationId/login" element={<OrganizationLoginPage />} />
              <Route path="verify-password" element={<VerifyPasswordPage />} />
              <Route path="verify-passkey" element={<VerifyPasskeyPage />} />
              <Route path="register-secondary-factor" element={<RegisterSecondaryFactorPage />} />
              <Route path="register-passkey" element={<RegisterPasskeyPage />} />
              <Route path="register-authenticator-app" element={<RegisterAuthenticatorAppPage />} />
              <Route path="finish-login" element={<FinishLoginPage />} />
            </Route>

            <Route
              path="/organization"
              element={
                <DashboardPage>
                  <OrganizationSettingsPage />
                </DashboardPage>
              }
            />
            <Route
              path="/organization/saml-connections/:samlConnectionId"
              element={
                <DashboardPage>
                  <EditSAMLConnectionsPage />
                </DashboardPage>
              }
            />
            <Route
              path="/settings"
              element={
                <DashboardPage>
                  <UserSettingsPage />
                </DashboardPage>
              }
            />
            <Route path="*" element={<NotFoundPage />} />
          </Routes>
        </BrowserRouter>
      </QueryClientProvider>
    </TransportProvider>
  );
}

export function App() {
  return <AppWithRoutes />;
}
