import React, { FC, useEffect } from 'react'

const encoder = new TextEncoder()

const RegisterPasskey: FC = () => {
  const getPasskeyOptionsMutation = useMutation(getPasskeyOptions)
  // const registerPasskeyMutation = useMutation(registerPasskey)

  const registerPasskey = async (): Promise<any> => {
    if (!navigator.credentials) {
      throw new Error('WebAuthn not supported')
    }

    const { data: passkeyOptions } =
      await getPasskeyOptionsMutation.mutateAsync({})
    const credentialOptions: PublicKeyCredentialCreationOptions = {
      challenge: encoder.encode(passkeyOptions.challenge).buffer,
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

    const credential = await navigator.credentials.create({
      publicKey: credentialOptions,
    })

    if (!credential) {
      throw new Error('No credential returned')
    }

    return credential
  }

  useEffect(() => {
    ;(async () => {
      const credential = await registerPasskey()
      console.log(credential)
    })()
  }, [])

  return <div>Register webauthn</div>
}

export default RegisterPasskey
