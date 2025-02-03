import React, { FC, useEffect } from 'react'
import { useNavigate } from 'react-router'
import { useMutation } from '@connectrpc/connect-query'
import { toast } from 'sonner'

import { useLayout } from '@/lib/settings'
import { cn } from '@/lib/utils'
import { LoginLayouts } from '@/lib/views'
import { parseErrorMessage } from '@/lib/errors'

import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'

const VerifyPasskey: FC = () => {
  const encoder = new TextEncoder()
  const layout = useLayout()
  const navigate = useNavigate()

  // const issuePasskeyChallengeMutation = useMutation(issuePasskeyChallenge)
  // const verifyPasskeyMutation = useMutation(verifyPasskeyChallenge)

  const authenticateWithPasskey = async () => {
    try {
      // const challengeResponse = await issuePasskeyChallengeMutation.mutateAsync({})
      // const requestOptions: PublicKeyCredentialRequestOptions = {
      //   challenge: encoder.encode(challengeResponse.challenge).buffer,
      //   allowCredentials: [],
      //   userVerification: 'preferred',
      //   rpId: challengeResponse.rpId,
      //   timeout: 60000,
      // }
      // const credential = await navigator.credentials.get({
      //   publicKey: requestOptions,
      // })
      // console.log(credential)
      // await verifyPasskeyMutation.mutateAsync({
      // })
      // const { accessToken, refreshToken } = await exchangeIntermediateSessionForSessionMutation.mutateAsync({})
      // navigate('/settings')
    } catch (error) {
      const message = parseErrorMessage(error)

      toast.error(message)
    }
  }

  useEffect(() => {}, [])

  return (
    <Card
      className={cn(
        'w-full max-w-sm',
        layout !== LoginLayouts.Centered && 'shadow-none border-0',
      )}
    >
      <CardHeader>
        <CardTitle className="text-center">Verify Passkey</CardTitle>
      </CardHeader>
      <CardContent>
        <p className="text-center text-sm text-muted-foreground">
          Follow the prompts on your device to continue logging in with your
          Passkey.
        </p>
      </CardContent>
    </Card>
  )
}

export default VerifyPasskey
