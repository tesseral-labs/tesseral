import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from '@/components/ui/breadcrumb'
import { Link } from 'react-router-dom'
import React, { useEffect, useState } from 'react'
import { useNavigate, useParams } from 'react-router'
import { useMutation, useQuery } from '@connectrpc/connect-query'
import {
  createProjectAPIKey,
  deleteSAMLConnection,
  deleteProjectAPIKey,
  getOrganization,
  getSAMLConnection,
  getProjectAPIKey,
  revokeProjectAPIKey,
  updateProjectAPIKey,
} from '@/gen/openauth/backend/v1/backend-BackendService_connectquery'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { DateTime } from 'luxon'
import { timestampDate } from '@bufbuild/protobuf/wkt'
import { toast } from 'sonner'
import {
  AlertDialog,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from '@/components/ui/alert-dialog'
import { Button } from '@/components/ui/button'
import { z } from 'zod'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { SecretCopier } from '@/components/SecretCopier'
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { PageCodeSubtitle, PageDescription, PageTitle } from '@/components/page'

export function ViewProjectAPIKeyPage() {
  const { projectApiKeyId } = useParams()
  const { data: getProjectApiKeyResponse } = useQuery(getProjectAPIKey, {
    id: projectApiKeyId,
  })
  return (
    <div>
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
              <Link to="/project-settings">Project Settings</Link>
            </BreadcrumbLink>
          </BreadcrumbItem>
          <BreadcrumbSeparator />
          <BreadcrumbItem>
            <BreadcrumbLink asChild>
              <Link to="/project-settings/api-keys">API Keys</Link>
            </BreadcrumbLink>
          </BreadcrumbItem>
          <BreadcrumbSeparator />
          <BreadcrumbItem>
            <BreadcrumbPage>
              {getProjectApiKeyResponse?.projectApiKey?.displayName}
            </BreadcrumbPage>
          </BreadcrumbItem>
        </BreadcrumbList>
      </Breadcrumb>

      <PageTitle>
        {getProjectApiKeyResponse?.projectApiKey?.displayName}
      </PageTitle>
      <PageCodeSubtitle>{projectApiKeyId}</PageCodeSubtitle>
      <PageDescription>
        A Project API key is how your backend talks to the Tesseral Backend API.
        Lorem ipsum dolor.
      </PageDescription>

      <Card className="my-8">
        <CardHeader className="flex-row justify-between items-center">
          <div className="flex flex-col space-y-1 5">
            <CardTitle>Configuration</CardTitle>
            <CardDescription>Lorem ipsum dolor.</CardDescription>
          </div>
          <EditProjectAPIKeyButton />
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-3 gap-x-2 text-sm">
            <div className="border-r border-gray-200 pr-8 flex flex-col gap-4">
              <div>
                <div className="font-semibold">Display Name</div>
                <div className="truncate">
                  {getProjectApiKeyResponse?.projectApiKey?.displayName}
                </div>
              </div>
              <div>
                <div className="font-semibold">Revoked</div>
                <div className="truncate">
                  {getProjectApiKeyResponse?.projectApiKey?.revoked
                    ? 'Yes'
                    : 'No'}
                </div>
              </div>
            </div>
            <div className="border-r border-gray-200 pr-8 pl-8 flex flex-col gap-4">
              <div>
                <div className="font-semibold">Created</div>
                <div className="truncate">
                  {getProjectApiKeyResponse?.projectApiKey?.createTime &&
                    DateTime.fromJSDate(
                      timestampDate(
                        getProjectApiKeyResponse?.projectApiKey?.createTime,
                      ),
                    ).toRelative()}
                </div>
              </div>
            </div>
            <div className="border-gray-200 pl-8 flex flex-col gap-4">
              <div>
                <div className="font-semibold">Updated</div>
                <div className="truncate">
                  {getProjectApiKeyResponse?.projectApiKey?.updateTime &&
                    DateTime.fromJSDate(
                      timestampDate(
                        getProjectApiKeyResponse?.projectApiKey?.updateTime,
                      ),
                    ).toRelative()}
                </div>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      <DangerZoneCard />
    </div>
  )
}

const schema = z.object({
  displayName: z.string(),
})

function EditProjectAPIKeyButton() {
  const { projectApiKeyId } = useParams()
  const { data: getProjectAPIKeyResponse, refetch } = useQuery(
    getProjectAPIKey,
    {
      id: projectApiKeyId,
    },
  )
  const updateProjectAPIKeyMutation = useMutation(updateProjectAPIKey)
  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      displayName: '',
    },
  })
  useEffect(() => {
    if (getProjectAPIKeyResponse?.projectApiKey) {
      form.reset({
        displayName: getProjectAPIKeyResponse.projectApiKey.displayName,
      })
    }
  }, [getProjectAPIKeyResponse])

  const [open, setOpen] = useState(false)

  async function handleSubmit(values: z.infer<typeof schema>) {
    await updateProjectAPIKeyMutation.mutateAsync({
      id: projectApiKeyId,
      projectApiKey: {
        displayName: values.displayName,
      },
    })
    await refetch()
    setOpen(false)
  }

  return (
    <AlertDialog open={open} onOpenChange={setOpen}>
      <AlertDialogTrigger>
        <Button variant="outline">Edit</Button>
      </AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Edit Project API Key</AlertDialogTitle>
          <AlertDialogDescription>
            Edit Project API Key settings.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(handleSubmit)}>
            <FormField
              control={form.control}
              name="displayName"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Display Name</FormLabel>
                  <FormControl>
                    <Input className="max-w-96" {...field} />
                  </FormControl>
                  <FormDescription>
                    A human-friendly name for the Project API Key.
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />
            <AlertDialogFooter className="mt-8">
              <AlertDialogCancel>Cancel</AlertDialogCancel>
              <Button type="submit">Save</Button>
            </AlertDialogFooter>
          </form>
        </Form>
      </AlertDialogContent>
    </AlertDialog>
  )
}

function DangerZoneCard() {
  const { projectApiKeyId } = useParams()
  const { data: getProjectApiKeyResponse, refetch } = useQuery(
    getProjectAPIKey,
    {
      id: projectApiKeyId,
    },
  )

  const [confirmRevokeOpen, setConfirmRevokeOpen] = useState(false)
  function handleRevoke() {
    setConfirmRevokeOpen(true)
  }

  const revokeProjectApiKeyMutation = useMutation(revokeProjectAPIKey)
  async function handleConfirmRevoke() {
    await revokeProjectApiKeyMutation.mutateAsync({
      id: projectApiKeyId,
    })

    await refetch()
    toast.success('Project API Key revoked')
    setConfirmRevokeOpen(false)
  }

  const [confirmDeleteOpen, setConfirmDeleteOpen] = useState(false)

  function handleDelete() {
    setConfirmDeleteOpen(true)
  }

  const deleteProjectApiKeyMutation = useMutation(deleteProjectAPIKey)
  const navigate = useNavigate()
  const handleConfirmDelete = async () => {
    await deleteProjectApiKeyMutation.mutateAsync({
      id: projectApiKeyId,
    })

    toast.success('Project API Key deleted')
    navigate(`/project-settings/api-keys`)
  }

  return (
    <>
      <AlertDialog open={confirmRevokeOpen} onOpenChange={setConfirmRevokeOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>
              Revoke {getProjectApiKeyResponse?.projectApiKey?.displayName}?
            </AlertDialogTitle>
            <AlertDialogDescription>
              Revoking a Project API Key cannot be undone. Backend API calls
              from {getProjectApiKeyResponse?.projectApiKey?.displayName} will
              stop working. This cannot be undone.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <Button variant="destructive" onClick={handleConfirmRevoke}>
              Revoke Project API Key
            </Button>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      <AlertDialog open={confirmDeleteOpen} onOpenChange={setConfirmDeleteOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>
              Delete {getProjectApiKeyResponse?.projectApiKey?.displayName}?
            </AlertDialogTitle>
            <AlertDialogDescription>
              Deleting a Project API Key cannot be undone.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <Button variant="destructive" onClick={handleConfirmDelete}>
              Delete Project API Key
            </Button>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      <Card className="border-destructive">
        <CardHeader>
          <CardTitle>Danger Zone</CardTitle>
        </CardHeader>

        <CardContent className="space-y-8">
          <div className="flex justify-between items-center">
            <div>
              <div className="text-sm font-semibold">
                Revoke Project API Key
              </div>
              <p className="text-sm">
                Revoke this Project API Key. Backend API calls from this key
                will stop working. This cannot be undone.
              </p>
            </div>

            <Button
              variant="destructive"
              disabled={getProjectApiKeyResponse?.projectApiKey?.revoked}
              onClick={handleRevoke}
            >
              Revoke Project API Key
            </Button>
          </div>

          <div className="flex justify-between items-center">
            <div>
              <div className="text-sm font-semibold">
                Delete Project API Key
              </div>
              <p className="text-sm">
                Delete this Project API Key. You must revoke this Project API
                Key first.
              </p>
            </div>

            <Button
              variant="destructive"
              disabled={!getProjectApiKeyResponse?.projectApiKey?.revoked}
              onClick={handleDelete}
            >
              Delete Project API Key
            </Button>
          </div>
        </CardContent>
      </Card>
    </>
  )
}
