import React, { FC, MouseEvent, useEffect, useState } from 'react';
import { DateTime } from 'luxon';
import { useUser } from '@/lib/auth';
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import { useMutation, useQuery } from '@connectrpc/connect-query';
import {
  createUserInvite,
  getOrganization,
  listSAMLConnections,
  listUserInvites,
  listUsers,
  updateOrganization,
  updateUser,
} from '@/gen/tesseral/frontend/v1/frontend-FrontendService_connectquery';
import {
  Table,
  TableBody,
  TableCell,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { Button } from '@/components/ui/button';
import { Link } from 'react-router-dom';
import { Switch } from '@/components/ui/switch';
import { parseErrorMessage } from '@/lib/errors';
import { toast } from 'sonner';
import Loader from '@/components/ui/loader';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogTitle,
  DialogTrigger,
  DialogHeader,
} from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';

const OrganizationSettingsPage: FC = () => {
  const user = useUser();

  const { data: usersData, refetch: refetchUsers } = useQuery(listUsers);
  const { data: organizationRes, refetch: refetchOrganization } =
    useQuery(getOrganization);
  const { data: samlConnectionsData } = useQuery(listSAMLConnections);
  const createUserInviteMutation = useMutation(createUserInvite);
  const updateOrganizationMutation = useMutation(updateOrganization);
  const updateUserMutation = useMutation(updateUser);
  const { data: userInvites, refetch: refetchUserInvites } =
    useQuery(listUserInvites);

  const [creatingUserInvite, setCreatingUserInvite] = useState(false);
  const [editingLoginSettings, setEditingLoginSettings] = useState(false);
  const [inviteUserVisible, setInviteUserVisible] = useState(false);
  const [inviteeEmail, setInviteeEmail] = useState('');
  const [inviteeIsOwner, setInviteeIsOwner] = useState(false);
  const [logInWithAuthenticatorApp, setLogInWithAuthenticatorApp] = useState(
    organizationRes?.organization?.logInWithAuthenticatorApp,
  );
  const [logInWithEmail, setLogInWithEmail] = useState(
    organizationRes?.organization?.logInWithEmail,
  );
  const [logInWithGoogle, setLogInWithGoogle] = useState(
    organizationRes?.organization?.logInWithGoogle,
  );
  const [logInWithMicrosoft, setLogInWithMicrosoft] = useState(
    organizationRes?.organization?.logInWithMicrosoft,
  );
  const [logInWithPasskey, setLogInWithPasskey] = useState(
    organizationRes?.organization?.logInWithPasskey,
  );
  const [logInWithPassword, setLogInWithPassword] = useState(
    organizationRes?.organization?.logInWithPassword,
  );
  const [requireMFA, setRequireMFA] = useState(
    organizationRes?.organization?.requireMfa,
  );
  const [submittingLoginSettings, setSubmittingLoginSettings] = useState(false);

  const changeUserRole = async (userId: string, isOwner: boolean) => {
    await updateUserMutation.mutateAsync({
      id: userId,
      user: {
        owner: isOwner,
      },
    });

    await refetchUsers();
  };

  const resetLoginSettings = () => {
    setLogInWithEmail(organizationRes?.organization?.logInWithEmail);
    setLogInWithGoogle(organizationRes?.organization?.logInWithGoogle);
    setLogInWithMicrosoft(organizationRes?.organization?.logInWithMicrosoft);
    setLogInWithPassword(organizationRes?.organization?.logInWithPassword);
  };

  const submitLoginSettings = async () => {
    setSubmittingLoginSettings(true);
    try {
      await updateOrganizationMutation.mutateAsync({
        organization: {
          logInWithEmail,
          logInWithGoogle,
          logInWithMicrosoft,
          logInWithPassword,
        },
      });

      setEditingLoginSettings(false);
      await refetchOrganization();
      resetLoginSettings();
      setSubmittingLoginSettings(false);
    } catch (error) {
      setSubmittingLoginSettings(false);
      const message = parseErrorMessage(error);
      toast.error('Could not update organization settings', {
        description: message,
      });
    }
  };

  const submitUserInvite = async () => {
    setCreatingUserInvite(true);
    try {
      await createUserInviteMutation.mutateAsync({
        userInvite: {
          email: inviteeEmail,
          owner: inviteeIsOwner,
        },
      });

      await refetchUserInvites();

      setInviteUserVisible(false);
      setCreatingUserInvite(false);
      setInviteeEmail('');
      setInviteeIsOwner(false);
    } catch (error) {
      const message = parseErrorMessage(error);
      toast.error('Could not invite user', {
        description: message,
      });
    }
  };

  useEffect(() => {
    if (organizationRes?.organization) {
      resetLoginSettings();
    }
  }, [organizationRes]);

  return (
    <div className="dark:text-foreground">
      <div className="mb-4">
        <h1 className="text-2xl font-bold mb-2">
          {organizationRes?.organization?.displayName}
        </h1>
        <span className="text-xs border px-2 py-1 rounded text-gray-400 dark:text-gray-700 bg-gray-200 dark:bg-gray-900 dark:border-gray-800">
          {organizationRes?.organization?.id}
        </span>
      </div>

      <Card className="mt-8">
        <CardHeader>
          <CardTitle>General configuration</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 gap-x-2 text-sm md:grid-cols-2 lg:grid-cols-3">
            <div className="border-gray-200 pr-8 dark:border-gray-700 lg:border-r">
              <div className="font-semibold mb-2">Display Name</div>
              <div className="text-sm text-gray-500">
                {organizationRes?.organization?.displayName}
              </div>
            </div>
            <div className=" border-gray-200 mt-8 pr-8 dark:border-gray-700 lg:px-8 lg:border-r md:mt-0">
              <div className="font-semibold mb-2">Created</div>
              <div className="text-sm text-gray-500">
                {organizationRes?.organization?.createTime &&
                  DateTime.fromSeconds(
                    parseInt(
                      `${organizationRes?.organization?.createTime.seconds}`,
                    ),
                  ).toRelative()}
              </div>
            </div>
            <div className="pr-8 mt-8 lg:px-8 lg:mt-0">
              <div className="font-semibold mb-2">Last updated</div>
              <div className="text-sm text-gray-500">
                {organizationRes?.organization?.updateTime
                  ? DateTime.fromSeconds(
                      parseInt(
                        `${organizationRes.organization.updateTime.seconds}`,
                      ),
                    ).toRelative()
                  : '—'}
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      <Card className="mt-8">
        <CardHeader>
          <CardTitle>
            <div className="flex items-center justify-between w-full ">
              <span>Log in settings</span>
              {editingLoginSettings ? (
                <div className="">
                  <Button
                    onClick={() => {
                      setEditingLoginSettings(false);
                      resetLoginSettings();
                    }}
                    variant="outline"
                  >
                    Cancel
                  </Button>
                  <Button
                    className="ml-2"
                    disabled={submittingLoginSettings}
                    onClick={submitLoginSettings}
                  >
                    {submittingLoginSettings && <Loader />}
                    Save
                  </Button>
                </div>
              ) : (
                <Button
                  onClick={() => setEditingLoginSettings(true)}
                  variant="outline"
                >
                  Edit
                </Button>
              )}
            </div>
          </CardTitle>
          <CardDescription>
            These settings control how users can log in to your organization.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <p className="font-semibold">Basic log in settings</p>
          <Table>
            <TableBody>
              <TableRow>
                <TableCell>
                  <p>Log in with Email</p>
                  <p className="text-xs text-muted-foreground">
                    Users can log in with their email address.
                  </p>
                </TableCell>
                <TableCell className="text-right">
                  <Switch
                    checked={logInWithEmail}
                    disabled={!editingLoginSettings}
                    onCheckedChange={setLogInWithEmail}
                  />
                </TableCell>
              </TableRow>
              <TableRow>
                <TableCell>
                  <p>Log in with Password</p>
                  <p className="text-xs text-muted-foreground">
                    Users will be prompted to enter a password during login.
                  </p>
                </TableCell>
                <TableCell className="text-right">
                  <Switch
                    checked={logInWithPassword}
                    disabled={!editingLoginSettings}
                    onCheckedChange={setLogInWithPassword}
                  />
                </TableCell>
              </TableRow>
              <TableRow>
                <TableCell>
                  <p>Log in with Google</p>
                  <p className="text-xs text-muted-foreground">
                    Users can log in with their Google account.
                  </p>
                </TableCell>
                <TableCell className="text-right">
                  <Switch
                    checked={logInWithGoogle}
                    disabled={!editingLoginSettings}
                    onCheckedChange={setLogInWithGoogle}
                  />
                </TableCell>
              </TableRow>
              <TableRow>
                <TableCell>
                  <p>Log in with Microsoft</p>
                  <p className="text-xs text-muted-foreground">
                    Users can log in with their Microsoft account.
                  </p>
                </TableCell>
                <TableCell className="text-right">
                  <Switch
                    checked={logInWithMicrosoft}
                    disabled={!editingLoginSettings}
                    onCheckedChange={setLogInWithMicrosoft}
                  />
                </TableCell>
              </TableRow>
            </TableBody>
          </Table>

          <p className="mt-4 font-semibold">MFA settings</p>
          <Table>
            <TableBody>
              <TableRow>
                <TableCell>
                  <p>Require MFA</p>
                  <p className="text-xs text-muted-foreground">
                    Require users to set up MFA during login.
                  </p>
                </TableCell>
                <TableCell className="text-right">
                  <Switch
                    checked={requireMFA}
                    disabled={!editingLoginSettings}
                    onCheckedChange={setRequireMFA}
                  />
                </TableCell>
              </TableRow>
              <TableRow>
                <TableCell>
                  <p>Log in with Authenticator App</p>
                  <p className="text-xs text-muted-foreground">
                    Users can log in with an authenticator app.
                  </p>
                </TableCell>
                <TableCell className="text-right">
                  <Switch
                    checked={logInWithAuthenticatorApp}
                    disabled={!editingLoginSettings}
                    onCheckedChange={setLogInWithAuthenticatorApp}
                  />
                </TableCell>
              </TableRow>
              <TableRow>
                <TableCell>
                  <p>Log in with Passkey</p>
                  <p className="text-xs text-muted-foreground">
                    Users can log in with a passkey.
                  </p>
                </TableCell>
                <TableCell className="text-right">
                  <Switch
                    checked={logInWithPasskey}
                    disabled={!editingLoginSettings}
                    onCheckedChange={setLogInWithPasskey}
                  />
                </TableCell>
              </TableRow>
            </TableBody>
          </Table>
        </CardContent>
      </Card>

      <Card className="mt-8">
        <CardHeader>
          <CardTitle>
            <div className="flex justify-between items-center w-full">
              <span>Users</span>
              <div>
                <Dialog
                  open={inviteUserVisible}
                  onOpenChange={setInviteUserVisible}
                >
                  <DialogTrigger asChild>
                    <Button variant="outline">Invite a User</Button>
                  </DialogTrigger>

                  <DialogContent>
                    <DialogHeader>
                      <DialogTitle>Invite a User</DialogTitle>
                      <DialogDescription></DialogDescription>
                    </DialogHeader>
                    <form onSubmit={submitUserInvite}>
                      <div className="mb-4">
                        <Label>Email</Label>
                        <Input
                          onChange={(e) => setInviteeEmail(e.target.value)}
                          type="email"
                          value={inviteeEmail}
                        />
                      </div>
                      <div className="mb-4 flex items-center">
                        <Label>Owner?</Label>
                        <Switch
                          checked={inviteeIsOwner}
                          className="ml-2"
                          onCheckedChange={setInviteeIsOwner}
                        />
                      </div>
                      <div className="flex justify-end">
                        <Button
                          variant="outline"
                          onClick={() => {
                            setInviteUserVisible(false);
                          }}
                        >
                          {creatingUserInvite && <Loader />}
                          Cancel
                        </Button>
                        <Button className="ml-2">Invite</Button>
                      </div>
                    </form>
                  </DialogContent>
                </Dialog>
              </div>
            </div>
          </CardTitle>
          <CardDescription></CardDescription>
        </CardHeader>
        <CardContent>
          <Tabs defaultValue="registered">
            <TabsList className="grid grid-cols-2 max-w-sm">
              <TabsTrigger value="registered">Registered</TabsTrigger>
              <TabsTrigger value="invited">Invited</TabsTrigger>
            </TabsList>
            <TabsContent value="registered">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableCell>ID</TableCell>
                    <TableCell>Email</TableCell>
                    <TableCell>Role</TableCell>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {usersData?.users.map((u) => (
                    <TableRow key={u.id}>
                      <TableCell className="flex items-center">
                        {u.id}
                      </TableCell>
                      <TableCell className="text-gray-500">{u.email}</TableCell>
                      <TableCell className="text-gray-500">
                        {u.owner ? 'Owner' : 'Member'}

                        {u.owner && u.id !== user?.id && (
                          <div
                            className="ml-2 rounded cursor-pointer text-primary border-border px-4 py-2 inline-block"
                            onClick={async (e: MouseEvent<HTMLSpanElement>) => {
                              e.stopPropagation();
                              e.preventDefault();

                              await changeUserRole(u.id, !u.owner);
                            }}
                          >
                            Make {u.owner ? 'Member' : 'Owner'}
                          </div>
                        )}
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </TabsContent>

            <TabsContent value="invited">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableCell>Email</TableCell>
                    <TableCell>Role</TableCell>
                    <TableCell>Created</TableCell>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {userInvites?.userInvites.map((i) => (
                    <TableRow key={i.id}>
                      <TableCell>{i.email}</TableCell>
                      <TableCell>{i.owner ? 'Owner' : 'Member'}</TableCell>
                      <TableCell>
                        {DateTime.fromSeconds(
                          parseInt(`${i.createTime?.seconds}`),
                        ).toRelative()}
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </TabsContent>
          </Tabs>
        </CardContent>
      </Card>

      <Card className="mt-8">
        <CardHeader>
          <CardTitle>SAML Connections</CardTitle>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableCell>ID</TableCell>
                <TableCell>IDP Entity ID</TableCell>
                <TableCell>IDP Redirect URL</TableCell>
                <TableCell>IDP X509 Certificate</TableCell>
                <TableCell></TableCell>
              </TableRow>
            </TableHeader>
            <TableBody>
              {samlConnectionsData?.samlConnections.map((c) => (
                <TableRow key={c.id}>
                  <TableCell className="flex items-center">{c.id}</TableCell>
                  <TableCell className="text-gray-500">
                    {c.idpEntityId}
                  </TableCell>
                  <TableCell className="text-gray-500">
                    {c.idpRedirectUrl}
                  </TableCell>
                  <TableCell className="text-gray-500">
                    {c.idpX509Certificate ? (
                      <a
                        className="font-medium underline underline-offset-2 decoration-muted-foreground/40"
                        download={`Certificate ${c.id}.crt`}
                        href={`data:text/plain;base64,${btoa(c.idpX509Certificate)}`}
                      >
                        Download (.crt)
                      </a>
                    ) : (
                      '-'
                    )}
                  </TableCell>
                  <TableCell>
                    <Link to={`/organization/saml-connections/${c.id}`}>
                      <Button variant="outline">Edit</Button>
                    </Link>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </CardContent>
      </Card>
    </div>
  );
};

export default OrganizationSettingsPage;
