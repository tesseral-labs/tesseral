import React, { Dispatch, FC, useState } from 'react'
import { LoginView } from '@/lib/views'
import { Title } from '@/components/Title'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Navigate, useLocation, useNavigate } from 'react-router'
import {
  exchangeIntermediateSessionForSession,
  verifyPassword,
} from '@/gen/openauth/intermediate/v1/intermediate-IntermediateService_connectquery'
import { useIntermediateSession } from '@/lib/auth'
import { useMutation } from '@connectrpc/connect-query'
import { setAccessToken, setRefreshToken } from '@/auth'
import { Input } from '@/components/ui/input'

interface VerifyPasswordViewProps {
  setView: Dispatch<React.SetStateAction<LoginView>>
}

const VerifyPasswordView: FC<VerifyPasswordViewProps> = ({ setView }) => {
  const intermediateSession = useIntermediateSession()
  const { state } = useLocation()
  const navigate = useNavigate()

  const [password, setPassword] = useState<string>('')

  const exchangeIntermediateSessionForSessionMutation = useMutation(
    exchangeIntermediateSessionForSession,
  )
  const verifyPasswordMutation = useMutation(verifyPassword)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    try {
      await verifyPasswordMutation.mutateAsync({
        password,
        organizationId: intermediateSession?.organizationId,
      })

      const { accessToken, refreshToken } =
        await exchangeIntermediateSessionForSessionMutation.mutateAsync({
          organizationId: state?.organizationId,
        })

      setAccessToken(accessToken)
      setRefreshToken(refreshToken)

      navigate('/project-settings')
    } catch (error) {
      // TODO: Show an error message to the user
      console.error(error)
    }
  }

  return (
    <>
      <Title title="Verify Email Address" />

      <Card className="w-[clamp(320px,50%,420px)] mx-auto">
        <CardHeader>
          <CardTitle className="text-center uppercase text-foreground font-semibold text-sm tracking-wide mt-2">
            Password Verification
          </CardTitle>
          <p className="text-sm text-center mt-2 text-gray-500">
            Please enter your password to continue logging in.
          </p>
        </CardHeader>
        <CardContent className="flex flex-col items-center justify-center w-full">
          <form className="flex flex-col items-center" onSubmit={handleSubmit}>
            <Input
              className="text-sm bg-input rounded border border-border focus:border-primary w-[clamp(240px,50%,100%)] mb-2"
              id="password"
              placeholder="Enter your password"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
            />
            <Button
              className="text-sm rounded border border-border focus:border-primary w-[clamp(240px,50%,100%)] mb-2"
              type="submit"
            >
              Log In
            </Button>
          </form>
        </CardContent>
      </Card>
    </>
  )
}

export default VerifyPasswordView
