import { setAccessToken, setRefreshToken } from '@/auth'
import { Title } from '@/components/Title'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import {
  InputOTP,
  InputOTPGroup,
  InputOTPSeparator,
  InputOTPSlot,
} from '@/components/ui/input-otp'
import {
  exchangeIntermediateSessionForSession,
  verifyAuthenticatorApp,
} from '@/gen/openauth/intermediate/v1/intermediate-IntermediateService_connectquery'
import { parseErrorMessage } from '@/lib/errors'
import { useLayout } from '@/lib/settings'
import { cn } from '@/lib/utils'
import { LoginLayouts } from '@/lib/views'
import { useMutation } from '@connectrpc/connect-query'
import React, { FC, useState } from 'react'
import { useNavigate } from 'react-router'
import { toast } from 'sonner'

const VerifyAuthenticatorApp: FC = () => {
  const layout = useLayout()
  const navigate = useNavigate()

  const exchangeIntermediateSessionForSessionMutation = useMutation(
    exchangeIntermediateSessionForSession,
  )
  const verifyAuthenticatorAppMutation = useMutation(verifyAuthenticatorApp)

  const [code, setCode] = useState('')

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    try {
      await verifyAuthenticatorAppMutation.mutateAsync({
        totpCode: code,
      })

      const { accessToken, refreshToken } =
        await exchangeIntermediateSessionForSessionMutation.mutateAsync({})

      setAccessToken(accessToken)
      setRefreshToken(refreshToken)

      navigate('/settings')
    } catch (error) {
      const message = parseErrorMessage(error)

      toast.error(message)
    }
  }

  return (
    <>
      <Title title="Verify your time-based one-time password" />

      <Card
        className={cn(
          'w-full max-w-sm',
          layout !== LoginLayouts.Centered && 'shadow-none border-0',
        )}
      >
        <CardHeader>
          <CardTitle className="text-center">
            Verify Authenticator App
          </CardTitle>
        </CardHeader>
        <CardContent>
          <p className="mt-4 text-sm text-center text-muted-foreground">
            Enter the 6-digit code from your authenticator app
          </p>

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

            <Button className="mt-4" type="submit">
              Submit
            </Button>
          </form>
        </CardContent>
      </Card>
    </>
  )
}

export default VerifyAuthenticatorApp
