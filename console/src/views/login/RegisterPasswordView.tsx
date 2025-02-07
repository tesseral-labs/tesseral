import React, { FC, useState } from 'react'
import { Title } from '@/components/Title'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import {
  exchangeIntermediateSessionForSession,
  registerPassword,
} from '@/gen/openauth/intermediate/v1/intermediate-IntermediateService_connectquery'
import { useMutation } from '@connectrpc/connect-query'
import { useNavigate } from 'react-router'
import { setAccessToken, setRefreshToken } from '@/auth'
import { Input } from '@/components/ui/input'
import { LoginView } from '@/lib/views'
import { parseErrorMessage } from '@/lib/errors'
import { toast } from 'sonner'
import { AuthType, useAuthType, useIntermediateOrganization } from '@/lib/auth'
import Loader from '@/components/ui/loader'

interface RegisterPasswordViewProps {
  setView: React.Dispatch<React.SetStateAction<LoginView>>
}

const RegisterPasswordView: FC<RegisterPasswordViewProps> = ({ setView }) => {
  const authType = useAuthType()
  const organization = useIntermediateOrganization()
  const navigate = useNavigate()

  const [password, setPassword] = useState<string>('')
  const [submitting, setSubmitting] = useState<boolean>(false)

  const exchangeIntermediateSessionForSessionMutation = useMutation(
    exchangeIntermediateSessionForSession,
  )
  const registerPasswordMutation = useMutation(registerPassword)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setSubmitting(true)

    try {
      await registerPasswordMutation.mutateAsync({
        password,
      })

      if (organization?.requireMfa) {
        setView(LoginView.ChooseAdditionalFactor)
        return
      }

      const { accessToken, refreshToken } =
        await exchangeIntermediateSessionForSessionMutation.mutateAsync({})

      setAccessToken(accessToken)
      setRefreshToken(refreshToken)
      setSubmitting(false)

      navigate('/settings')
    } catch (error) {
      setSubmitting(false)
      const message = parseErrorMessage(error)
      toast.error('Could not set password', {
        description: message,
      })
    }
  }

  return (
    <>
      <Title title="Set your password" />

      <Card className="w-full max-w-sm">
        <CardHeader>
          <CardTitle>Set your password</CardTitle>
          <CardDescription>
            Please set your password to continue{' '}
            {authType === AuthType.SignUp ? 'signing up' : 'logging in'}.
          </CardDescription>
        </CardHeader>
        <CardContent>
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
              disabled={password.length < 1 || submitting}
              type="submit"
            >
              {submitting && <Loader />}
              Continue
            </Button>
          </form>
        </CardContent>
      </Card>
    </>
  )
}

export default RegisterPasswordView
