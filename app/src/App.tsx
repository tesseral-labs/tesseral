import React, { FC } from 'react'

import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { createConnectTransport } from '@connectrpc/connect-web'
import { type Transport } from '@connectrpc/connect'
import { TransportProvider } from '@connectrpc/connect-query'
import { BrowserRouter, Route, Routes } from 'react-router-dom'
import SessionInfoPage from './pages/SessionInfoPage'
import NotFoundPage from './pages/NotFound'
import { ListOrganizationsPage } from '@/pages/organizations/ListOrganizationsPage'
import { useAccessToken } from '@/lib/use-access-token'
import { Container } from '@/pages/Container'
import { ViewOrganizationPage } from '@/pages/organizations/ViewOrganizationPage'
import { ViewUserPage } from '@/pages/users/ViewUserPage'
import { ListProjectAPIKeysPage } from '@/pages/project-api-keys/ListProjectAPIKeysPage'
import { ViewProjectPage } from '@/pages/project/ViewProjectPage'
import { OrganizationUsersTab } from '@/pages/organizations/OrganizationUsersTab'
import { OrganizationSAMLConnectionsTab } from '@/pages/organizations/OrganizationSAMLConnectionsTab'
import { OrganizationSCIMAPIKeysTab } from '@/pages/organizations/OrganizationSCIMAPIKeysTab'
import { OrganizationDetailsTab } from '@/pages/organizations/OrganizationDetailsTab'
import { EditOrganizationPage } from '@/pages/organizations/EditOrganizationPage'
import { ViewSAMLConnectionPage } from '@/pages/saml-connections/ViewSAMLConnectionPage'
import { Toaster } from '@/components/ui/sonner'
import { Edit } from 'lucide-react'
import { EditSAMLConnectionPage } from '@/pages/saml-connections/EditSAMLConnectionPage'

const queryClient = new QueryClient()

function useTransport(): Transport {
  const accessToken = useAccessToken()

  return createConnectTransport({
    baseUrl: `http://auth.app.tesseral.example.com/api/internal/connect`,
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
        return next(req)
      },
    ],
  })
}

function AppWithinQueryClient() {
  const transport = useTransport()
  return (
    <TransportProvider transport={transport}>
      <BrowserRouter>
        <Routes>
          <Route path="" element={<Container />}>
            <Route path="" element={<ViewProjectPage />} />

            <Route
              path="project-api-keys"
              element={<ListProjectAPIKeysPage />}
            />

            <Route path="organizations" element={<ListOrganizationsPage />} />

            <Route
              path="organizations/:organizationId"
              element={<ViewOrganizationPage />}
            >
              <Route path="" element={<OrganizationDetailsTab />} />
              <Route path="users" element={<OrganizationUsersTab />} />
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
          </Route>
          <Route path="*" element={<NotFoundPage />} />
        </Routes>
      </BrowserRouter>
    </TransportProvider>
  )
}

const App: FC = () => {
  return (
    <QueryClientProvider client={queryClient}>
      <AppWithinQueryClient />
      <Toaster />
    </QueryClientProvider>
  )
}

export default App
