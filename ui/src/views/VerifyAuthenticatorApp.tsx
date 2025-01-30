import { Title } from '@/components/Title'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
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
          <p className="mt-4 w-[300px] text-sm text-center">
            Enter the 6-digit code from your authenticator app and
          </p>

          <form className="mt-8" onSubmit={handleSubmit}>
            <Input
              className="mb-2"
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

export default VerifyAuthenticatorApp
