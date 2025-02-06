import React, { useState } from 'react'

import { LoginViews } from '@/lib/views'
import CreateOrganization from '@/views/CreateOrganization'
import ChooseOrganization from '@/views/ChooseOrganization'
import Login from '@/views/Login'
import RegisterPassword from '@/views/RegisterPassword'
import VerifyEmail from '@/views/VerifyEmail'
import VerifyPassword from '@/views/VerifyPassword'
import RegisterAuthenticatorApp from '@/views/RegisterAuthenticatorApp'
import RegisterPasskey from '@/views/RegisterPasskey'
import VerifyAuthenticatorApp from '@/views/VerifyAuthenticatorApp'
import VerifyPasskey from '@/views/VerifyPasskey'
import ChooseAdditionalFactor from '@/views/ChooseAdditionalFactor'
import { Organization } from '@/gen/openauth/intermediate/v1/intermediate_pb'
import { IntermediateOrganizationContextProvider } from '@/lib/auth'
import ChooseOrganizationPrimaryLoginFactor from '@/views/ChooseOrganizationPrimaryLoginFactor'

const LoginPage = () => {
  const [intermediateOrganization, setIntermediateOrganization] =
    useState<Organization>()
  const [view, setView] = useState<LoginViews>(LoginViews.Login)

  return (
    <IntermediateOrganizationContextProvider value={intermediateOrganization}>
      {view === LoginViews.ChooseAdditionalFactor && (
        <ChooseAdditionalFactor setView={setView} />
      )}
      {view === LoginViews.ChooseOrganization && (
        <ChooseOrganization
          setIntermediateOrganization={setIntermediateOrganization}
          setView={setView}
        />
      )}
      {view === LoginViews.ChooseOrganizationPrimaryLoginFactor && (
        <ChooseOrganizationPrimaryLoginFactor setView={setView} />
      )}
      {view === LoginViews.CreateOrganization && (
        <CreateOrganization setView={setView} />
      )}
      {view === LoginViews.Login && <Login setView={setView} />}
      {view === LoginViews.RegisterPassword && (
        <RegisterPassword setView={setView} />
      )}
      {view === LoginViews.RegisterAuthenticatorApp && (
        <RegisterAuthenticatorApp />
      )}
      {view === LoginViews.RegisterPasskey && <RegisterPasskey />}
      {view === LoginViews.VerifyEmail && <VerifyEmail setView={setView} />}
      {view === LoginViews.VerifyPassword && (
        <VerifyPassword setView={setView} />
      )}
      {view === LoginViews.VerifyAuthenticatorApp && <VerifyAuthenticatorApp />}
      {view === LoginViews.VerifyPasskey && <VerifyPasskey />}
    </IntermediateOrganizationContextProvider>
  )
}

export default LoginPage
