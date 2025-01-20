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
import { exchangeIntermediateSessionForNewOrganizationSession } from '@/gen/openauth/intermediate/v1/intermediate-IntermediateService_connectquery'
import { useMutation } from '@connectrpc/connect-query'
import React, { useState } from 'react'
import { useNavigate } from 'react-router'
import {
  getAccessToken
} from '@/gen/openauth/frontend/v1/frontend-FrontendService_connectquery'

const CreateOrganization = () => {
  const navigate = useNavigate()

  const [displayName, setDisplayName] = useState<string>('')

  const exchangeIntermediateSessionForNewOrganizationMutation = useMutation(
    exchangeIntermediateSessionForNewOrganizationSession,
  )

  const getAccessTokenMutation = useMutation(getAccessToken)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    try {
      const { refreshToken } =
        await exchangeIntermediateSessionForNewOrganizationMutation.mutateAsync({
          displayName,
        })

      const { accessToken } = await getAccessTokenMutation.mutateAsync({})

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
            <input
              className="text-sm rounded border border-border focus:border-primary w-[clamp(240px,50%,100%)] mb-2"
              id="displayName"
              placeholder="Acme, Inc."
              value={displayName}
              onChange={(e) => setDisplayName(e.target.value)}
            />
            <Button
              className="text-sm rounded border border-border focus:border-primary w-[clamp(240px,50%,100%)] mb-2"
              type="submit"
            >
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
