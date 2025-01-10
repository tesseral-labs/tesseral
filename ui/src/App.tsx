import React, { FC } from 'react'
import { Route, Routes } from 'react-router'
import { BrowserRouter } from 'react-router-dom'

import { Transport } from '@connectrpc/connect'
import { TransportProvider } from '@connectrpc/connect-query'
import { createConnectTransport } from '@connectrpc/connect-web'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'

import { API_URL, PROJECT_ID } from '@/config'

import GoogleOAuthCallbackPage from '@/pages/GoogleOAuthCallbackPage'
import LoginPage from '@/pages/LoginPage'
import MicrosoftOAuthCallbackPage from '@/pages/MicrosoftOAuthCallbackPage'
import NotFoundPage from '@/pages/NotFound'
import SessionInfoPage from '@/pages/SessionInfoPage'

import Page from '@/components/Page'

const queryClient = new QueryClient()

function useTransport(): Transport {
  return createConnectTransport({
    baseUrl: `${API_URL}/internal/connect`,
    fetch: (input, init) => fetch(input, { ...init, credentials: 'include' }),
    interceptors: [
      (next) => async (req) => {
        // TODO: When we figure out how to get the project ID from the server, we should remove this logic.
        req.header.set('X-TODO-OpenAuth-Project-ID', PROJECT_ID)

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
              <Route
                path="/google-oauth-callback"
                element={<GoogleOAuthCallbackPage />}
              />
              <Route
                path="/microsoft-oauth-callback"
                element={<MicrosoftOAuthCallbackPage />}
              />
              <Route path="/login" element={<LoginPage />} />
              <Route path="/session-info" element={<SessionInfoPage />} />
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
