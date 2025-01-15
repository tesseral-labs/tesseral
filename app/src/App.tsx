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
            <Route path="/" element={<ViewProjectPage />} />

            <Route
              path="/project-api-keys"
              element={<ListProjectAPIKeysPage />}
            />

            <Route path="/organizations" element={<ListOrganizationsPage />} />
            <Route
              path="/organizations/:organizationId"
              element={<ViewOrganizationPage />}
            />

            <Route
              path="/organizations/:organizationId/users/:userId"
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
    </QueryClientProvider>
  )
}

export default App
