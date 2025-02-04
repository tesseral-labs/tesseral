import React, { FC, FormEvent, MouseEvent, useEffect, useState } from 'react'
import QRCode from 'qrcode'
import { useUser } from '@/lib/auth'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { useMutation, useQuery } from '@connectrpc/connect-query'
import {
  deleteMyPasskey,
  getAuthenticatorAppOptions,
  getPasskeyOptions,
  listMyPasskeys,
  registerAuthenticatorApp,
  registerPasskey,
  setPassword as setUserPassword,
  whoami,
} from '@/gen/openauth/frontend/v1/frontend-FrontendService_connectquery'
import { Input } from '@/components/ui/input'
import {
  Table,
  TableBody,
  TableCell,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog'
import { base64urlEncode } from '@/lib/utils'
import { parseErrorMessage } from '@/lib/errors'
import { toast } from 'sonner'
import {
  InputOTP,
  InputOTPGroup,
  InputOTPSeparator,
  InputOTPSlot,
} from '@/components/ui/input-otp'
import Loader from '@/components/ui/loader'
import { CheckCircle, PlusCircle } from 'lucide-react'
import { set } from 'react-hook-form'

const UserSettingsPage: FC = () => {
  const encoder = new TextEncoder()
  const user = useUser()

  const { data: whoamiRes } = useQuery(whoami)
  const deleteMyPasskeyMutation = useMutation(deleteMyPasskey)
  const setPasswordMutation = useMutation(setUserPassword)
  const getAuthenticatorAppOptionsMutation = useMutation(
    getAuthenticatorAppOptions,
  )
  const getPasskeyOptionsMutation = useMutation(getPasskeyOptions)
  const { data: listMyPasskeysRes, refetch: refetchMyPasskeys } =
    useQuery(listMyPasskeys)
  const registerAuthenticatorAppMutation = useMutation(registerAuthenticatorApp)
  const registerPasskeyMutation = useMutation(registerPasskey)

  const [authenticatorAppCode, setAuthenticatorAppCode] = useState('')
  const [authenticatorAppDialogOpen, setAuthenticatorAppDialogOpen] =
    useState(false)
  const [editingEmail, setEditingEmail] = useState(false)
  const [editingPassword, setEditingPassword] = useState(false)
  const [email, setEmail] = useState(whoamiRes?.user?.email || '')
  const [password, setPassword] = useState('')
  const [qrImage, setQRImage] = useState<string | null>(null)
  const [registeringAuthenticatorApp, setRegisteringAuthenticatorApp] =
    useState(false)

  const handleDeletePasskey = async (id: string) => {
    try {
      await deleteMyPasskeyMutation.mutateAsync({
        id,
      })
    } catch (error) {
      const message = parseErrorMessage(error)
      toast.error('Could not delete passkey', {
        description: message,
      })
    }

    try {
      await refetchMyPasskeys()
    } catch (error) {
      const message = parseErrorMessage(error)
      toast.error('Could not fetch updated passkey list', {
        description: message,
      })
    }
  }

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

    const authenticatorAppOptions =
      await getAuthenticatorAppOptionsMutation.mutateAsync({})
    const qrImage = await QRCode.toDataURL(authenticatorAppOptions.otpauthUri, {
      errorCorrectionLevel: 'H',
    })
    setQRImage(qrImage)

    return true
  }

  const handleRegisterAuthenticatorApp = async (
    e: FormEvent<HTMLFormElement>,
  ) => {
    e.preventDefault()
    setRegisteringAuthenticatorApp(true)

    try {
      await registerAuthenticatorAppMutation.mutateAsync({
        totpCode: authenticatorAppCode,
      })
      setRegisteringAuthenticatorApp(false)
      setAuthenticatorAppDialogOpen(false)
      setAuthenticatorAppCode('')
      setQRImage(null)
    } catch (error) {
      setRegisteringAuthenticatorApp(false)
      const message = parseErrorMessage(error)
      toast.error('Could not register authenticator app', {
        description: message,
      })
    }
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
    if (whoamiRes && whoamiRes.user?.email) {
      setEmail(whoamiRes.user.email)
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
              <div className="text-sm text-gray-500">{whoamiRes?.user?.id}</div>
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
          <CardTitle>
            <div className="grid grid-cols-2 gap-4">
              <span>Authenticator App</span>
              <div className="flex flex-row items-end justify-end">
                <Dialog
                  open={authenticatorAppDialogOpen}
                  onOpenChange={setAuthenticatorAppDialogOpen}
                >
                  <DialogTrigger asChild>
                    <Button
                      onClick={handleAuthenticatorAppClick}
                      variant="outline"
                    >
                      <PlusCircle />
                      Register Authenticator App
                    </Button>
                  </DialogTrigger>
                  <DialogContent>
                    <DialogHeader>
                      <DialogTitle>Register Authenticator App</DialogTitle>
                      <DialogDescription>
                        Scan the QR code with your authenticator app to
                        register.
                      </DialogDescription>

                      <div className="flex flex-row justify-center w-full">
                        {qrImage ? (
                          <div className="border rounded-lg w-full">
                            <img className="w-full" src={qrImage} />
                          </div>
                        ) : (
                          <div className="my-8">
                            <Loader />
                          </div>
                        )}
                      </div>

                      <form
                        className="mt-8 flex flex-col items-center w-full"
                        onSubmit={handleRegisterAuthenticatorApp}
                      >
                        <InputOTP
                          maxLength={6}
                          onChange={(value) => setAuthenticatorAppCode(value)}
                        >
                          <InputOTPGroup>
                            <InputOTPSlot index={0} />
                            <InputOTPSlot index={1} />
                            <InputOTPSlot index={2} />
                          </InputOTPGroup>
                          <InputOTPSeparator />
                          <InputOTPGroup>
                            <InputOTPSlot index={3} />
                            <InputOTPSlot index={4} />
                            <InputOTPSlot index={5} />
                          </InputOTPGroup>
                        </InputOTP>

                        <Button
                          className="mt-4"
                          disabled={registeringAuthenticatorApp}
                          type="submit"
                        >
                          {registeringAuthenticatorApp && <Loader />}
                          Submit
                        </Button>
                      </form>
                    </DialogHeader>
                  </DialogContent>
                </Dialog>
              </div>
            </div>
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="text-sm flex flex-row items-center">
            {whoamiRes?.user?.hasAuthenticatorApp ? (
              <>
                <CheckCircle />
                <span className="ml-2">Registered</span>
              </>
            ) : (
              'Not registered'
            )}
          </div>
        </CardContent>
      </Card>
      <Card className="mt-4">
        <CardHeader>
          <CardTitle>
            <div className="grid grid-cols-2 gap-4">
              <span>Passkeys</span>
              <div className="flex flex-row items-end justify-end">
                <Button onClick={handleRegisterPasskeyClick} variant="outline">
                  <PlusCircle />
                  Register Passkey
                </Button>
              </div>
            </div>
          </CardTitle>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow className="font-bold">
                <TableCell>ID</TableCell>
                <TableCell>Type</TableCell>
                <TableCell className="flex flex-col items-end"></TableCell>
              </TableRow>
            </TableHeader>
            <TableBody>
              {listMyPasskeysRes?.passkeys.map((passkey) => (
                <TableRow key={passkey.id}>
                  <TableCell className="text-sm">{passkey.id}</TableCell>
                  <TableCell>Passkey</TableCell>
                  <TableCell className="text-right">
                    <Dialog>
                      <DialogTrigger asChild>
                        <Button variant="destructive">Delete</Button>
                      </DialogTrigger>
                      <DialogContent>
                        <DialogHeader>
                          <DialogTitle>Are you sure?</DialogTitle>
                          <DialogDescription>
                            Once deleted, you'll no longer be able to log in
                            with this passkey.
                          </DialogDescription>
                        </DialogHeader>
                        <DialogFooter>
                          <Button
                            className="mr-2"
                            onClick={(e: MouseEvent<HTMLButtonElement>) => {
                              e.preventDefault()
                              handleDeletePasskey(passkey.id)
                            }}
                            variant="destructive"
                          >
                            Delete
                          </Button>
                        </DialogFooter>
                      </DialogContent>
                    </Dialog>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </CardContent>
      </Card>
    </div>
  )
}

export default UserSettingsPage
