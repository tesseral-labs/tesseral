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
import React, { FC, useState } from 'react'

const VerifyAuthenticatorApp: FC = () => {
  const [code, setCode] = useState('')

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    console.log('submitting', code)
  }

  return (
    <>
      <Title title="Register your time-based one-time password" />
      <Card>
        <CardHeader>
          <CardTitle>Register Authenticator App</CardTitle>
        </CardHeader>
        <CardContent>
          <p className="mt-4 text-sm text-center text-muted-foreground">
            Enter the 6-digit code from your authenticator app and
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
