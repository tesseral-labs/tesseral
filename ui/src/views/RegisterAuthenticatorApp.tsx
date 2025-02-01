import React, { FC, useEffect, useState } from 'react'
import QRCode from 'qrcode'
import { Title } from '@/components/Title'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { useMutation, useQuery } from '@connectrpc/connect-query'
import {
  exchangeIntermediateSessionForSession,
  getAuthenticatorAppOptions,
  registerAuthenticatorApp,
  whoami,
} from '@/gen/openauth/intermediate/v1/intermediate-IntermediateService_connectquery'
import { base32Encode } from '@/lib/utils'
import { useIntermediateOrganization } from '@/lib/auth'
import {
  InputOTP,
  InputOTPGroup,
  InputOTPSeparator,
  InputOTPSlot,
} from '@/components/ui/input-otp'
import { setAccessToken, setRefreshToken } from '@/auth'
import { useNavigate } from 'react-router'

const RegisterAuthenticatorApp: FC = () => {
  const navigate = useNavigate()
  const organization = useIntermediateOrganization()
  const [qrcode, setQRCode] = useState<string | null>(null)
  const [code, setCode] = useState<string>('')

  const exchangeIntermediateSessionForSessionMutation = useMutation(
    exchangeIntermediateSessionForSession,
  )
  const getAuthenticatorAppOptionsMutation = useMutation(
    getAuthenticatorAppOptions,
  )
  const registerAuthenticatorAppMutation = useMutation(registerAuthenticatorApp)
  const { data: whoamiRes } = useQuery(whoami)

  const generateQRCode = async (): Promise<string> => {
    const authenticatorAppOptions =
      await getAuthenticatorAppOptionsMutation.mutateAsync({})
    const secret = base32Encode(authenticatorAppOptions.secret)
    const url = `otpauth://totp/${organization?.displayName}:${whoamiRes?.intermediateSession?.email}?secret=${secret}&issuer=${organization?.displayName}`

    return QRCode.toDataURL(url, {
      errorCorrectionLevel: 'H',
    })
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    try {
      await registerAuthenticatorAppMutation.mutateAsync({
        totpCode: code,
      })

      const { accessToken, refreshToken } =
        await exchangeIntermediateSessionForSessionMutation.mutateAsync({})

      setAccessToken(accessToken)
      setRefreshToken(refreshToken)

      navigate('/settings')
    } catch (error) {
      console.error(error)
    }
  }

  useEffect(() => {
    ;(async () => {
      const qrcode = await generateQRCode()
      setQRCode(qrcode)
    })()
  }, [])

  return (
    <>
      <Title title="Register your time-based one-time password" />
      <Card className="max-w-sm">
        <CardHeader>
          <CardTitle>Register Authenticator App</CardTitle>
        </CardHeader>
        <CardContent>
          {qrcode && (
            <div className="border rounded-lg w-full mr-auto">
              <img className="w-full" src={qrcode} />
            </div>
          )}

          <p className="mt-4 text-sm text-center">
            Scan this QR code using your authenticator app and enter the
            resulting 6-digit code.
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

export default RegisterAuthenticatorApp
