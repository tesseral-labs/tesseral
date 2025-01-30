import React, { FC, useEffect } from 'react'

const encoder = new TextEncoder()

const exampleCredentialOptions: PublicKeyCredentialCreationOptions = {
  challenge: encoder.encode('random-base64-url-string').buffer,
  rp: { name: 'Tesseral Localhost' },
  user: {
    id: encoder.encode('random-user-id-in-base64url').buffer,
    name: 'root@app.tesseral.example.com',
    displayName: 'Tesseral User',
  },
  pubKeyCredParams: [
    { type: 'public-key', alg: -7 }, // ECDSA with SHA-256
    { type: 'public-key', alg: -257 }, // RSA with SHA-256
  ],
  timeout: 60000,
  attestation: 'direct',
}

const RegisterPasskey: FC = () => {
  // const RegisterPasskeyMutation = useMutation()

  const RegisterPasskey = async (
    credentialOptions: PublicKeyCredentialCreationOptions,
  ): Promise<any> => {
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
      const credential = await RegisterPasskey(exampleCredentialOptions)
      console.log(credential)
    })()
  }, [])

  return <div>Register webauthn</div>
}

export default RegisterPasskey
