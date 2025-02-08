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
import { useMutation, useQuery } from '@connectrpc/connect-query'
// import { exchangeIntermediateSessionForNewOrganizationSession } from '@/gen/openauth/intermediate/v1/intermediate-IntermediateService_connectquery'
import { refresh } from '@/gen/openauth/frontend/v1/frontend-FrontendService_connectquery'
import { setAccessToken, setRefreshToken } from '@/auth'
import { Input } from '@/components/ui/input'
import {
  createProject,
  exchangeIntermediateSessionForSession,
  setOrganization,
  whoami,
} from '@/gen/openauth/intermediate/v1/intermediate-IntermediateService_connectquery'
import { parseErrorMessage } from '@/lib/errors'
import { toast } from 'sonner'
import Loader from '@/components/ui/loader'
import useSettings from '@/lib/settings'

interface CreateProjectViewProps {
  setView: Dispatch<React.SetStateAction<LoginView>>
}

const CreateProjectView: FC<CreateProjectViewProps> = ({ setView }) => {
  const navigate = useNavigate()
  const settings = useSettings()

  const [displayName, setDisplayName] = useState<string>('')
  const [submitting, setSubmitting] = useState<boolean>(false)

  const createProjectMutation = useMutation(createProject)
  const exchangeIntermediateSessionForSessionMutation = useMutation(
    exchangeIntermediateSessionForSession,
  )
  const refreshMutation = useMutation(refresh)
  const setOrganizationMutation = useMutation(setOrganization)
  const { data: whoamiRes, refetch: refetchWhoami } = useQuery(whoami)

  const deriveNextView = (): LoginView | undefined => {
    if (settings?.logInWithPassword) {
      return LoginView.RegisterPassword
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setSubmitting(true)

    try {
      const projectRes = await createProjectMutation.mutateAsync({
        displayName,
      })

      await setOrganizationMutation.mutateAsync({
        organizationId: projectRes?.project?.organizationId,
      })

      console.log(`whoamiRes`, whoamiRes)
      await refetchWhoami()
      console.log(`refetched whoamiRes`, whoamiRes)

      const nextView = deriveNextView()
      if (nextView) {
        setView(nextView)
        return
      }

      const { refreshToken } =
        await exchangeIntermediateSessionForSessionMutation.mutateAsync({})

      const { accessToken } = await refreshMutation.mutateAsync({})

      setRefreshToken(refreshToken)
      setAccessToken(accessToken)

      setSubmitting(false)
      navigate('/')
    } catch (error) {
      setSubmitting(false)
      const message = parseErrorMessage(error)
      toast.error(message)
    }
  }

  return (
    <>
      <Title title="Create a new Project" />

      <Card className="w-full max-w-sm">
        <CardHeader>
          <CardTitle>Create a new Project</CardTitle>
        </CardHeader>
        <CardContent>
          <form className="w-full" onSubmit={handleSubmit}>
            <Input
              id="displayName"
              placeholder="Acme, Inc."
              value={displayName}
              onChange={(e) => setDisplayName(e.target.value)}
            />
            <Button
              className="mt-2 w-full"
              disabled={displayName.length < 1 || submitting}
              type="submit"
            >
              {submitting && <Loader />}
              Create Project
            </Button>
          </form>
        </CardContent>
      </Card>
    </>
  )
}

export default CreateProjectView
