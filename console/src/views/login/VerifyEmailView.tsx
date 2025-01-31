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
import { useLocation } from 'react-router'
import { useIntermediateSession } from '@/lib/auth'
import { useMutation } from '@connectrpc/connect-query'
import { verifyEmailChallenge } from '@/gen/openauth/intermediate/v1/intermediate-IntermediateService_connectquery'
import { Input } from '@/components/ui/input'

interface VerifyEmailViewProps {
  setView: Dispatch<React.SetStateAction<LoginView>>
}

const VerifyEmailView: FC<VerifyEmailViewProps> = ({ setView }) => {
  const intermediateSession = useIntermediateSession()
  const { state } = useLocation()

  const [challengeCode, setChallengeCode] = useState<string>('')
  const [email, setEmail] = useState<string>('')

  const verifyEmailChallengeMutation = useMutation(verifyEmailChallenge)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    try {
      await verifyEmailChallengeMutation.mutateAsync({
        emailVerificationChallengeId: state?.challengeId,
        code: challengeCode,
      })

      setView(LoginView.ChooseProject)
    } catch (error) {
      console.error(error)
    }
  }

  // useEffect(() => {
  //   if (intermediateSession !== null && !intermediateSession) {
  //     setView(LoginView.StartLogin)
  //   }
  // }, [intermediateSession])

  return (
    <>
      <Title title="Verify Email Address" />

      <Card className="w-[clamp(320px,50%,420px)] mx-auto">
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
            <Input
              className="text-sm bg-input rounded border border-border focus:border-primary w-[clamp(240px,50%,100%)] mb-2"
              id="challengeCode"
              placeholder="Enter your challenge code"
              value={challengeCode}
              onChange={(e) => setChallengeCode(e.target.value)}
            />
            <Button
              className="text-sm rounded border border-border focus:border-primary w-[clamp(240px,50%,100%)] mb-2"
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

export default VerifyEmailView
