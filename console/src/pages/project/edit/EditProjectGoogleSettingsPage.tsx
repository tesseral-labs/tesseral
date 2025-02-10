import { PageCodeSubtitle, PageDescription, PageTitle } from '@/components/page'
import { Title } from '@/components/Title'
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from '@/components/ui/breadcrumb'
import { Button } from '@/components/ui/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Switch } from '@/components/ui/switch'
import {
  getProject,
  updateProject,
} from '@/gen/openauth/backend/v1/backend-BackendService_connectquery'
import { parseErrorMessage } from '@/lib/errors'
import { useMutation, useQuery } from '@connectrpc/connect-query'
import React, { FC, FormEvent, useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { toast } from 'sonner'

const EditProjectGoogleSettingsPage: FC = () => {
  const { data: getProjectResponse, refetch: refetchProject } = useQuery(
    getProject,
    {},
  )
  const updateProjectMutation = useMutation(updateProject)

  const [logInWithGoogle, setLogInWithGoogle] = useState(
    getProjectResponse?.project?.logInWithGoogle,
  )
  const [googleOauthClientId, setGoogleOauthClientId] = useState('')
  const [googleOauthClientSecret, setGoogleOauthClientSecret] = useState('')

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault()

    try {
      await updateProjectMutation.mutateAsync({
        project: {
          logInWithGoogle,
          googleOauthClientId,
          googleOauthClientSecret,
        },
      })
      toast.success('Google settings saved')

      const { data: refetchedProjectResponse } = await refetchProject()

      setLogInWithGoogle(refetchedProjectResponse?.project?.logInWithGoogle)
      setGoogleOauthClientId(
        refetchedProjectResponse?.project?.googleOauthClientId || '',
      )
      setGoogleOauthClientSecret(
        refetchedProjectResponse?.project?.googleOauthClientSecret || '',
      )
    } catch (error) {
      const message = parseErrorMessage(error)
      toast.error('Failed to save Google settings', {
        description: message,
      })
    }
  }

  useEffect(() => {
    setLogInWithGoogle(getProjectResponse?.project?.logInWithGoogle)
  }, [getProjectResponse])

  return (
    <div>
      <Title title="Edit Project Google Settings" />

      <Breadcrumb>
        <BreadcrumbList>
          <BreadcrumbItem>
            <BreadcrumbLink asChild>
              <Link to="/">Home</Link>
            </BreadcrumbLink>
          </BreadcrumbItem>
          <BreadcrumbSeparator />
          <BreadcrumbItem>
            <BreadcrumbLink asChild>
              <Link to="/project-settings">Project settings</Link>
            </BreadcrumbLink>
          </BreadcrumbItem>
          <BreadcrumbSeparator />
          <BreadcrumbItem>
            <BreadcrumbPage>Log in with Google settings</BreadcrumbPage>
          </BreadcrumbItem>
        </BreadcrumbList>
      </Breadcrumb>

      <PageTitle>Log in with Google settings</PageTitle>
      <PageCodeSubtitle>{getProjectResponse?.project?.id}</PageCodeSubtitle>
      <PageDescription>
        Edit the Google log in settings for your Project.
      </PageDescription>

      <Card className="mt-8">
        <CardHeader>
          <CardTitle>Log in with Google settings</CardTitle>
          <CardDescription>Log in with Google settings</CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit}>
            <div className="grid grid-cols-2 gap-8">
              <div>
                <Label>Log in with Google</Label>
                <p className="text-sm text-muted-foreground">
                  Enable or disable log in with Google within your Project.
                </p>
              </div>
              <div className="text-right">
                <Switch
                  checked={logInWithGoogle}
                  onCheckedChange={setLogInWithGoogle}
                />
              </div>
            </div>
            <div className="grid grid-cols-2 gap-8 mt-4 pt-4 border-t">
              <div>
                <Label>Google OAuth Client ID</Label>
                <p className="text-sm text-muted-foreground">
                  The OAuth Client ID for your Google application.
                </p>
              </div>
              <div className="text-right">
                <Input
                  onChange={(e) => setGoogleOauthClientId(e.target.value)}
                  placeholder={getProjectResponse?.project?.googleOauthClientId}
                  value={googleOauthClientId}
                />
              </div>
            </div>
            <div className="grid grid-cols-2 gap-8 mt-4 pt-4 border-t">
              <div>
                <Label>Google OAuth Client Secret</Label>
                <p className="text-sm text-muted-foreground">
                  The OAuth Client Secret for your Google application.
                </p>
              </div>
              <div className="text-right">
                <Input
                  onChange={(e) => setGoogleOauthClientSecret(e.target.value)}
                  placeholder={
                    getProjectResponse?.project?.googleOauthClientId
                      ? '<encrypted>'
                      : ''
                  }
                  value={googleOauthClientSecret}
                />
              </div>
            </div>
            <div className="text-right mt-8">
              <Link to="/project-settings">
                <Button variant="outline" className="mr-4">
                  Cancel
                </Button>
              </Link>
              <Button>Save</Button>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  )
}

export default EditProjectGoogleSettingsPage
