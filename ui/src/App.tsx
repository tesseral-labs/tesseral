import React, { FC } from 'react'
import { Route, Routes } from 'react-router'
import { BrowserRouter } from 'react-router-dom'

import LoginPage from '@/pages/LoginPage'
import NotFoundPage from '@/pages/NotFound'
import OrganizationsPage from '@/pages/OrganizationsPage'
import Page from './components/Page'

const App: FC = () => {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Page />}>
          <Route path="/login" element={<LoginPage />} />
          <Route path="/organizations" element={<OrganizationsPage />} />
        </Route>

        <Route path="*" element={<NotFoundPage />} />
      </Routes>
    </BrowserRouter>
  )
}

export default App
