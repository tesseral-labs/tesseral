import React, { FC, useEffect, useState } from 'react'

import {
  AuthType,
  AuthTypeContextProvider,
  IntermediateOrganizationContextProvider,
} from '@/lib/auth'
import { LoginView } from '@/lib/views'

import ChooseProjectView from '@/views/login/ChooseProjectView'
import CreateProjectView from '@/views/login/CreateProjectView'
import StartLoginView from '@/views/login/StartLoginView'
import VerifyEmailView from '@/views/login/VerifyEmailView'
import VerifyPasswordView from '@/views/login/VerifyPasswordView'
import { Organization } from '@/gen/openauth/intermediate/v1/intermediate_pb'
import ChooseAdditionalFactorView from '@/views/login/ChooseAdditionalFactorView'
import VerifyAuthenticatorAppView from '@/views/login/VerifyAuthenticatorAppView'
import VerifyPasskeyView from '@/views/login/VerifyPasskeyView'
import RegisterAuthenticatorAppView from '@/views/login/RegisterAuthenticatorAppView'
import RegisterPasskeyView from '@/views/login/RegisterPasskeyView'
import RegisterPasswordView from '@/views/login/RegisterPasswordView'
import { useSearchParams } from 'react-router-dom'

interface LoginPageProps {
  authType?: AuthType
}

const LoginPage: FC<LoginPageProps> = ({ authType = AuthType.LogIn }) => {
  const [searchParams] = useSearchParams()
  const [intermediateOrganization, setIntermediateOrganization] =
    useState<Organization>()
  const [view, setView] = useState<LoginView>(LoginView.StartLogin)

  useEffect(() => {
    if (searchParams.get('view')) {
      setView(searchParams.get('view') as LoginView)
    }
  }, [searchParams])

  return (
    <AuthTypeContextProvider value={authType}>
      <IntermediateOrganizationContextProvider value={intermediateOrganization}>
        <div className="w-screen min-h-screen flex flex-col justify-center items-center bg-indigo-600">
          <div className="min-w-[320px] max-w-[576px]">
            <div className="w-full mb-8">
              <img
                className="h-full max-h-10 mx-auto"
                src="/images/tesseral-logo-white.svg"
                alt="tesseral"
              />
            </div>
            <div className="w-full">
              {view === LoginView.ChooseAdditionalFactor ? (
                <ChooseAdditionalFactorView setView={setView} />
              ) : null}
              {view === LoginView.ChooseProject ? (
                <ChooseProjectView
                  setIntermediateOrganization={setIntermediateOrganization}
                  setView={setView}
                />
              ) : null}
              {view === LoginView.CreateProject ? (
                <CreateProjectView setView={setView} />
              ) : null}
              {view === LoginView.RegisterAuthenticatorApp ? (
                <RegisterAuthenticatorAppView />
              ) : null}
              {view === LoginView.RegisterPasskey ? (
                <RegisterPasskeyView />
              ) : null}
              {view === LoginView.RegisterPassword ? (
                <RegisterPasswordView setView={setView} />
              ) : null}
              {view === LoginView.StartLogin ? (
                <StartLoginView setView={setView} />
              ) : null}
              {view === LoginView.VerifyAuthenticatorApp ? (
                <VerifyAuthenticatorAppView />
              ) : null}
              {view === LoginView.VerifyEmail ? (
                <VerifyEmailView setView={setView} />
              ) : null}
              {view === LoginView.VerifyPasskey ? <VerifyPasskeyView /> : null}
              {view === LoginView.VerifyPassword ? (
                <VerifyPasswordView setView={setView} />
              ) : null}
            </div>
          </div>
        </div>
      </IntermediateOrganizationContextProvider>
    </AuthTypeContextProvider>
  )
}

export default LoginPage
