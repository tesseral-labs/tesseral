import { useParams } from 'react-router';
import { useMutation, useQuery } from '@connectrpc/connect-query';
import {
  createUserImpersonationToken,
  getOrganization,
  getProject,
  getUser,
  listPasskeys,
  listSessions,
  updateUser,
} from '@/gen/tesseral/backend/v1/backend-BackendService_connectquery';
import React, { FC, useEffect, useState } from 'react';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { Link } from 'react-router-dom';
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from '@/components/ui/breadcrumb';
import {
  PageCodeSubtitle,
  PageDescription,
  PageTitle,
} from '@/components/page';
import {
  DetailsGrid,
  DetailsGridColumn,
  DetailsGridEntry,
  DetailsGridKey,
  DetailsGridValue,
} from '@/components/details-grid';
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import { DateTime } from 'luxon';
import { timestampDate } from '@bufbuild/protobuf/wkt';
import { Button } from '@/components/ui/button';
import {
  AlertDialog,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from '@/components/ui/alert-dialog';
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
} from '@/components/ui/form';
import { useForm } from 'react-hook-form';
import { z } from 'zod';
import { Input } from '@/components/ui/input';
import { Switch } from '@/components/ui/switch';
import { toast } from 'sonner';
import { User } from '@/gen/tesseral/backend/v1/models_pb';

export const ViewUserPage = () => {
  const { organizationId, userId } = useParams();
  const { data: getOrganizationResponse } = useQuery(getOrganization, {
    id: organizationId,
  });
  const { data: getUserResponse, refetch } = useQuery(getUser, {
    id: userId,
  });
  const { data: listSessionsResponse } = useQuery(listSessions, {
    userId,
  });
  const { data: listPasskeysResponse } = useQuery(listPasskeys, {
    userId,
  });

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
              <Link to="/organizations">Organizations</Link>
            </BreadcrumbLink>
          </BreadcrumbItem>
          <BreadcrumbSeparator />
          <BreadcrumbItem>
            <BreadcrumbLink asChild>
              <Link to={`/organizations/${organizationId}`}>
                {getOrganizationResponse?.organization?.displayName}
              </Link>
            </BreadcrumbLink>
          </BreadcrumbItem>
          <BreadcrumbSeparator />
          <BreadcrumbItem>
            <BreadcrumbLink asChild>
              <Link to={`/organizations/${organizationId}/users`}>Users</Link>
            </BreadcrumbLink>
          </BreadcrumbItem>
          <BreadcrumbSeparator />
          <BreadcrumbItem>
            <BreadcrumbPage>{getUserResponse?.user?.email}</BreadcrumbPage>
          </BreadcrumbItem>
        </BreadcrumbList>
      </Breadcrumb>

      <PageTitle>{getUserResponse?.user?.email}</PageTitle>
      <PageCodeSubtitle>{userId}</PageCodeSubtitle>
      <PageDescription>
        A User is what people using your product log into.
      </PageDescription>

      <Card className="my-8">
        <CardHeader className="flex flex-row items-center justify-between space-y-4">
          <div>
            <CardTitle>General settings</CardTitle>
            <CardDescription>Basic settings for this User.</CardDescription>
          </div>
          <EditUserSettingsButton />
        </CardHeader>
        <CardContent>
          <DetailsGrid>
            <DetailsGridColumn>
              <DetailsGridEntry>
                <DetailsGridKey>Email</DetailsGridKey>
                <DetailsGridValue>
                  {getUserResponse?.user?.email}
                </DetailsGridValue>
              </DetailsGridEntry>
              <DetailsGridEntry>
                <DetailsGridKey>Owner</DetailsGridKey>
                <DetailsGridValue>
                  {getUserResponse?.user?.owner ? 'Yes' : 'No'}
                </DetailsGridValue>
              </DetailsGridEntry>
              <DetailsGridEntry>
                <DetailsGridKey>Authenticator App</DetailsGridKey>
                <DetailsGridValue>
                  {getUserResponse?.user?.hasAuthenticatorApp
                    ? 'Enabled'
                    : 'Not Enabled'}
                </DetailsGridValue>
              </DetailsGridEntry>
            </DetailsGridColumn>
            <DetailsGridColumn>
              <DetailsGridEntry>
                <DetailsGridKey>Google User ID</DetailsGridKey>
                <DetailsGridValue>
                  {getUserResponse?.user?.googleUserId || '-'}
                </DetailsGridValue>
              </DetailsGridEntry>
              <DetailsGridEntry>
                <DetailsGridKey>Microsoft User ID</DetailsGridKey>
                <DetailsGridValue>
                  {getUserResponse?.user?.microsoftUserId || '-'}
                </DetailsGridValue>
              </DetailsGridEntry>
            </DetailsGridColumn>
            <DetailsGridColumn>
              <DetailsGridEntry>
                <DetailsGridKey>Created</DetailsGridKey>
                <DetailsGridValue>
                  {getUserResponse?.user?.createTime &&
                    DateTime.fromJSDate(
                      timestampDate(getUserResponse?.user?.createTime),
                    ).toRelative()}
                </DetailsGridValue>
              </DetailsGridEntry>
              <DetailsGridEntry>
                <DetailsGridKey>Updated</DetailsGridKey>
                <DetailsGridValue>
                  {getUserResponse?.user?.updateTime &&
                    DateTime.fromJSDate(
                      timestampDate(getUserResponse?.user?.updateTime),
                    ).toRelative()}
                </DetailsGridValue>
              </DetailsGridEntry>
            </DetailsGridColumn>
          </DetailsGrid>
        </CardContent>
      </Card>

      <Card className="my-8">
        <CardHeader>
          <CardTitle>Sessions</CardTitle>
          <CardDescription>
            Every time your Users log in or perform an action, that's associated
            with a Session.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>ID</TableHead>
                <TableHead>Created</TableHead>
                <TableHead>Last Active</TableHead>
                <TableHead>Expiration</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {listSessionsResponse?.sessions?.map((session) => (
                <TableRow key={session.id}>
                  <TableCell className="font-medium font-mono">
                    {session.id}
                  </TableCell>
                  <TableCell>
                    {session?.createTime &&
                      DateTime.fromJSDate(
                        timestampDate(session.createTime),
                      ).toRelative()}
                  </TableCell>
                  <TableCell>
                    {session?.lastActiveTime &&
                      DateTime.fromJSDate(
                        timestampDate(session.lastActiveTime),
                      ).toRelative()}
                  </TableCell>
                  <TableCell>
                    {session?.expireTime &&
                      DateTime.fromJSDate(
                        timestampDate(session.expireTime),
                      ).toRelative()}
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Passkeys</CardTitle>
          <CardDescription>Passkeys associated with this User.</CardDescription>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>ID</TableHead>
                <TableHead>Created</TableHead>
                <TableHead>Updated</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {listPasskeysResponse?.passkeys?.map((passkey) => (
                <TableRow key={passkey.id}>
                  <TableCell className="font-medium font-mono">
                    <Link
                      className="font-mono font-medium underline underline-offset-2 decoration-muted-foreground/40"
                      to={`/organizations/${organizationId}/users/${userId}/passkeys/${passkey.id}`}
                    >
                      {passkey.id}
                    </Link>
                  </TableCell>
                  <TableCell>
                    {passkey?.createTime &&
                      DateTime.fromJSDate(
                        timestampDate(passkey.createTime),
                      ).toRelative()}
                  </TableCell>
                  <TableCell>
                    {passkey?.updateTime &&
                      DateTime.fromJSDate(
                        timestampDate(passkey.updateTime),
                      ).toRelative()}
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </CardContent>
      </Card>

      <DangerZoneCard />
    </div>
  );
};

const DangerZoneCard = () => {
  const { userId } = useParams();
  const createUserImpersonationTokenMutation = useMutation(
    createUserImpersonationToken,
  );
  const { data: project } = useQuery(getProject);

  const handleImpersonate = async () => {
    const { userImpersonationToken } =
      await createUserImpersonationTokenMutation.mutateAsync({
        userImpersonationToken: {
          impersonatedId: userId,
        },
      });

    window.location.href = `https://${project?.project?.vaultDomain}/impersonate?secret-user-impersonation-token=${userImpersonationToken?.secretToken}`;
  };

  return (
    <Card className="mt-8 border-destructive">
      <CardHeader>
        <CardTitle>Danger Zone</CardTitle>
      </CardHeader>

      <CardContent>
        <div className="flex justify-between items-center">
          <div>
            <div className="text-sm font-semibold">Impersonate User</div>
            <p className="text-sm">
              Impersonate this User. You will be logged in as this User. You can
              end the impersonated session by logging out.
            </p>
          </div>

          <Button variant="destructive" onClick={handleImpersonate}>
            Impersonate User
          </Button>
        </div>
      </CardContent>
    </Card>
  );
};

const schema = z.object({
  email: z.string().email(),
  owner: z.boolean(),
  googleUserId: z.string().optional(),
  microsoftUserId: z.string().optional(),
});

const EditUserSettingsButton: FC = () => {
  const { userId } = useParams();
  const form = useForm<z.infer<typeof schema>>({
    defaultValues: {
      email: '',
      owner: false,
      googleUserId: '',
      microsoftUserId: '',
    },
  });
  const { data: getUserResponse, refetch } = useQuery(getUser, {
    id: userId,
  });
  const updateUserMutation = useMutation(updateUser);

  const [open, setOpen] = useState(false);

  useEffect(() => {
    if (getUserResponse?.user) {
      form.reset({
        email: getUserResponse.user.email,
        owner: getUserResponse.user.owner,
        googleUserId: getUserResponse.user.googleUserId || '',
        microsoftUserId: getUserResponse.user.microsoftUserId || '',
      });
    }
  }, [getUserResponse]);

  const handleSubmit = async (data: z.infer<typeof schema>) => {
    const updatedUser: Partial<User> = {
      email: data.email,
      owner: data.owner,
    };

    if (data.googleUserId) {
      updatedUser.googleUserId = data.googleUserId;
    }
    if (data.microsoftUserId) {
      updatedUser.microsoftUserId = data.microsoftUserId;
    }

    await updateUserMutation.mutateAsync({
      id: userId,
      user: updatedUser as User,
    });

    await refetch();

    setOpen(false);
    toast.success('User settings updated successfully.');
  };

  return (
    <AlertDialog open={open} onOpenChange={setOpen}>
      <AlertDialogTrigger asChild>
        <Button variant="outline">Edit</Button>
      </AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Edit User Settings</AlertDialogTitle>
        </AlertDialogHeader>
        <Form {...form}>
          <form
            className="space-y-8"
            onSubmit={form.handleSubmit(handleSubmit)}
          >
            <FormField
              control={form.control}
              name="email"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Email</FormLabel>
                  <FormControl>
                    <Input
                      type="email"
                      placeholder="jane.doe@example.com"
                      {...field}
                    />
                  </FormControl>
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="googleUserId"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Google User ID</FormLabel>
                  <FormControl>
                    <Input
                      type="text"
                      placeholder="Google User ID"
                      {...field}
                    />
                  </FormControl>
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="microsoftUserId"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Microsoft User ID</FormLabel>
                  <FormControl>
                    <Input
                      type="text"
                      placeholder="Microsoft User ID"
                      {...field}
                    />
                  </FormControl>
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="owner"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Owner</FormLabel>
                  <FormControl>
                    <Switch
                      className="block"
                      checked={field.value}
                      onCheckedChange={field.onChange}
                    />
                  </FormControl>
                </FormItem>
              )}
            />

            <AlertDialogFooter>
              <AlertDialogCancel>Cancel</AlertDialogCancel>
              <Button type="submit">Save</Button>
            </AlertDialogFooter>
          </form>
        </Form>
      </AlertDialogContent>
    </AlertDialog>
  );
};
