import { useMutation, useQuery } from '@connectrpc/connect-query'
import {
  createProjectAPIKey,
  listProjectAPIKeys,
} from '@/gen/openauth/backend/v1/backend-BackendService_connectquery'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { Link } from 'react-router-dom'
import React, { useState } from 'react'
import { DateTime } from 'luxon'
import { timestampDate } from '@bufbuild/protobuf/wkt'
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from '@/components/ui/breadcrumb'
import { PageDescription, PageTitle } from '@/components/page'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { z } from 'zod'
import { useNavigate, useParams } from 'react-router'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
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

export function ListProjectAPIKeysPage() {
  const { data: listProjectAPIKeysResponse } = useQuery(listProjectAPIKeys, {})

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
            <BreadcrumbPage>Project API Keys</BreadcrumbPage>
          </BreadcrumbItem>
        </BreadcrumbList>
      </Breadcrumb>

      <div className="flex justify-between items-center">
        <PageTitle>Project API Keys</PageTitle>
        <CreateProjectAPIKeyButton />
      </div>
      <PageDescription>
        A Project API Key is how your backend talks to the Tesseral Backend API.
        Lorem ipsum dolor.
      </PageDescription>

      <Card className="mt-8 overflow-hidden">
        <CardContent className="-m-6 mt-0">
          <Table>
            <TableHeader className="bg-gray-50">
              <TableRow>
                <TableCell>Display Name</TableCell>
                <TableHead>ID</TableHead>
                <TableHead>Status</TableHead>
                <TableHead>Created At</TableHead>
                <TableHead>Updated At</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {listProjectAPIKeysResponse?.projectApiKeys?.map(
                (projectAPIKey) => (
                  <TableRow key={projectAPIKey.id}>
                    <TableCell className="font-medium">
                      <Link
                        className="font-medium underline underline-offset-2 decoration-muted-foreground/40"
                        to={`/project-api-keys/${projectAPIKey.id}`}
                      >
                        {projectAPIKey.displayName}
                      </Link>
                    </TableCell>
                    <TableCell className="font-mono">
                      {projectAPIKey.id}
                    </TableCell>
                    <TableCell>
                      {projectAPIKey?.revoked ? 'Revoked' : 'Active'}
                    </TableCell>
                    <TableCell>
                      {DateTime.fromJSDate(
                        timestampDate(projectAPIKey.createTime!),
                      ).toRelative()}
                    </TableCell>
                    <TableCell>
                      {DateTime.fromJSDate(
                        timestampDate(projectAPIKey.updateTime!),
                      ).toRelative()}
                    </TableCell>
                  </TableRow>
                ),
              )}
            </TableBody>
          </Table>
        </CardContent>
      </Card>
    </div>
  )
}

const schema = z.object({
  displayName: z.string(),
})

function CreateProjectAPIKeyButton() {
  const createProjectAPIKeyMutation = useMutation(createProjectAPIKey)
  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      displayName: '',
    },
  })
  const navigate = useNavigate()
  const [createOpen, setCreateOpen] = useState(false)
  const [projectAPIKeyID, setProjectAPIKeyID] = useState('')
  const [secretToken, setSecretToken] = useState('')

  async function handleSubmit(values: z.infer<typeof schema>) {
    const { projectApiKey } = await createProjectAPIKeyMutation.mutateAsync({
      projectApiKey: {
        displayName: values.displayName,
      },
    })

    setCreateOpen(false)
    setProjectAPIKeyID(projectApiKey!.id)
    setSecretToken(projectApiKey!.secretToken)
  }

  function handleClose() {
    navigate(`/project-api-keys/${projectAPIKeyID}`)
  }

  return (
    <>
      <AlertDialog open={!!secretToken}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Project API Key Created</AlertDialogTitle>
            <AlertDialogDescription>
              Project API Key was created successfully.
            </AlertDialogDescription>
          </AlertDialogHeader>

          <div className="text-sm font-medium leading-none">
            Project Secret Token
          </div>

          <SecretCopier
            placeholder="tesseral_secret_key_•••••••••••••••••••••••••"
            secret={secretToken}
          />

          <div className="text-sm text-muted-foreground">
            Store this secret as TESSERAL_API_KEY in your secrets manager. You
            will not be able to see this secret token again later.
          </div>

          <AlertDialogFooter>
            <AlertDialogCancel onClick={handleClose}>Close</AlertDialogCancel>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      <AlertDialog open={createOpen} onOpenChange={setCreateOpen}>
        <AlertDialogTrigger>
          <Button variant="outline">Create</Button>
        </AlertDialogTrigger>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Create Project API Key</AlertDialogTitle>
            <AlertDialogDescription>
              A Project API key is how your backend talks to the Tesseral
              Backend API. Lorem ipsum dolor.
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
                <Button type="submit">Create Project API Key</Button>
              </AlertDialogFooter>
            </form>
          </Form>
        </AlertDialogContent>
      </AlertDialog>
    </>
  )
}
