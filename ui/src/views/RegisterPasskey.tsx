import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import {
  getPasskeyOptions,
  registerPasskey,
} from '@/gen/openauth/intermediate/v1/intermediate-IntermediateService_connectquery'
import { base64urlEncode } from '@/lib/utils'
import { useMutation } from '@connectrpc/connect-query'
import React, { FC, useEffect } from 'react'
import { useNavigate } from 'react-router'

const encoder = new TextEncoder()

const RegisterPasskey: FC = () => {
  const navigate = useNavigate()

  const getPasskeyOptionsMutation = useMutation(getPasskeyOptions)
  const registerPasskeyMutation = useMutation(registerPasskey)

  const registerCredential = async (): Promise<any> => {
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

    const response = await registerPasskeyMutation.mutateAsync({
      attestationObject: base64urlEncode(
        (credential.response as AuthenticatorAttestationResponse)
          .attestationObject,
      ),
    })

    navigate('/settings')
  }

  useEffect(() => {
    ;(async () => {
      const credential = await registerCredential()
      console.log(credential)
    })()
  }, [])

  return (
    <Card>
      <CardHeader>
        <CardTitle>Register a Passkey</CardTitle>
      </CardHeader>
      <CardContent>
        <p>Follow the prompts on your device to register a new Passkey.</p>
      </CardContent>
    </Card>
  )
}

export default RegisterPasskey
