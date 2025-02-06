import React, { useState } from 'react'
import { useMutation, useQuery } from '@connectrpc/connect-query'
import {
  createSCIMAPIKey,
  createUserInvite,
  listOrganizations,
  listUserInvites,
  listUsers,
} from '@/gen/openauth/backend/v1/backend-BackendService_connectquery'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { Link } from 'react-router-dom'
import { DateTime } from 'luxon'
import { timestampDate } from '@bufbuild/protobuf/wkt'
import { useNavigate, useParams } from 'react-router'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { z } from 'zod'
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
import { Switch } from '@/components/ui/switch'

export function OrganizationUserInvitesTab() {
  const { organizationId } = useParams()
  const { data: listUserInvitesResponse } = useQuery(listUserInvites, {
    organizationId,
  })

  return (
    <Card>
      <CardHeader className="flex-row justify-between items-center">
        <div className="flex flex-col space-y-1 5">
          <CardTitle>User Invites</CardTitle>
          <CardDescription>
            A user invite lets outside collaborators join an organization. Lorem
            ipsum dolor.
          </CardDescription>
        </div>
        <CreateUserInviteButton />
      </CardHeader>
      <CardContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Email</TableHead>
              <TableHead>ID</TableHead>
              <TableHead>Created At</TableHead>
              <TableHead>Updated At</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {listUserInvitesResponse?.userInvites?.map((userInvite) => (
              <TableRow key={userInvite.id}>
                <TableCell>
                  <Link
                    className="font-medium underline underline-offset-2 decoration-muted-foreground/40"
                    to={`/organizations/${organizationId}/user-invites/${userInvite.id}`}
                  >
                    {userInvite.email}
                  </Link>

                  {userInvite.owner && (
                    <Badge variant="outline" className="ml-2">
                      Owner
                    </Badge>
                  )}
                </TableCell>
                <TableCell className="font-mono">{userInvite.id}</TableCell>
                <TableCell>
                  {DateTime.fromJSDate(
                    timestampDate(userInvite.createTime!),
                  ).toRelative()}
                </TableCell>
                <TableCell>
                  {DateTime.fromJSDate(
                    timestampDate(userInvite.updateTime!),
                  ).toRelative()}
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  )
}

const schema = z.object({
  email: z.string().email(),
  owner: z.boolean(),
})

function CreateUserInviteButton() {
  const { organizationId } = useParams()
  const createUserInviteMutation = useMutation(createUserInvite)
  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      email: '',
      owner: false,
    },
  })
  const navigate = useNavigate()
  const [open, setOpen] = useState(false)

  async function handleSubmit(values: z.infer<typeof schema>) {
    const { userInvite } = await createUserInviteMutation.mutateAsync({
      userInvite: {
        organizationId,
        email: values.email,
        owner: values.owner,
      },
    })

    navigate(`/organizations/${organizationId}/user-invites/${userInvite!.id}`)
  }

  return (
    <AlertDialog open={open} onOpenChange={setOpen}>
      <AlertDialogTrigger>
        <Button variant="outline">Create</Button>
      </AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Create User Invite</AlertDialogTitle>
          <AlertDialogDescription>
            A user invite lets outside collaborators join an organization. Lorem
            ipsum dolor.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <Form {...form}>
          <form
            onSubmit={form.handleSubmit(handleSubmit)}
            className="space-y-8"
          >
            <FormField
              control={form.control}
              name="email"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Email</FormLabel>
                  <FormControl>
                    <Input className="max-w-96" {...field} />
                  </FormControl>
                  <FormDescription>
                    The outside collaborator's email. The collaborator will need
                    to verify this email before being able to join the
                    organization.
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="owner"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Invite as owner</FormLabel>
                  <FormControl>
                    <Switch
                      className="block"
                      checked={field.value}
                      onCheckedChange={field.onChange}
                    />
                  </FormControl>
                  <FormDescription>
                    Whether the collaborator will join as an owner.
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />
            <AlertDialogFooter className="mt-8">
              <AlertDialogCancel>Cancel</AlertDialogCancel>
              <Button type="submit">Create User Invite</Button>
            </AlertDialogFooter>
          </form>
        </Form>
      </AlertDialogContent>
    </AlertDialog>
  )
}
