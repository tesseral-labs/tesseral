import React, { useState } from 'react'

import { LoginViews } from '@/lib/views'
import { useLocation } from 'react-router'
import CreateOrganization from '@/views/CreateOrganization'
import Organizations from '@/views/Organizations'
import EmailVerification from '@/views/EmailVerification'
import Login from '@/views/Login'
import PasswordVerification from '@/views/PasswordVerification'

const LoginPage = () => {
  const location = useLocation()
  const { view = LoginViews.Login } =
    (location.state as { view: LoginViews }) || {}

  return (
    <>
      {view === LoginViews.CreateOrganization && <CreateOrganization />}
      {view === LoginViews.EmailVerification && <EmailVerification />}
      {view === LoginViews.Login && <Login />}
      {view === LoginViews.Organizations && <Organizations />}
      {view === LoginViews.PasswordVerification && <PasswordVerification />}
    </>
  )
}

export default LoginPage
