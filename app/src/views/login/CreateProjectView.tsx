import React, { Dispatch, FC, useState } from 'react'
import { LoginView } from '@/lib/views'
import { Title } from '@/components/Title'
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { useNavigate } from 'react-router'
import { useMutation } from '@connectrpc/connect-query'
// import { exchangeIntermediateSessionForNewOrganizationSession } from '@/gen/openauth/intermediate/v1/intermediate-IntermediateService_connectquery'
import { refresh } from '@/gen/openauth/frontend/v1/frontend-FrontendService_connectquery'
import { setAccessToken, setRefreshToken } from '@/auth'
import { Input } from '@/components/ui/input'

interface CreateProjectViewProps {
  setView: Dispatch<React.SetStateAction<LoginView>>
}

const CreateProjectView: FC<CreateProjectViewProps> = ({ setView }) => {
  const navigate = useNavigate()

  const [displayName, setDisplayName] = useState<string>('')

  // const exchangeIntermediateSessionForNewOrganizationMutation = useMutation(
  //   exchangeIntermediateSessionForNewOrganizationSession,
  // )

  const refreshMutation = useMutation(refresh)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    try {
      // const { refreshToken } =
      //   await exchangeIntermediateSessionForNewProjectSession.mutateAsync(
      //     {
      //       displayName,
      //     },
      //   )

      // const { accessToken } = await refreshMutation.mutateAsync({})

      // setRefreshToken(refreshToken)
      // setAccessToken(accessToken)

      navigate('/project-settings')
    } catch (error) {
      console.error(error)
    }
  }

  return (
    <>
      <Title title="Create a new Project" />

      <Card className="w-[clamp(320px,50%,420px)] mx-auto">
        <CardHeader>
          <CardTitle className="text-center uppercase text-foreground font-semibold text-sm tracking-wide mt-2">
            Create a new Project
          </CardTitle>
        </CardHeader>
        <CardContent className="flex flex-col items-center justify-center w-full">
          <form className="flex flex-col items-center" onSubmit={handleSubmit}>
            <Input
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
              Create Project
            </Button>
          </form>
        </CardContent>
        <CardFooter></CardFooter>
      </Card>
    </>
  )
}

export default CreateProjectView
