import React, { useEffect } from 'react'
import { Title } from '@/components/Title'
import { useNavigate, useSearchParams } from 'react-router-dom'
import {
  issueEmailVerificationChallenge,
  redeemGoogleOAuthCode,
  whoami,
} from '@/gen/openauth/intermediate/v1/intermediate-IntermediateService_connectquery'
import { useMutation, useQuery } from '@connectrpc/connect-query'

const GoogleOAuthCallbackPage = () => {
  const navigate = useNavigate()

  const [searchParams] = useSearchParams()

  const issueEmailVerificationChallengeMutation = useMutation(
    issueEmailVerificationChallenge,
  )
  const redeemGoogleOAuthCodeMutation = useMutation(redeemGoogleOAuthCode)
  const whoamiQuery = useQuery(whoami)

  useEffect(() => {
    ;(async () => {
      const code = searchParams.get('code')
      const state = searchParams.get('state')

      if (code && state) {
        try {
          // Redeem the Google OAuth code.
          await redeemGoogleOAuthCodeMutation.mutateAsync({
            code,
            state,
          })

          // Fetch the whoami query to determine if the user has verified their email.
          const { data } = await whoamiQuery.refetch()
          if (!data) {
            throw new Error('No data returned from whoami query')
          }

          // If the user has verified their email, navigate to the organizations page.
          if (data.isEmailVerified) {
            navigate('/organizations')
            return
          }

          // Issue a new email verification challenge.
          const emailVerificationChallengeResponse =
            await issueEmailVerificationChallengeMutation.mutateAsync({})
          if (!emailVerificationChallengeResponse) {
            throw new Error(
              'No data returned from issueEmailVerificationChallenge query',
            )
          }

          // Navigate to the email verification page.
          navigate(
            `/verify-email?challenge_id=${emailVerificationChallengeResponse.emailVerificationChallengeId}`,
          )
        } catch (error) {
          // TODO: Handle errors on screen once an error handling strategy is in place.
          console.error(error)
        }
      }
    })()
  }, [])

  return (
    <div className="flex flex-col items-center justify-center h-screen">
      <Title title="Verifying Google OAuth Credentials..." />

      <div className="space-y-4 text-center"></div>
    </div>
  )
}

export default GoogleOAuthCallbackPage
