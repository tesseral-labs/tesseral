import React, { useEffect } from 'react'
import { Title } from '@/components/Title'
import { useNavigate, useSearchParams } from 'react-router-dom'
import { redeemMicrosoftOAuthCode } from '@/gen/openauth/intermediate/v1/intermediate-IntermediateService_connectquery'
import { useMutation } from '@connectrpc/connect-query'

const MicrosoftOAuthCallbackPage = () => {
  const navigate = useNavigate()
  const [searchParams] = useSearchParams()

  const redeemMicrosoftOAuthCodeMutation = useMutation(redeemMicrosoftOAuthCode)

  useEffect(() => {
    ;(async () => {
      const code = searchParams.get('code')
      const state = searchParams.get('state')

      if (code && state) {
        try {
          const { emailVerificationChallengeId, shouldVerifyEmail } =
            await redeemMicrosoftOAuthCodeMutation.mutateAsync({
              code,
              state,
            })

          if (shouldVerifyEmail && emailVerificationChallengeId) {
            navigate(
              `/verify-email?challenge_id=${emailVerificationChallengeId}`,
            )
            return
          }

          navigate('/organizations')
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
