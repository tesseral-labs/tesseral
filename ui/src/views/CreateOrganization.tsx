import { setAccessToken, setRefreshToken } from '@/auth'
import { Title } from '@/components/Title'
import { Button } from '@/components/ui/button'
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import {
  createOrganization,
  exchangeIntermediateSessionForSession,
  whoami,
} from '@/gen/openauth/intermediate/v1/intermediate-IntermediateService_connectquery'
import { useMutation, useQuery } from '@connectrpc/connect-query'
import React, { Dispatch, FC, SetStateAction, useState } from 'react'
import { useNavigate } from 'react-router'
import { refresh } from '@/gen/openauth/frontend/v1/frontend-FrontendService_connectquery'
import { LoginLayouts, LoginViews } from '@/lib/views'
import { Input } from '@/components/ui/input'
import { Organization } from '@/gen/openauth/intermediate/v1/intermediate_pb'
import { useLayout } from '@/lib/settings'
import { cn } from '@/lib/utils'
import { parseErrorMessage } from '@/lib/errors'
import { toast } from 'sonner'

interface CreateOrganizationProps {
  setView: Dispatch<SetStateAction<LoginViews>>
}

const CreateOrganization: FC<CreateOrganizationProps> = ({ setView }) => {
  const layout = useLayout()
  const navigate = useNavigate()
  const { data: whoamiRes } = useQuery(whoami)

  const [displayName, setDisplayName] = useState<string>('')

  const createOrganizationMutation = useMutation(createOrganization)
  const exchangeIntermediateSessionForSessionMutation = useMutation(
    exchangeIntermediateSessionForSession,
  )

  const refreshMutation = useMutation(refresh)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    try {
      await createOrganizationMutation.mutateAsync({
        displayName,
      })

      if (
        !whoamiRes?.intermediateSession?.googleUserId &&
        !whoamiRes?.intermediateSession?.microsoftUserId
      ) {
        setView(LoginViews.RegisterPassword)
        return
      }

      const { refreshToken } =
        await exchangeIntermediateSessionForSessionMutation.mutateAsync({})

      const { accessToken } = await refreshMutation.mutateAsync({})

      setRefreshToken(refreshToken)
      setAccessToken(accessToken)

      navigate('/settings')
    } catch (error) {
      const message = parseErrorMessage(error)

      toast.error('Could not create organization', {
        description: message,
      })
    }
  }

  return (
    <>
      <Title title="Create a new Organization" />

      <Card
        className={cn(
          'w-full max-w-sm',
          layout !== LoginLayouts.Centered && 'shadow-none border-0',
        )}
      >
        <CardHeader>
          <CardTitle className="text-center">
            Create a new Organization
          </CardTitle>
        </CardHeader>
        <CardContent>
          <form
            className="flex flex-col items-center w-full"
            onSubmit={handleSubmit}
          >
            <Input
              className="w-full mb-2"
              id="displayName"
              placeholder="Acme, Inc."
              value={displayName}
              onChange={(e) => setDisplayName(e.target.value)}
            />
            <Button className="w-full" type="submit">
              Create Organization
            </Button>
          </form>
        </CardContent>
        <CardFooter></CardFooter>
      </Card>
    </>
  )
}

export default CreateOrganization
