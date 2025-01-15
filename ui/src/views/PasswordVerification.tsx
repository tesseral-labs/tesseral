import React, { useState } from 'react'
import { Title } from '@/components/Title'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import {
  exchangeIntermediateSessionForSession,
  verifyPassword,
} from '@/gen/openauth/intermediate/v1/intermediate-IntermediateService_connectquery'
import { useMutation } from '@connectrpc/connect-query'
import { useLocation, useNavigate, useParams } from 'react-router'
import { setAccessToken, setRefreshToken } from '@/auth'

const PasswordVerification = () => {
  const navigate = useNavigate()
  const { state } = useLocation()
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
        organizationId: state?.organizationId,
      })

      const { accessToken, refreshToken } =
        await exchangeIntermediateSessionForSessionMutation.mutateAsync({
          organizationId: state?.organizationId,
        })

      setAccessToken(accessToken)
      setRefreshToken(refreshToken)

      navigate('/session-info')
    } catch (error) {
      // TODO: Show an error message to the user
      console.error(error)
    }
  }

  return (
    <>
      <Title title="Verify Email Address" />

      <Card className="w-[clamp(320px,50%,420px)]">
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
            <input
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

export default PasswordVerification
