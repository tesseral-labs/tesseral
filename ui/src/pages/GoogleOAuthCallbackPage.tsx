import React, { useEffect } from 'react'
import { Title } from '@/components/Title'
import { useNavigate, useSearchParams } from 'react-router-dom'
import { redeemGoogleOAuthCode } from '@/gen/openauth/intermediate/v1/intermediate-IntermediateService_connectquery'
import { useMutation } from '@connectrpc/connect-query'

const GoogleOAuthCallbackPage = () => {
  const navigate = useNavigate()
  const [searchParams] = useSearchParams()

  const redeemGoogleOAuthCodeMutation = useMutation(redeemGoogleOAuthCode)

  useEffect(() => {
    ;(async () => {
      const code = searchParams.get('code')
      const state = searchParams.get('state')

      if (code && state) {
        try {
          await redeemGoogleOAuthCodeMutation.mutateAsync({
            code,
            state,
          })

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
      <Title title="Verifying Google OAuth Credentials..." />

      <div className="space-y-4 text-center"></div>
    </div>
  )
}

export default GoogleOAuthCallbackPage
