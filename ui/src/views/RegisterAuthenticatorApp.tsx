import React, { FC, useEffect, useState } from 'react'
import QRCode from 'qrcode'
import { Title } from '@/components/Title'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'

const RegisterAuthenticatorApp: FC = () => {
  const secretValue = 'testValue'
  const [qrcode, setQRCode] = useState<string | null>(null)
  const [code, setCode] = useState<string>('')

  const generateQRCode = async (value: string): Promise<string> => {
    return QRCode.toDataURL(value, {
      errorCorrectionLevel: 'H',
    })
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    console.log('submitting', code)
  }

  useEffect(() => {
    ;(async () => {
      const qrcode = await generateQRCode(secretValue)
      setQRCode(qrcode)
    })()
  }, [])

  return (
    <>
      <Title title="Register your time-based one-time password" />
      <Card>
        <CardHeader>
          <CardTitle>Register Authenticator App</CardTitle>
        </CardHeader>
        <CardContent>
          {qrcode && (
            <div className="border rounded-lg w-[300px] mr-auto">
              <img className="w-full" src={qrcode} />
            </div>
          )}

          <p className="mt-4 w-[300px] text-sm text-center">
            Scan this QR code using your authenticator app and enter the
            resulting 6-digit code.
          </p>

          <form className="mt-8" onSubmit={handleSubmit}>
            <Input
              onChange={(e) => setCode(e.target.value)}
              placeholder="6-digit code"
              value={code}
            />

            <Button className="mt-4 w-full" type="submit">
              Submit
            </Button>
          </form>
        </CardContent>
      </Card>
    </>
  )
}

export default RegisterAuthenticatorApp
