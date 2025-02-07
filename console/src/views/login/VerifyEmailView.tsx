import React, { Dispatch, FC, useEffect, useState } from 'react'
import { LoginView } from '@/lib/views'
import { Title } from '@/components/Title'
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { useIntermediateSession } from '@/lib/auth'
import { useMutation } from '@connectrpc/connect-query'
import { verifyEmailChallenge } from '@/gen/openauth/intermediate/v1/intermediate-IntermediateService_connectquery'
import {
  InputOTP,
  InputOTPGroup,
  InputOTPSlot,
} from '@/components/ui/input-otp'

interface VerifyEmailViewProps {
  setView: Dispatch<React.SetStateAction<LoginView>>
}

const VerifyEmailView: FC<VerifyEmailViewProps> = ({ setView }) => {
  const intermediateSession = useIntermediateSession()

  const [challengeCode, setChallengeCode] = useState<string>('')

  const verifyEmailChallengeMutation = useMutation(verifyEmailChallenge)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    try {
      await verifyEmailChallengeMutation.mutateAsync({
        code: challengeCode,
      })

      setView(LoginView.ChooseProject)
    } catch (error) {
      console.error(error)
    }
  }

  return (
    <>
      <Title title="Verify Email Address" />

      <Card className="w-full max-w-sm">
        <CardHeader>
          <CardTitle className="text-center uppercase text-foreground font-semibold text-sm tracking-wide mt-2">
            Verify Email Address
          </CardTitle>
          <p className="text-sm text-center mt-2 text-gray-500">
            Please enter the verification code sent to{' '}
            <b>{intermediateSession?.email}</b> below.
          </p>
        </CardHeader>
        <CardContent className="flex flex-col items-center justify-center w-full">
          <form className="flex flex-col items-center" onSubmit={handleSubmit}>
            <InputOTP maxLength={6} onChange={setChallengeCode}>
              <InputOTPGroup>
                <InputOTPSlot index={0} />
                <InputOTPSlot index={1} />
                <InputOTPSlot index={2} />
                <InputOTPSlot index={3} />
                <InputOTPSlot index={4} />
                <InputOTPSlot index={5} />
              </InputOTPGroup>
            </InputOTP>

            <Button
              className="mt-4"
              disabled={challengeCode.length < 6}
              type="submit"
            >
              Continue
            </Button>
          </form>
        </CardContent>
        <CardFooter></CardFooter>
      </Card>
    </>
  )
}

export default VerifyEmailView
