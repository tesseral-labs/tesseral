import React, { FC, FormEvent, useState } from 'react'
import { useUser } from '@/lib/auth'
import { Card } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { useMutation } from '@connectrpc/connect-query'
import {
  setPassword as setUserPassword,
  updateUser,
} from '@/gen/openauth/frontend/v1/frontend-FrontendService_connectquery'

const UserSettingsPage: FC = () => {
  const user = useUser()
  const setPasswordMutation = useMutation(setUserPassword)
  const updateUserMutation = useMutation(updateUser)

  const [editingEmail, setEditingEmail] = useState(false)
  const [editingPassword, setEditingPassword] = useState(false)
  const [email, setEmail] = useState(user?.email || '')
  const [password, setPassword] = useState('')

  const handleEmailSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault()

    try {
      await updateUserMutation.mutateAsync({
        user: {
          email,
        },
      })
    } catch (error) {
      console.error(error)
    }
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

  return (
    <div className="dark:text-foreground">
      <h1 className="text-2xl font-bold mb-4">User Settings</h1>
      <Card className="p-4">
        <form onSubmit={handleEmailSubmit}>
          <label className="block w-full text-sm font-semibold mb-2">
            Email
          </label>
          <input
            className="text-sm bg-input rounded border border-border focus:border-primary w-[240px] mb-2 disabled:text-gray-400 disabled:bg-gray-200"
            disabled={!editingEmail}
            onChange={(e) => setEmail(e.target.value)}
            type="email"
            placeholder="jane.doe@example.com"
            value={email}
          />
          <div>
            {editingEmail ? (
              <>
                <Button
                  className="text-sm rounded border border-border focus:border-primary mb-2 mr-2"
                  onClick={() => setEditingEmail(false)}
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
                onClick={() => setEditingEmail(true)}
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
          <input
            className="text-sm bg-input rounded border border-border focus:border-primary w-[240px] mb-2 disabled:text-gray-400 disabled:bg-gray-200"
            disabled={!editingPassword}
            onChange={(e) => setPassword(e.target.value)}
            type="password"
            placeholder="•••••••••••••"
            value={password}
          />
          <div>
            {editingPassword ? (
              <>
                <Button
                  className="text-sm rounded border border-border focus:border-primary mb-2 mr-2"
                  onClick={() => setEditingPassword(false)}
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
                onClick={() => setEditingPassword(true)}
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
