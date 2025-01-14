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
          <Route path="/" element={<SessionInfoPage />} />
          <Route path="/organizations" element={<ListOrganizationsPage />} />
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
