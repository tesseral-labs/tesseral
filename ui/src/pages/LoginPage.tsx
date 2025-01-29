import React, { useState } from 'react'

import { LoginViews } from '@/lib/views'
import { useLocation } from 'react-router'
import CreateOrganization from '@/views/CreateOrganization'
import Organizations from '@/views/Organizations'
import EmailVerification from '@/views/EmailVerification'
import Login from '@/views/Login'
import PasswordVerification from '@/views/PasswordVerification'
import RegisterPassword from '@/views/RegisterPassword'

const LoginPage = () => {
  const location = useLocation()
  const [view, setView] = useState<LoginViews>(LoginViews.Login)

  return (
    <>
      {view === LoginViews.CreateOrganization && (
        <CreateOrganization setView={setView} />
      )}
      {view === LoginViews.EmailVerification && (
        <EmailVerification setView={setView} />
      )}
      {view === LoginViews.Login && <Login setView={setView} />}
      {view === LoginViews.Organizations && <Organizations setView={setView} />}
      {view === LoginViews.PasswordVerification && <PasswordVerification />}
      {view === LoginViews.RegisterPassword && <RegisterPassword />}
    </>
  )
}

export default LoginPage
