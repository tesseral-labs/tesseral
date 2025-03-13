import { Transport } from "@connectrpc/connect";
import { TransportProvider } from "@connectrpc/connect-query";
import { createConnectTransport } from "@connectrpc/connect-web";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import React, { FC } from "react";
import { Route, Routes } from "react-router";
import { BrowserRouter } from "react-router-dom";

import Page from "@/components/Page";
import NotFoundPage from "@/pages/NotFound";

import DashboardPage from "./components/DashboardPage";
import EditSAMLConnectionsPage from "./pages/dashboard/EditSAMLConnectionsPage";
import OrganizationSettingsPage from "./pages/dashboard/OrganizationSettingsPage";
import UserSettingsPage from "./pages/dashboard/UserSettingsPage";

const queryClient = new QueryClient();

const useTransport = (): Transport => {
  return createConnectTransport({
    baseUrl: `/api/internal/connect`,
    fetch: (input, init) => fetch(input, { ...init, credentials: "include" }),
  });
};

const AppWithRoutes: FC = () => {
  const transport = useTransport();

  return (
    <TransportProvider transport={transport}>
      <QueryClientProvider client={queryClient}>
        <BrowserRouter>
          <Routes>
            <Route path="" element={<Page />}></Route>

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
};

const App: FC = () => {
  return <AppWithRoutes />;
};

export default App;
