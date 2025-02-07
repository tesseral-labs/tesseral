import { setAccessToken, setRefreshToken } from '@/auth'
import { Title } from '@/components/Title'
import { Button } from '@/components/ui/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import {
  InputOTP,
  InputOTPGroup,
  InputOTPSeparator,
  InputOTPSlot,
} from '@/components/ui/input-otp'
import Loader from '@/components/ui/loader'
import {
  exchangeIntermediateSessionForSession,
  verifyAuthenticatorApp,
} from '@/gen/openauth/intermediate/v1/intermediate-IntermediateService_connectquery'
import { parseErrorMessage } from '@/lib/errors'
import { cn } from '@/lib/utils'
import { useMutation } from '@connectrpc/connect-query'
import React, { FC, useState } from 'react'
import { useNavigate } from 'react-router'
import { toast } from 'sonner'

const VerifyAuthenticatorAppView: FC = () => {
  const navigate = useNavigate()

  const exchangeIntermediateSessionForSessionMutation = useMutation(
    exchangeIntermediateSessionForSession,
  )
  const verifyAuthenticatorAppMutation = useMutation(verifyAuthenticatorApp)

  const [code, setCode] = useState('')
  const [submitting, setSubmitting] = useState<boolean>(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setSubmitting(true)

    try {
      await verifyAuthenticatorAppMutation.mutateAsync({
        totpCode: code,
      })

      const { accessToken, refreshToken } =
        await exchangeIntermediateSessionForSessionMutation.mutateAsync({})

      setAccessToken(accessToken)
      setRefreshToken(refreshToken)
      setSubmitting(false)

      navigate('/')
    } catch (error) {
      setSubmitting(false)
      const message = parseErrorMessage(error)
      toast.error(message)
    }
  }

  return (
    <>
      <Title title="Verify your time-based one-time password" />

      <Card className="w-full max-w-sm">
        <CardHeader>
          <CardTitle>Verify Authenticator App</CardTitle>
          <CardDescription>
            Enter the 6-digit code from your authenticator app.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form
            className="mt-8 flex flex-col items-center w-full"
            onSubmit={handleSubmit}
          >
            <InputOTP maxLength={6} onChange={(value) => setCode(value)}>
              <InputOTPGroup>
                <InputOTPSlot index={0} />
                <InputOTPSlot index={1} />
                <InputOTPSlot index={2} />
              </InputOTPGroup>
              <InputOTPSeparator />
              <InputOTPGroup>
                <InputOTPSlot index={3} />
                <InputOTPSlot index={4} />
                <InputOTPSlot index={5} />
              </InputOTPGroup>
            </InputOTP>

            <Button
              className="mt-4"
              disabled={code.length < 6 || submitting}
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

export default VerifyAuthenticatorAppView
