import React, { FC, useEffect } from 'react'
import { useNavigate } from 'react-router'
import { useMutation } from '@connectrpc/connect-query'
import { toast } from 'sonner'

import { setAccessToken, setRefreshToken } from '@/auth'
import { Title } from '@/components/Title'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import {
  exchangeIntermediateSessionForSession,
  getPasskeyOptions,
  registerPasskey,
} from '@/gen/openauth/intermediate/v1/intermediate-IntermediateService_connectquery'
import { parseErrorMessage } from '@/lib/errors'
import { base64urlEncode, cn } from '@/lib/utils'
import { useIntermediateOrganization } from '@/lib/auth'

const RegisterPasskeyView: FC = () => {
  const encoder = new TextEncoder()
  const navigate = useNavigate()
  const organization = useIntermediateOrganization()

  const exchangeIntermediateSessionForSessionMutation = useMutation(
    exchangeIntermediateSessionForSession,
  )
  const getPasskeyOptionsMutation = useMutation(getPasskeyOptions)
  const registerPasskeyMutation = useMutation(registerPasskey)

  const registerCredential = async (): Promise<any> => {
    try {
      if (!navigator.credentials) {
        throw new Error('WebAuthn not supported')
      }

      const passkeyOptions = await getPasskeyOptionsMutation.mutateAsync({})
      const credentialOptions: PublicKeyCredentialCreationOptions = {
        challenge: new Uint8Array([0]).buffer,
        rp: {
          id: passkeyOptions.rpId,
          name: passkeyOptions.rpName,
        },
        user: {
          id: encoder.encode(passkeyOptions.userId).buffer,
          name: passkeyOptions.userDisplayName,
          displayName: passkeyOptions.userDisplayName,
        },
        pubKeyCredParams: [
          { type: 'public-key', alg: -7 }, // ECDSA with SHA-256
          { type: 'public-key', alg: -257 }, // RSA with SHA-256
        ],
        timeout: 60000,
        attestation: 'direct',
      }

      const credential = (await navigator.credentials.create({
        publicKey: credentialOptions,
      })) as PublicKeyCredential

      if (!credential) {
        throw new Error('No credential returned')
      }

      await registerPasskeyMutation.mutateAsync({
        attestationObject: base64urlEncode(
          (credential.response as AuthenticatorAttestationResponse)
            .attestationObject,
        ),
      })

      const { accessToken, refreshToken } =
        await exchangeIntermediateSessionForSessionMutation.mutateAsync({})

      setAccessToken(accessToken)
      setRefreshToken(refreshToken)

      navigate('/settings')
    } catch (error) {
      const message = parseErrorMessage(error)
      toast.error('Could not register passkey', {
        description: message,
      })
    }
  }

  useEffect(() => {
    ;(async () => {
      const credential = await registerCredential()
      console.log(credential)
    })()
  }, [])

  return (
    <>
      <Title title="Register a Passkey" />

      <Card className="w-full max-w-sm">
        <CardHeader>
          <CardTitle className="text-center">Register a Passkey</CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-center text-sm text-muted-foreground">
            Follow the prompts on your device to register a new Passkey.
          </p>
        </CardContent>
      </Card>
    </>
  )
}

export default RegisterPasskeyView
