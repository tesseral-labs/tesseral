import React, { useEffect } from 'react'
import { Title } from '@/components/Title'
import { useNavigate, useSearchParams } from 'react-router-dom'
import {
  issueEmailVerificationChallenge,
  redeemMicrosoftOAuthCode,
  whoami,
} from '@/gen/openauth/intermediate/v1/intermediate-IntermediateService_connectquery'
import { useMutation, useQuery } from '@connectrpc/connect-query'

const MicrosoftOAuthCallbackPage = () => {
  const navigate = useNavigate()
  const [searchParams] = useSearchParams()

  const issueEmailVerificationChallengeMutation = useMutation(
    issueEmailVerificationChallenge,
  )
  const redeemMicrosoftOAuthCodeMutation = useMutation(redeemMicrosoftOAuthCode)
  const whoamiQuery = useQuery(whoami)

  useEffect(() => {
    ;(async () => {
      const code = searchParams.get('code')
      const state = searchParams.get('state')

      if (code && state) {
        try {
          await redeemMicrosoftOAuthCodeMutation.mutateAsync({
            code,
            state,
            redirectUrl: `${window.location.origin}/microsoft-oauth-callback`,
          })

          const { data } = await whoamiQuery.refetch()
          if (!data) {
            throw new Error('No data returned from whoami query')
          }

          if (data.isEmailVerified) {
            navigate('/organizations')
            return
          }

          const { emailVerificationChallengeId } =
            await issueEmailVerificationChallengeMutation.mutateAsync({})

          navigate(`/verify-email?challenge_id=${emailVerificationChallengeId}`)
        } catch (error) {
          // TODO: Handle errors on screen once an error handling strategy is in place.
          console.error(error)
        }
      }
    })()
  }, [])

  return (
    <div className="flex flex-col items-center justify-center h-screen">
      <Title title="Verifying Microsoft OAuth Credentials..." />

      <div className="space-y-4 text-center"></div>
    </div>
  )
}

export default MicrosoftOAuthCallbackPage
