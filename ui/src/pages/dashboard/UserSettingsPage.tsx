import React, { FC, FormEvent, MouseEvent, useEffect, useState } from 'react'
import { useUser } from '@/lib/auth'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { useMutation } from '@connectrpc/connect-query'
import { setPassword as setUserPassword } from '@/gen/openauth/frontend/v1/frontend-FrontendService_connectquery'
import { Input } from '@/components/ui/input'

const UserSettingsPage: FC = () => {
  const user = useUser()
  const setPasswordMutation = useMutation(setUserPassword)

  const [editingEmail, setEditingEmail] = useState(false)
  const [editingPassword, setEditingPassword] = useState(false)
  const [email, setEmail] = useState(user?.email || '')
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

  useEffect(() => {
    if (user && user.email) {
      setEmail(user.email || '')
    }
  }, [user])

  return (
    <div className="dark:text-foreground">
      <h1 className="text-2xl font-bold mb-4">User Settings</h1>

      <Card>
        <CardHeader>
          <CardTitle className="text-xl">Basic information</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 gap-x-2 text-sm md:grid-cols-2 lg:grid-cols-3">
            <div className="pr-8 dark:border-gray-700 md:border-r">
              <div className="text-sm font-semibold mb-2">User ID</div>
              <div className="text-sm text-gray-500">{user?.id}</div>
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
      <Card className="p-4 mt-4">
        <form onSubmit={handleEmailSubmit}>
          <label className="block w-full text-sm font-semibold mb-2">
            Email
          </label>
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
      </Card>
      <Card className="p-4 mt-4">
        <form onSubmit={handlePasswordSubmit}>
          <label className="block w-full text-sm font-semibold mb-2">
            Password
          </label>
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
      </Card>
    </div>
  )
}

export default UserSettingsPage
