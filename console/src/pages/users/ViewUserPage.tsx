import { useParams } from 'react-router';
import { useMutation, useQuery } from '@connectrpc/connect-query';
import {
  createUserImpersonationToken,
  deleteUserRoleAssignment,
  getOrganization,
  getProject,
  getRole,
  getUser,
  listPasskeys,
  listSessions,
  listUserRoleAssignments,
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
  PageCodeSubtitle,
  PageContent,
  PageDescription,
  PageHeader,
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
  AlertDialogAction,
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
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form';
import { useForm } from 'react-hook-form';
import { z } from 'zod';
import { Input } from '@/components/ui/input';
import { Switch } from '@/components/ui/switch';
import { toast } from 'sonner';
import { User, UserRoleAssignment } from '@/gen/tesseral/backend/v1/models_pb';
import { AssignUserRolesButton } from './AssignUserRolesButton';

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
  const { data: listUserRoleAssignmentsResponse } = useQuery(
    listUserRoleAssignments,
    {
      userId,
    },
  );

  return (
    <>
      <PageHeader>
        <PageTitle>{getUserResponse?.user?.email}</PageTitle>
        <PageCodeSubtitle>{userId}</PageCodeSubtitle>
        <PageDescription>
          A User is what people using your product log into.
        </PageDescription>
      </PageHeader>
      <PageContent>
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
                  <DetailsGridKey>Display Name</DetailsGridKey>
                  <DetailsGridValue>
                    {getUserResponse?.user?.displayName || '-'}
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
                  <DetailsGridKey>Profile Picture URL</DetailsGridKey>
                  <DetailsGridValue>
                    {getUserResponse?.user?.profilePictureUrl || '-'}
                  </DetailsGridValue>
                </DetailsGridEntry>
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
              Every time your Users log in or perform an action, that's
              associated with a Session.
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
            <CardDescription>
              Passkeys associated with this User.
            </CardDescription>
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

        <Card className="mt-8">
          <CardHeader className="flex-row justify-between items-center gap-x-2">
            <div className="flex flex-col space-y-1.5">
              <CardTitle>Assigned Roles</CardTitle>
              <CardDescription>
                Roles this User has been assigned.
              </CardDescription>
            </div>

            <div className="shrink-0 space-x-4">
              <AssignUserRolesButton />
            </div>
          </CardHeader>
          <CardContent>
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Display Name</TableHead>
                  <TableHead>Description</TableHead>
                  <TableHead></TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {listUserRoleAssignmentsResponse?.userRoleAssignments?.map(
                  (userRoleAssignment) => (
                    <UserRoleAssignmentRow
                      key={userRoleAssignment.roleId}
                      userRoleAssignment={userRoleAssignment}
                    />
                  ),
                )}
              </TableBody>
            </Table>
          </CardContent>
        </Card>

        <DangerZoneCard />
      </PageContent>
    </>
  );
};

function UserRoleAssignmentRow({
  userRoleAssignment,
}: {
  userRoleAssignment: UserRoleAssignment;
}) {
  const { refetch } = useQuery(listUserRoleAssignments, {
    userId: userRoleAssignment.userId,
  });
  const { data: getRoleResponse } = useQuery(getRole, {
    id: userRoleAssignment.roleId,
  });

  const { data: getUserResponse } = useQuery(getUser, {
    id: userRoleAssignment.userId,
  });

  const { mutateAsync: deleteUserRoleAssignmentAsync } = useMutation(
    deleteUserRoleAssignment,
  );

  async function handleUnassign() {
    await deleteUserRoleAssignmentAsync({ id: userRoleAssignment.id });
    await refetch();
    toast.success('Role unassigned');
  }

  const [open, setOpen] = useState(false);

  return (
    <TableRow>
      <TableCell>
        <Link
          className="font-medium underline underline-offset-2 decoration-muted-foreground/40"
          to={`/roles/${userRoleAssignment.roleId}`}
        >
          {getRoleResponse?.role?.displayName}
        </Link>
      </TableCell>
      <TableCell>{getRoleResponse?.role?.description}</TableCell>
      <TableCell className="text-right">
        <AlertDialog open={open} onOpenChange={setOpen}>
          <AlertDialogTrigger asChild>
            <Button size="sm" variant="link">
              Unassign
            </Button>
          </AlertDialogTrigger>
          <AlertDialogContent>
            <AlertDialogHeader>
              <AlertDialogTitle>Unassign Role</AlertDialogTitle>
            </AlertDialogHeader>
            <AlertDialogDescription>
              Are you sure you want to unassign{' '}
              <span className="font-medium">
                {getUserResponse?.user?.email}
              </span>{' '}
              from the Role{' '}
              <span className="font-medium">
                {getRoleResponse?.role?.displayName}
              </span>
              ?
            </AlertDialogDescription>
            <AlertDialogFooter>
              <AlertDialogCancel>Cancel</AlertDialogCancel>
              <AlertDialogAction onClick={handleUnassign}>
                Unassign
              </AlertDialogAction>
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialog>
      </TableCell>
    </TableRow>
  );
}

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
  displayName: z.string().optional(),
  profilePictureUrl: z.string().optional(),
});

const EditUserSettingsButton: FC = () => {
  const { userId } = useParams();
  const form = useForm<z.infer<typeof schema>>({
    defaultValues: {
      email: '',
      owner: false,
      googleUserId: '',
      microsoftUserId: '',
      displayName: '',
      profilePictureUrl: '',
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
        displayName: getUserResponse.user.displayName || '',
        profilePictureUrl: getUserResponse.user.profilePictureUrl || '',
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
    if (data.displayName) {
      updatedUser.displayName = data.displayName;
    }
    if (data.profilePictureUrl) {
      updatedUser.profilePictureUrl = data.profilePictureUrl;
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
                  <FormDescription>
                    The User's email address. Must be unique within their
                    Organization.
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="displayName"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Display Name</FormLabel>
                  <FormControl>
                    <Input type="text" placeholder="John Doe" {...field} />
                  </FormControl>
                  <FormDescription>
                    The User's display name. This is typically their full
                    personal name. Optional.
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="profilePictureUrl"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Profile Picture URL</FormLabel>
                  <FormControl>
                    <Input type="text" placeholder="https://..." {...field} />
                  </FormControl>
                  <FormDescription>
                    The URL of the User's profile picture. Optional.
                  </FormDescription>
                  <FormMessage />
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
                  <FormDescription>
                    The User's Google-assigned ID. Optional.
                  </FormDescription>
                  <FormMessage />
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
                  <FormDescription>
                    The User's Microsoft-assigned ID. Optional.
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
                  <FormLabel>Owner</FormLabel>
                  <FormControl>
                    <Switch
                      className="block"
                      checked={field.value}
                      onCheckedChange={field.onChange}
                    />
                  </FormControl>
                  <FormDescription>
                    Whether the User is an Owner of their Organization.
                    Optional.
                  </FormDescription>
                  <FormMessage />
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
