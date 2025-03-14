import { Transport } from "@connectrpc/connect";
import { TransportProvider } from "@connectrpc/connect-query";
import { createConnectTransport } from "@connectrpc/connect-web";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import React from "react";
import { Route, Routes } from "react-router";
import { BrowserRouter } from "react-router-dom";

import { NotFoundPage } from "@/pages/NotFoundPage";
import { LoginPage } from "@/pages/login/LoginPage";

import { DashboardPage } from "./components/DashboardPage";
import { EditSAMLConnectionsPage } from "./pages/dashboard/EditSAMLConnectionsPage";
import { OrganizationSettingsPage } from "./pages/dashboard/OrganizationSettingsPage";
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
