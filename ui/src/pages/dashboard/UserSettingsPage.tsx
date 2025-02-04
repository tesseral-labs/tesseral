import React, { FC, FormEvent, MouseEvent, useEffect, useState } from 'react'
import { useUser } from '@/lib/auth'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { useMutation, useQuery } from '@connectrpc/connect-query'
import {
  getPasskeyOptions,
  registerPasskey,
  setPassword as setUserPassword,
  whoAmI,
} from '@/gen/openauth/frontend/v1/frontend-FrontendService_connectquery'
import { Input } from '@/components/ui/input'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableRow,
} from '@/components/ui/table'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog'
import { base64urlEncode } from '@/lib/utils'
import { parseErrorMessage } from '@/lib/errors'
import { toast } from 'sonner'

const UserSettingsPage: FC = () => {
  const encoder = new TextEncoder()
  const user = useUser()

  const { data: whoamiRes } = useQuery(whoAmI)
  const setPasswordMutation = useMutation(setUserPassword)
  // const getAuthenticatorAppOptionsMutation = useMutation(getAuthenticatorAppOptions)
  const getPasskeyOptionsMutation = useMutation(getPasskeyOptions)
  // const registerAuthenticatorAppMutation = useMutation(registerAuthenticatorApp)
  const registerPasskeyMutation = useMutation(registerPasskey)

  const [editingEmail, setEditingEmail] = useState(false)
  const [editingPassword, setEditingPassword] = useState(false)
  const [email, setEmail] = useState(whoamiRes?.email || '')
  const [password, setPassword] = useState('')

  const handleEmailSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault()

    // TODO: Kick off email validation and show a modal to verify the new email address
  }

  const handlePasswordSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault()

    try {
      setPasswordMutation.mutateAsync({
        password,
      })
      setEditingPassword(false)
    } catch (error) {
      console.error(error)
    }
  }

  const handleAuthenticatorAppClick = async (
    e: MouseEvent<HTMLButtonElement>,
  ) => {
    e.stopPropagation()

    return true
  }

  const handleRegisterPasskeyClick = async () => {
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
    } catch (error) {
      const message = parseErrorMessage(error)
      toast.error('Could not register passkey', {
        description: message,
      })
    }
  }

  useEffect(() => {
    if (whoamiRes && whoamiRes.email) {
      setEmail(whoamiRes.email || '')
    }
  }, [whoamiRes])

  return (
    <div className="dark:text-foreground">
      <h1 className="text-2xl font-bold mb-4">User Settings</h1>

      <Card>
        <CardHeader>
          <CardTitle>Basic information</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 gap-x-2 text-sm md:grid-cols-2 lg:grid-cols-3">
            <div className="pr-8 dark:border-gray-700 md:border-r">
              <div className="text-sm font-semibold mb-2">User ID</div>
              <div className="text-sm text-gray-500">{whoamiRes?.id}</div>
            </div>
            <div className="pr-8 mt-8 dark:border-gray-700 lg:border-r lg:px-8 md:mt-0">
              <div className="text-sm font-semibold mb-2">Google User ID</div>
              <div className="text-sm text-gray-500">
                {user?.googleUserId || '—'}
              </div>
            </div>
            <div className="pr-8 mt-8 lg:px-8 lg:mt-0">
              <div className="text-sm font-semibold mb-2">
                Microsoft User ID
              </div>
              <div className="text-sm text-gray-500">
                {user?.microsoftUserId || '—'}
              </div>
            </div>
          </div>
        </CardContent>
      </Card>
      <Card className="mt-4 pt-4">
        <CardContent>
          <form onSubmit={handleEmailSubmit}>
            <label className="block w-full font-semibold mb-2">Email</label>
            <Input
              className="max-w-xs"
              disabled={!editingEmail}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="jane.doe@example.com"
              type="email"
              value={email}
            />
            <div className="mt-2">
              {editingEmail ? (
                <>
                  <Button
                    className="text-sm rounded border border-border focus:border-primary mb-2 mr-2"
                    onClick={(e: MouseEvent<HTMLButtonElement>) => {
                      e.preventDefault()
                      e.stopPropagation()
                      setEditingEmail(false)
                    }}
                    variant="outline"
                  >
                    Cancel
                  </Button>
                  <Button
                    className="text-sm rounded border border-border focus:border-primary mb-2"
                    type="submit"
                  >
                    Save Email
                  </Button>
                </>
              ) : (
                <Button
                  className="text-sm rounded border border-border focus:border-primary mb-2"
                  onClick={(e: MouseEvent<HTMLButtonElement>) => {
                    e.preventDefault()
                    e.stopPropagation()
                    setEditingEmail(true)
                  }}
                  variant="outline"
                >
                  Change Email
                </Button>
              )}
            </div>
          </form>
        </CardContent>
      </Card>
      <Card className="mt-4 pt-4">
        <CardContent>
          <form onSubmit={handlePasswordSubmit}>
            <label className="block w-full font-semibold mb-2">Password</label>
            <Input
              className="max-w-xs"
              disabled={!editingPassword}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="•••••••••••••"
              type="password"
              value={password}
            />
            <div className="mt-2">
              {editingPassword ? (
                <>
                  <Button
                    className="text-sm rounded border border-border focus:border-primary mb-2 mr-2"
                    onClick={(e: MouseEvent<HTMLButtonElement>) => {
                      e.preventDefault()
                      e.stopPropagation()

                      setEditingPassword(false)
                    }}
                    variant="outline"
                  >
                    Cancel
                  </Button>
                  <Button
                    className="text-sm rounded border border-border focus:border-primary mb-2"
                    type="submit"
                  >
                    Save Password
                  </Button>
                </>
              ) : (
                <Button
                  className="text-sm rounded border border-border focus:border-primary mb-2"
                  onClick={(e: MouseEvent<HTMLButtonElement>) => {
                    e.preventDefault()
                    e.stopPropagation()
                    setEditingPassword(true)
                  }}
                  variant="outline"
                >
                  Change Password
                </Button>
              )}
            </div>
          </form>
        </CardContent>
      </Card>
      <Card className="mt-4">
        <CardHeader>
          <CardTitle>MFA</CardTitle>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHead>
              <TableRow className="font-bold">
                <TableCell>Method</TableCell>
                <TableCell>Registered</TableCell>
                <TableCell></TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              <TableRow>
                <TableCell>Authenticator App</TableCell>
                <TableCell>Yes</TableCell>
                <TableCell className="text-right">
                  <Dialog>
                    <DialogTrigger asChild>
                      <Button
                        className="text-sm rounded border border-border focus:border-primary"
                        onClick={handleAuthenticatorAppClick}
                        variant="outline"
                      >
                        Register
                      </Button>
                    </DialogTrigger>
                    <DialogContent>
                      <DialogHeader>
                        <DialogTitle>Register Authenticator App</DialogTitle>
                        <DialogDescription></DialogDescription>
                      </DialogHeader>
                    </DialogContent>
                  </Dialog>
                </TableCell>
              </TableRow>
              <TableRow>
                <TableCell>Passkey</TableCell>
                <TableCell>No</TableCell>
                <TableCell className="text-right">
                  <Button
                    className="text-sm rounded border border-border focus:border-primary"
                    onClick={handleRegisterPasskeyClick}
                    variant="outline"
                  >
                    Register
                  </Button>
                </TableCell>
              </TableRow>
            </TableBody>
          </Table>
        </CardContent>
      </Card>
    </div>
  )
}

export default UserSettingsPage
