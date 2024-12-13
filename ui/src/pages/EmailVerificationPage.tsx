import React, { useEffect, useState } from 'react'
import { useNavigate, useSearchParams } from 'react-router-dom'
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

const EmailVerificationPage = () => {
  const navigate = useNavigate()
  const [searchParams] = useSearchParams()

  const [challengeCode, setChallengeCode] = useState<string>('')
  const [challengeId, setChallengeId] = useState<string>('')
  const [email, setEmail] = useState<string>('')

  const { data: whoamiRes } = useQuery(whoami)
  const verifyEmailChallengeMutation = useMutation(verifyEmailChallenge)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    try {
      const response = await verifyEmailChallengeMutation.mutateAsync({
        emailVerificationChallengeId: challengeId,
        code: challengeCode,
      })

      console.log('response: ', response)

      navigate('/organizations')
    } catch (error) {
      console.error(error)
    }
  }

  useEffect(() => {
    ;(async () => {
      try {
        const challengeId = searchParams.get('challenge_id')
        if (challengeId) {
          setChallengeId(challengeId)
        }
      } catch (error) {
        console.error(error)
      }
    })()
  }, [])

  return (
    <>
      <Title title="Verify Email Address" />

      <Card className="w-[clamp(320px,50%,420px)]">
        <CardHeader>
          <CardTitle className="text-center uppercase text-foreground font-semibold text-sm tracking-wide mt-2">
            Verify Email Address
          </CardTitle>
        </CardHeader>
        <CardContent className="flex flex-col items-center justify-center w-full">
          <p className="text-center mb-3">
            Please enter the verification code sent to <b>{whoamiRes?.email}</b>{' '}
            below.
          </p>
          <form className="flex flex-col items-center" onSubmit={handleSubmit}>
            <input
              className="text-sm rounded border border-border focus:border-primary w-[clamp(240px,50%,100%)] mb-2"
              id="challengeCode"
              placeholder="Enter your challenge code"
              value={challengeCode}
              onChange={(e) => setChallengeCode(e.target.value)}
            />
            <Button
              className="text-sm rounded border border-border focus:border-primary w-[clamp(240px,50%,100%)] mb-2"
              type="submit"
            >
              Submit
            </Button>
          </form>
        </CardContent>
        <CardFooter></CardFooter>
      </Card>
    </>
  )
}

export default EmailVerificationPage
