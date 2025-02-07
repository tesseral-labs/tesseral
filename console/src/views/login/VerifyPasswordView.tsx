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
import {
  AuthType,
  useAuthType,
  useIntermediateOrganization,
  useIntermediateSession,
} from '@/lib/auth'
import { useMutation } from '@connectrpc/connect-query'
import { setAccessToken, setRefreshToken } from '@/auth'
import { Input } from '@/components/ui/input'
import { parseErrorMessage } from '@/lib/errors'
import { toast } from 'sonner'

interface VerifyPasswordViewProps {
  setView: Dispatch<React.SetStateAction<LoginView>>
}

const VerifyPasswordView: FC<VerifyPasswordViewProps> = ({ setView }) => {
  const authType = useAuthType()
  const intermediateSession = useIntermediateSession()
  const organization = useIntermediateOrganization()

  const { state } = useLocation()
  const navigate = useNavigate()

  const [password, setPassword] = useState<string>('')

  const exchangeIntermediateSessionForSessionMutation = useMutation(
    exchangeIntermediateSessionForSession,
  )
  const verifyPasswordMutation = useMutation(verifyPassword)

  const deriveNextView = (): LoginView | undefined => {
    console.log(`organization:`, organization)

    const hasMultipleSecondFactors =
      organization?.userHasAuthenticatorApp && organization?.userHasPasskey
    const hasSecondFactor =
      organization?.userHasAuthenticatorApp || organization?.userHasPasskey

    if (organization?.requireMfa) {
      if (hasMultipleSecondFactors || !hasSecondFactor) {
        return LoginView.ChooseAdditionalFactor
      } else if (organization?.userHasPasskey) {
        return LoginView.VerifyPasskey
      } else if (organization?.userHasAuthenticatorApp) {
        return LoginView.VerifyAuthenticatorApp
      }
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    try {
      await verifyPasswordMutation.mutateAsync({
        password,
        organizationId: intermediateSession?.organizationId,
      })

      const nextView = deriveNextView()

      console.log('nextView:', nextView)

      if (nextView) {
        setView(nextView)
        return
      }

      const { accessToken, refreshToken } =
        await exchangeIntermediateSessionForSessionMutation.mutateAsync({
          organizationId: state?.organizationId,
        })

      setAccessToken(accessToken)
      setRefreshToken(refreshToken)

      navigate('/project-settings')
    } catch (error) {
      const message = parseErrorMessage(error)
      toast.error('Could not verify password', {
        description: message,
      })
    }
  }

  return (
    <>
      <Title title="Verify Email Address" />

      <Card className="w-full max-w-sm">
        <CardHeader>
          <CardTitle>Password Verification</CardTitle>
          <p className="text-sm text-center mt-2 text-gray-500">
            Please enter your password to continue{' '}
            {authType === AuthType.SignUp ? 'signing up' : 'logging in'}.
          </p>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit}>
            <Input
              className="w-full"
              id="password"
              placeholder="Enter your password"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
            />
            <Button className="mt-2 w-full" type="submit">
              Continue
            </Button>
          </form>
        </CardContent>
      </Card>
    </>
  )
}

export default VerifyPasswordView
