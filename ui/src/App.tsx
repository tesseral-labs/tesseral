import React, { FC } from 'react'
import { Route, Routes } from 'react-router'
import { BrowserRouter } from 'react-router-dom'

import { Transport } from '@connectrpc/connect'
import { TransportProvider } from '@connectrpc/connect-query'
import { createConnectTransport } from '@connectrpc/connect-web'

import { getIntermediateSessionToken } from './auth'

import LoginPage from '@/pages/LoginPage'
import NotFoundPage from '@/pages/NotFound'
import OrganizationsPage from '@/pages/OrganizationsPage'
import Page from './components/Page'
import EmailVerificationPage from './pages/EmailVerificationPage'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { API_URL } from './config'

const queryClient = new QueryClient()

function useTransport(): Transport {
  return createConnectTransport({
    baseUrl: API_URL,
    interceptors: [
      (next) => async (req) => {
        req.header.set(
          'Authorization',
          `Bearer ${getIntermediateSessionToken() ?? ''}`,
        )
        return next(req)
      },
    ],
  })
}

const AppWithRoutes: FC = () => {
  const transport = useTransport()

  return (
    <TransportProvider transport={transport}>
      <QueryClientProvider client={queryClient}>
        <BrowserRouter>
          <Routes>
            <Route path="/" element={<Page />}>
              <Route path="/login" element={<LoginPage />} />
              <Route path="/organizations" element={<OrganizationsPage />} />
              <Route path="/verify-email" element={<EmailVerificationPage />} />
            </Route>

            <Route path="*" element={<NotFoundPage />} />
          </Routes>
        </BrowserRouter>
      </QueryClientProvider>
    </TransportProvider>
  )
}

const App: FC = () => {
  return <AppWithRoutes />
}

export default App
