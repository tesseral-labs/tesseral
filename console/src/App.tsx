import React, { FC } from 'react';

import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { createConnectTransport } from '@connectrpc/connect-web';
import { type Transport } from '@connectrpc/connect';
import { TransportProvider } from '@connectrpc/connect-query';
import { BrowserRouter, Route, Routes } from 'react-router-dom';
import NotFoundPage from './pages/NotFound';
import { ListOrganizationsPage } from '@/pages/organizations/ListOrganizationsPage';
import { useAccessToken } from '@/lib/use-access-token';
import { ViewOrganizationPage } from '@/pages/organizations/ViewOrganizationPage';
import { ViewUserPage } from '@/pages/users/ViewUserPage';
import { ListAPIKeysPage } from '@/pages/api-keys/ListAPIKeysPage';
import { ViewProjectSettingsPage } from '@/pages/project/ViewProjectSettingsPage';
import { OrganizationUsersTab } from '@/pages/organizations/OrganizationUsersTab';
import { OrganizationSAMLConnectionsTab } from '@/pages/organizations/OrganizationSAMLConnectionsTab';
import { OrganizationSCIMAPIKeysTab } from '@/pages/organizations/OrganizationSCIMAPIKeysTab';
import { OrganizationDetailsTab } from '@/pages/organizations/OrganizationDetailsTab';
import { EditOrganizationPage } from '@/pages/organizations/EditOrganizationPage';
import { ViewSAMLConnectionPage } from '@/pages/saml-connections/ViewSAMLConnectionPage';
import { Toaster } from '@/components/ui/sonner';
import { EditSAMLConnectionPage } from '@/pages/saml-connections/EditSAMLConnectionPage';
import { PageShell } from '@/components/page';
import { ViewSCIMAPIKeyPage } from '@/pages/scim-api-keys/ViewSCIMAPIKeyPage';
import { ViewProjectAPIKeyPage } from '@/pages/api-keys/ViewProjectAPIKey';
import { HomePage } from '@/pages/home/HomePage';
import { ProjectDetailsTab } from '@/pages/project/ProjectDetailsTab';
import LoginPage from './pages/login/LoginPage';
import { ViewPasskeyPage } from '@/pages/passkeys/ViewPasskeyPage';
import { OrganizationUserInvitesTab } from '@/pages/organizations/OrganizationUserInvitesTab';
import { ViewUserInvitePage } from '@/pages/user-invites/ViewUserInvitePage';
import { API_URL } from './config';
import { AuthType } from './lib/auth';
import GoogleOAuthCallbackPage from './pages/login/GoogleOAuthCallbackPage';
import MicrosoftOAuthCallbackPage from './pages/login/MicrosoftOAuthCallbackPage';
import { ViewPublishableKeyPage } from '@/pages/api-keys/ViewPublishableKeyPage';

const queryClient = new QueryClient();

const useTransport = (): Transport => {
  const accessToken = useAccessToken();

  return createConnectTransport({
    baseUrl: `${API_URL}/api/internal/connect`,
    fetch: (input, init) =>
      fetch(input, {
        ...init,
        headers: {
          Authorization: `Bearer ${accessToken}`,
          'Content-Type': 'application/json',
        },
        credentials: 'include',
      }),
    interceptors: [
      (next) => async (req) => {
        return next(req);
      },
    ],
  });
};

const AppWithinQueryClient = () => {
  const transport = useTransport();
  return (
    <TransportProvider transport={transport}>
      <BrowserRouter>
        <Routes>
          <Route
            path="/google-oauth-callback"
            element={<GoogleOAuthCallbackPage />}
          />
          <Route path="/login" element={<LoginPage />} />
          <Route
            path="/microsoft-oauth-callback"
            element={<MicrosoftOAuthCallbackPage />}
          />
          <Route
            path="/signup"
            element={<LoginPage authType={AuthType.SignUp} />}
          />

          <Route path="/" element={<PageShell />}>
            <Route path="" element={<HomePage />} />
            <Route
              path="project-settings"
              element={<ViewProjectSettingsPage />}
            >
              <Route path="" element={<ProjectDetailsTab />} />
            </Route>

            <Route
              path="project-settings/api-keys"
              element={<ListAPIKeysPage />}
            />

            <Route
              path="project-settings/api-keys/publishable-keys/:publishableKeyId"
              element={<ViewPublishableKeyPage />}
            />

            <Route
              path="project-settings/api-keys/project-api-keys/:projectApiKeyId"
              element={<ViewProjectAPIKeyPage />}
            />

            <Route path="organizations" element={<ListOrganizationsPage />} />

            <Route
              path="organizations/:organizationId"
              element={<ViewOrganizationPage />}
            >
              <Route path="" element={<OrganizationDetailsTab />} />
              <Route path="users" element={<OrganizationUsersTab />} />
              <Route
                path="user-invites"
                element={<OrganizationUserInvitesTab />}
              />
              <Route
                path="saml-connections"
                element={<OrganizationSAMLConnectionsTab />}
              />
              <Route
                path="scim-api-keys"
                element={<OrganizationSCIMAPIKeysTab />}
              />
            </Route>

            <Route
              path="organizations/:organizationId/edit"
              element={<EditOrganizationPage />}
            />
            <Route
              path="organizations/:organizationId/saml-connections/:samlConnectionId"
              element={<ViewSAMLConnectionPage />}
            />
            <Route
              path="organizations/:organizationId/saml-connections/:samlConnectionId/edit"
              element={<EditSAMLConnectionPage />}
            />
            <Route
              path="organizations/:organizationId/users/:userId"
              element={<ViewUserPage />}
            />
            <Route
              path="organizations/:organizationId/users/:userId/passkeys/:passkeyId"
              element={<ViewPasskeyPage />}
            />
            <Route
              path="organizations/:organizationId/user-invites/:userInviteId"
              element={<ViewUserInvitePage />}
            />
            <Route
              path="organizations/:organizationId/scim-api-keys/:scimApiKeyId"
              element={<ViewSCIMAPIKeyPage />}
            />
          </Route>

          <Route path="*" element={<NotFoundPage />} />
        </Routes>
      </BrowserRouter>
    </TransportProvider>
  );
};

const App: FC = () => {
  return (
    <QueryClientProvider client={queryClient}>
      <AppWithinQueryClient />
      <Toaster />
    </QueryClientProvider>
  );
};

export default App;
