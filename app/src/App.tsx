import React, { FC } from 'react'

import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { API_URL, DOGFOOD_PROJECT_ID } from './config'
import { createConnectTransport } from '@connectrpc/connect-web'
import { type Transport } from '@connectrpc/connect'
import { TransportProvider } from '@connectrpc/connect-query'
import { BrowserRouter, Route, Routes } from 'react-router-dom'
import SessionInfoPage from './pages/SessionInfoPage'
import NotFoundPage from './pages/NotFound'

const queryClient = new QueryClient()

function useTransport(): Transport {
  return createConnectTransport({
    baseUrl: `${API_URL}/api/internal/connect`,
    fetch: (input, init) => fetch(input, { ...init, credentials: 'include' }),
    interceptors: [
      (next) => async (req) => {
        // TODO: When we figure out how to get the project ID from the server, we should remove this logic.
        req.header.set('X-TODO-OpenAuth-Project-ID', DOGFOOD_PROJECT_ID)

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
            <Route path="/" element={<SessionInfoPage />} />
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
