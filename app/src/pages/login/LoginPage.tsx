import React, { FC, useState } from 'react'

import ChooseProjectView from '@/views/login/ChooseProjectView'
import CreateProjectView from '@/views/login/CreateProjectView'
import StartLoginView from '@/views/login/StartLoginView'
import VerifyEmailView from '@/views/login/VerifyEmailView'
import VerifyPasswordView from '@/views/login/VerifyPasswordView'
import { LoginView } from '@/lib/views'

const LoginPage: FC = () => {
  const [view, setView] = useState<LoginView>(LoginView.StartLogin)

  return (
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
          {view === LoginView.ChooseProject ? (
            <ChooseProjectView setView={setView} />
          ) : null}
          {view === LoginView.CreateProject ? (
            <CreateProjectView setView={setView} />
          ) : null}
          {view === LoginView.StartLogin ? (
            <StartLoginView setView={setView} />
          ) : null}
          {view === LoginView.VerifyEmail ? (
            <VerifyEmailView setView={setView} />
          ) : null}
          {view === LoginView.VerifyPassword ? (
            <VerifyPasswordView setView={setView} />
          ) : null}
        </div>
      </div>
    </div>
  )
}

export default LoginPage
