import React, { Dispatch, FC, SetStateAction, useState } from 'react'
import { Title } from '@/components/Title'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import {
  exchangeIntermediateSessionForSession,
  verifyPassword,
  whoami,
} from '@/gen/openauth/intermediate/v1/intermediate-IntermediateService_connectquery'
import { useMutation, useQuery } from '@connectrpc/connect-query'
import { useLocation, useNavigate, useParams } from 'react-router'
import { setAccessToken, setRefreshToken } from '@/auth'
import { Input } from '@/components/ui/input'
import { useIntermediateOrganization } from '@/lib/auth'
import { LoginViews } from '@/lib/views'

interface VerifyPasswordProps {
  setView: Dispatch<SetStateAction<LoginViews>>
}

const VerifyPassword: FC<VerifyPasswordProps> = ({ setView }) => {
  const organization = useIntermediateOrganization()
  const navigate = useNavigate()
  const { state } = useLocation()
  const [password, setPassword] = useState<string>('')

  const { data: whoamiRes } = useQuery(whoami)

  const exchangeIntermediateSessionForSessionMutation = useMutation(
    exchangeIntermediateSessionForSession,
  )
  const verifyPasswordMutation = useMutation(verifyPassword)

  const deriveNextView = (): LoginViews | undefined => {
    console.log(`organization`, organization)

    if (
      organization?.requireMfa &&
      !whoamiRes?.intermediateSession?.googleUserId &&
      !whoamiRes?.intermediateSession?.microsoftUserId
    ) {
      return LoginViews.ChooseAdditionalFactor
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    try {
      await verifyPasswordMutation.mutateAsync({
        password,
        organizationId: state?.organizationId,
      })

      const nextView = deriveNextView()

      console.log(`nextView: ${nextView}`)

      if (!!nextView) {
        setView(nextView)
        return
      }

      const { accessToken, refreshToken } =
        await exchangeIntermediateSessionForSessionMutation.mutateAsync({
          organizationId: state?.organizationId,
        })

      setAccessToken(accessToken)
      setRefreshToken(refreshToken)

      navigate('/settings')
    } catch (error) {
      // TODO: Show an error message to the user
      console.error(error)
    }
  }

  return (
    <>
      <Title title="Verify Email Address" />

      <Card>
        <CardHeader>
          <CardTitle className="text-center">Password Verification</CardTitle>
          <p className="text-sm text-center mt-2 text-gray-500">
            Please enter your password to continue logging in.
          </p>
        </CardHeader>
        <CardContent className="flex flex-col items-center justify-center w-full">
          <form
            className="flex flex-col items-center w-full"
            onSubmit={handleSubmit}
          >
            <Input
              className="w-full mb-2"
              id="password"
              placeholder="Enter your password"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
            />
            <Button
              className="w-full"
              disabled={password.length < 1}
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

export default VerifyPassword
