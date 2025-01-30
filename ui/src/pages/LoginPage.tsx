import React, { useState } from 'react'

import { LoginViews } from '@/lib/views'
import CreateOrganization from '@/views/CreateOrganization'
import Organizations from '@/views/Organizations'
import Login from '@/views/Login'
import RegisterPassword from '@/views/RegisterPassword'
import VerifyEmail from '@/views/VerifyEmail'
import VerifyPassword from '@/views/VerifyPassword'
import RegisterAuthenticatorApp from '@/views/RegisterAuthenticatorApp'
import RegisterPasskey from '@/views/RegisterPasskey'
import VerifyAuthenticatorApp from '@/views/VerifyAuthenticatorApp'
import VerifyPasskey from '@/views/VerifyPasskey'

const LoginPage = () => {
  const [view, setView] = useState<LoginViews>(LoginViews.Login)

  return (
    <>
      {view === LoginViews.CreateOrganization && (
        <CreateOrganization setView={setView} />
      )}
      {view === LoginViews.Login && <Login setView={setView} />}
      {view === LoginViews.Organizations && <Organizations setView={setView} />}
      {view === LoginViews.RegisterPassword && <RegisterPassword />}
      {view === LoginViews.RegisterAuthenticatorApp && (
        <RegisterAuthenticatorApp />
      )}
      {view === LoginViews.RegisterPasskey && <RegisterPasskey />}
      {view === LoginViews.VerifyEmail && <VerifyEmail setView={setView} />}
      {view === LoginViews.VerifyPassword && <VerifyPassword />}
      {view === LoginViews.VerifyAuthenticatorApp && <VerifyAuthenticatorApp />}
      {view === LoginViews.VerifyPasskey && <VerifyPasskey />}
    </>
  )
}

export default LoginPage
