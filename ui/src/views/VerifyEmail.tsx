import React, { Dispatch, FC, SetStateAction, useState } from 'react'
import { useMutation, useQuery } from '@connectrpc/connect-query'

import { Title } from '@/components/Title'
import { Button } from '@/components/ui/button'
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import {
  verifyEmailChallenge,
  whoami,
} from '@/gen/openauth/intermediate/v1/intermediate-IntermediateService_connectquery'
import { LoginViews } from '@/lib/views'
import { Input } from '@/components/ui/input'
import {
  InputOTP,
  InputOTPGroup,
  InputOTPSeparator,
  InputOTPSlot,
} from '@/components/ui/input-otp'

interface VerifyEmailProps {
  setView: Dispatch<SetStateAction<LoginViews>>
}

const VerifyEmail: FC<VerifyEmailProps> = ({ setView }) => {
  const [challengeCode, setChallengeCode] = useState<string>('')
  const [email, setEmail] = useState<string>('')

  const { data: whoamiRes } = useQuery(whoami)
  const verifyEmailChallengeMutation = useMutation(verifyEmailChallenge)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    try {
      await verifyEmailChallengeMutation.mutateAsync({
        code: challengeCode,
      })

      setView(LoginViews.ChooseOrganization)
    } catch (error) {
      console.error(error)
    }
  }

  return (
    <>
      <Title title="Verify Email Address" />

      <Card className="max-w-sm">
        <CardHeader>
          <CardTitle className="text-center">Verify Email Address</CardTitle>

          <p className="text-sm mt-2 text-muted-foreground text-center">
            Please enter the verification code sent to{' '}
            <b>{whoamiRes?.intermediateSession?.email}</b> below.
          </p>
        </CardHeader>
        <CardContent className="flex flex-col items-center justify-center w-full">
          <form className="flex flex-col items-center" onSubmit={handleSubmit}>
            <InputOTP
              maxLength={6}
              id="challengeCode"
              value={challengeCode}
              onChange={(value) => setChallengeCode(value)}
            >
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
              className="w-full mt-4"
              disabled={challengeCode.length < 6}
              type="submit"
            >
              Verify Email Address
            </Button>
          </form>
        </CardContent>
        <CardFooter></CardFooter>
      </Card>
    </>
  )
}

export default VerifyEmail
