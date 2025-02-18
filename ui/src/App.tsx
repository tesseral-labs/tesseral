import React, { FC } from 'react';
import { Navigate, Route, Routes } from 'react-router';
import { BrowserRouter } from 'react-router-dom';

import { Transport } from '@connectrpc/connect';
import { TransportProvider } from '@connectrpc/connect-query';
import { createConnectTransport } from '@connectrpc/connect-web';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

import GoogleOAuthCallbackPage from '@/pages/GoogleOAuthCallbackPage';
import LoginPage from '@/pages/LoginPage';
import MicrosoftOAuthCallbackPage from '@/pages/MicrosoftOAuthCallbackPage';
import NotFoundPage from '@/pages/NotFound';

import Page from '@/components/Page';
import UserSettingsPage from './pages/dashboard/UserSettingsPage';
import DashboardPage from './components/DashboardPage';
import OrganizationSettingsPage from './pages/dashboard/OrganizationSettingsPage';
import EditSAMLConnectionsPage from './pages/dashboard/EditSAMLConnectionsPage';
import RegisterPasskey from './views/RegisterPasskey';
import { AuthType } from './lib/auth';

const queryClient = new QueryClient();

const useTransport = (): Transport => {
  return createConnectTransport({
    baseUrl: `/api/internal/connect`,
    fetch: (input, init) => fetch(input, { ...init, credentials: 'include' }),
  });
};

const AppWithRoutes: FC = () => {
  const transport = useTransport();

  return (
    <TransportProvider transport={transport}>
      <QueryClientProvider client={queryClient}>
        <BrowserRouter>
          <Routes>
            <Route path="/" element={<Page />}>
              <Route index element={<Navigate to="login" replace />} />

              <Route path="passkey-test" element={<RegisterPasskey />} />

              <Route
                path="google-oauth-callback"
                element={<GoogleOAuthCallbackPage />}
              />
              <Route
                path="microsoft-oauth-callback"
                element={<MicrosoftOAuthCallbackPage />}
              />
              <Route path="login" element={<LoginPage />} />
              <Route
                path="signup"
                element={<LoginPage authType={AuthType.SignUp} />}
              />
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
};

const App: FC = () => {
  return <AppWithRoutes />;
};

export default App;
