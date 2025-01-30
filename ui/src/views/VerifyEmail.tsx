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

      setView(LoginViews.Organizations)
    } catch (error) {
      console.error(error)
    }
  }

  return (
    <>
      <Title title="Verify Email Address" />

      <Card className="w-[clamp(320px,50%,420px)]">
        <CardHeader>
          <CardTitle className="text-center uppercase text-foreground font-semibold text-sm tracking-wide mt-2">
            Verify Email Address
          </CardTitle>
          <p className="text-sm text-center mt-2 text-gray-500">
            Please enter the verification code sent to{' '}
            <b>{whoamiRes?.intermediateSession?.email}</b> below.
          </p>
        </CardHeader>
        <CardContent className="flex flex-col items-center justify-center w-full">
          <form className="flex flex-col items-center" onSubmit={handleSubmit}>
            <Input
              className="w-[clamp(240px,50%,100%)] mb-2"
              id="challengeCode"
              placeholder="Enter your challenge code"
              value={challengeCode}
              onChange={(e) => setChallengeCode(e.target.value)}
            />
            <Button className="w-[clamp(240px,50%,100%)]" type="submit">
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
