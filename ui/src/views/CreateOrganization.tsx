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
import { LoginViews } from '@/lib/views'
import { Input } from '@/components/ui/input'

interface CreateOrganizationProps {
  setView: Dispatch<SetStateAction<LoginViews>>
}

const CreateOrganization: FC<CreateOrganizationProps> = ({ setView }) => {
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
      console.error(error)
    }
  }

  return (
    <>
      <Title title="Create a new Organization" />

      <Card className="w-[clamp(320px,50%,420px)]">
        <CardHeader>
          <CardTitle className="text-center uppercase text-foreground font-semibold text-sm tracking-wide mt-2">
            Create a new Organization
          </CardTitle>
        </CardHeader>
        <CardContent className="flex flex-col items-center justify-center w-full">
          <form className="flex flex-col items-center" onSubmit={handleSubmit}>
            <Input
              className="w-[clamp(240px,50%,100%)] mb-2"
              id="displayName"
              placeholder="Acme, Inc."
              value={displayName}
              onChange={(e) => setDisplayName(e.target.value)}
            />
            <Button className="w-[clamp(240px,50%,100%)]" type="submit">
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
