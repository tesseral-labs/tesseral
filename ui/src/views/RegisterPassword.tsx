import React, { FC, useState } from 'react'
import { Title } from '@/components/Title'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import {
  exchangeIntermediateSessionForSession,
  registerPassword,
} from '@/gen/openauth/intermediate/v1/intermediate-IntermediateService_connectquery'
import { useMutation } from '@connectrpc/connect-query'
import { useNavigate } from 'react-router'
import { setAccessToken, setRefreshToken } from '@/auth'
import { Input } from '@/components/ui/input'
import { useLayout } from '@/lib/settings'
import { cn } from '@/lib/utils'
import { LoginLayouts, LoginViews } from '@/lib/views'
import { parseErrorMessage } from '@/lib/errors'
import { toast } from 'sonner'
import { useIntermediateOrganization } from '@/lib/auth'

interface RegisterPasswordProps {
  setView: React.Dispatch<React.SetStateAction<LoginViews>>
}

const RegisterPassword: FC<RegisterPasswordProps> = ({ setView }) => {
  const organization = useIntermediateOrganization()
  const layout = useLayout()
  const navigate = useNavigate()
  const [password, setPassword] = useState<string>('')

  const exchangeIntermediateSessionForSessionMutation = useMutation(
    exchangeIntermediateSessionForSession,
  )
  const registerPasswordMutation = useMutation(registerPassword)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    try {
      await registerPasswordMutation.mutateAsync({
        password,
      })

      if (organization?.requireMfa) {
        setView(LoginViews.ChooseAdditionalFactor)
        return
      }

      const { accessToken, refreshToken } =
        await exchangeIntermediateSessionForSessionMutation.mutateAsync({})

      setAccessToken(accessToken)
      setRefreshToken(refreshToken)

      navigate('/settings')
    } catch (error) {
      const message = parseErrorMessage(error)

      toast.error('Could not set password', {
        description: message,
      })
    }
  }

  return (
    <>
      <Title title="Set your password" />

      <Card
        className={cn(
          'w-full max-w-sm',
          layout !== LoginLayouts.Centered && 'shadow-none border-0',
        )}
      >
        <CardHeader>
          <CardTitle className="text-center">Set your password</CardTitle>
          <p className="text-sm text-center mt-2 text-gray-500">
            Please set your password to continue signing up.
          </p>
        </CardHeader>
        <CardContent className="flex flex-col items-center justify-center w-full">
          <form
            className="flex flex-col items-center w-full"
            onSubmit={handleSubmit}
          >
            <Input
              className="w-full"
              id="password"
              placeholder="Enter your password"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
            />
            <Button
              className="w-full mt-4"
              disabled={password.length < 1}
              type="submit"
            >
              Set password
            </Button>
          </form>
        </CardContent>
      </Card>
    </>
  )
}

export default RegisterPassword
