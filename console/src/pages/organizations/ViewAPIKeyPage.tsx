import React, { useEffect, useState } from 'react';
import {
  PageCodeSubtitle,
  PageContent,
  PageDescription,
  PageHeader,
  PageTitle,
} from '@/components/page';
import {
  deleteAPIKey,
  deleteAPIKeyRoleAssignment,
  getAPIKey,
  getRole,
  listAPIKeyRoleAssignments,
  revokeAPIKey,
  updateAPIKey,
} from '@/gen/tesseral/backend/v1/backend-BackendService_connectquery';
import { useMutation, useQuery } from '@connectrpc/connect-query';
import { useNavigate, useParams } from 'react-router';
import {
  ConsoleCard,
  ConsoleCardDetails,
  ConsoleCardContent,
  ConsoleCardDescription,
  ConsoleCardHeader,
  ConsoleCardTitle,
  ConsoleCardTableContent,
} from '@/components/ui/console-card';
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
import { Button } from '@/components/ui/button';
import { useForm } from 'react-hook-form';
import { z } from 'zod';
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form';
import { Input } from '@/components/ui/input';
import {
  DetailsGrid,
  DetailsGridColumn,
  DetailsGridEntry,
  DetailsGridKey,
  DetailsGridValue,
} from '@/components/details-grid';
import { DateTime } from 'luxon';
import { timestampDate } from '@bufbuild/protobuf/wkt';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { toast } from 'sonner';
import { AddAPIKeyRoleButton } from './AddAPIKeyRoleButton';
import { APIKeyRoleAssignment } from '@/gen/tesseral/backend/v1/models_pb';
import { zodResolver } from '@hookform/resolvers/zod';

export const ViewAPIKeyPage = () => {
  const { apiKeyId } = useParams();

  const { data: getAPIKeyResponse } = useQuery(getAPIKey, {
    id: apiKeyId,
  });
  const { data: listApiKeyRoleAssignmentsResponse } = useQuery(
    listAPIKeyRoleAssignments,
    {
      apiKeyId,
    },
  );

  return (
    <>
      <PageHeader>
        <PageTitle>{getAPIKeyResponse?.apiKey?.displayName}</PageTitle>
        <PageCodeSubtitle>{getAPIKeyResponse?.apiKey?.id}</PageCodeSubtitle>
        <PageDescription>View and manage your API key.</PageDescription>
      </PageHeader>
      <PageContent>
        <div className="space-y-8">
          <ConsoleCard>
            <ConsoleCardHeader>
              <ConsoleCardDetails>
                <ConsoleCardTitle>API Key Details</ConsoleCardTitle>
                <ConsoleCardDescription></ConsoleCardDescription>
              </ConsoleCardDetails>

              <EditAPIKeyButton />
            </ConsoleCardHeader>
            <ConsoleCardContent>
              <DetailsGrid>
                <DetailsGridColumn>
                  <DetailsGridEntry>
                    <DetailsGridKey>Created at</DetailsGridKey>
                    <DetailsGridValue>
                      {getAPIKeyResponse?.apiKey?.createTime &&
                        DateTime.fromJSDate(
                          timestampDate(getAPIKeyResponse?.apiKey?.createTime),
                        ).toRelative()}
                    </DetailsGridValue>
                  </DetailsGridEntry>
                </DetailsGridColumn>
                <DetailsGridColumn>
                  <DetailsGridEntry>
                    <DetailsGridKey>Updated at</DetailsGridKey>
                    <DetailsGridValue>
                      {getAPIKeyResponse?.apiKey?.updateTime &&
                        DateTime.fromJSDate(
                          timestampDate(getAPIKeyResponse?.apiKey?.updateTime),
                        ).toRelative()}
                    </DetailsGridValue>
                  </DetailsGridEntry>
                </DetailsGridColumn>
                <DetailsGridColumn>
                  <DetailsGridEntry>
                    <DetailsGridKey>Status</DetailsGridKey>
                    <DetailsGridValue>
                      {getAPIKeyResponse?.apiKey?.revoked
                        ? 'Revoked'
                        : 'Active'}
                    </DetailsGridValue>
                  </DetailsGridEntry>
                </DetailsGridColumn>
                <DetailsGridColumn>
                  <DetailsGridEntry>
                    <DetailsGridKey>Expires at</DetailsGridKey>
                    <DetailsGridValue>
                      {getAPIKeyResponse?.apiKey?.expireTime ? (
                        <>
                          {DateTime.fromJSDate(
                            timestampDate(
                              getAPIKeyResponse?.apiKey?.expireTime,
                            ),
                          ).toRelative()}
                        </>
                      ) : (
                        <>{'Never'}</>
                      )}
                    </DetailsGridValue>
                  </DetailsGridEntry>
                </DetailsGridColumn>
              </DetailsGrid>
            </ConsoleCardContent>
          </ConsoleCard>

          <ConsoleCard>
            <ConsoleCardHeader>
              <ConsoleCardDetails>
                <ConsoleCardTitle>API Key Roles</ConsoleCardTitle>
                <ConsoleCardDescription>
                  Manage the roles associated with this API key.
                </ConsoleCardDescription>
              </ConsoleCardDetails>

              <AddAPIKeyRoleButton />
            </ConsoleCardHeader>
            <ConsoleCardTableContent>
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Role</TableHead>
                    <TableHead>Actions</TableHead>
                    <TableHead></TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {listApiKeyRoleAssignmentsResponse?.apiKeyRoleAssignments?.map(
                    (roleAssignment) => (
                      <APIKeyRoleAssignmentRow
                        key={roleAssignment.id}
                        apiKeyRoleAssignment={roleAssignment}
                      />
                    ),
                  )}
                </TableBody>
              </Table>
            </ConsoleCardTableContent>
          </ConsoleCard>

          <DangerZoneCard />
        </div>
      </PageContent>
    </>
  );
};

const APIKeyRoleAssignmentRow = ({
  apiKeyRoleAssignment,
}: {
  apiKeyRoleAssignment: APIKeyRoleAssignment;
}) => {
  const { data: getRoleResponse } = useQuery(getRole, {
    id: apiKeyRoleAssignment.roleId,
  });

  return (
    <TableRow key={getRoleResponse?.role?.id}>
      <TableCell>{getRoleResponse?.role?.displayName}</TableCell>
      <TableCell className="space-x-2">
        {getRoleResponse?.role?.actions.map((action) => (
          <span
            key={action}
            className="p-1 text-xs text-mono bg-muted text-muted-foreground rounded"
          >
            {action}
          </span>
        ))}
      </TableCell>
      <TableCell className="text-right">
        <RemoveRoleButton id={apiKeyRoleAssignment.id} />
      </TableCell>
    </TableRow>
  );
};

const schema = z.object({
  displayName: z.string().min(1, { message: 'Display name is required' }),
});

const EditAPIKeyButton = () => {
  const { apiKeyId } = useParams();
  const { data: getAPIKeyResponse } = useQuery(getAPIKey, {
    id: apiKeyId,
  });
  const updateAPIKeyMutation = useMutation(updateAPIKey);

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      displayName: '',
    },
  });

  const handleSubmit = async (data: z.infer<typeof schema>) => {
    await updateAPIKeyMutation.mutateAsync({
      id: apiKeyId,
      apiKey: {
        displayName: data.displayName,
      },
    });

    toast.success('API key updated successfully');
  };

  useEffect(() => {
    if (getAPIKeyResponse?.apiKey) {
      form.reset({
        displayName: getAPIKeyResponse.apiKey.displayName,
      });
    }
  }, [getAPIKeyResponse]);

  return (
    <AlertDialog>
      <AlertDialogTrigger asChild>
        <Button variant="outline">Edit</Button>
      </AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Edit API Key</AlertDialogTitle>
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
                    <Input placeholder="Display Name" {...field} />
                  </FormControl>
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
  );
};

const RemoveRoleButton = ({ id }: { id: string }) => {
  const { apiKeyId } = useParams();
  const [open, setOpen] = useState(false);
  const { refetch } = useQuery(listAPIKeyRoleAssignments, {
    apiKeyId,
  });
  const deleteAPIKeyRoleAssignmentMutation = useMutation(
    deleteAPIKeyRoleAssignment,
  );

  const handleDelete = async () => {
    await deleteAPIKeyRoleAssignmentMutation.mutateAsync({
      id,
    });

    toast.success('Role removed successfully');
    await refetch();
    setOpen(false);
  };

  return (
    <AlertDialog open={open} onOpenChange={setOpen}>
      <AlertDialogTrigger asChild>
        <Button variant="destructive">Remove</Button>
      </AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Remove Role</AlertDialogTitle>
          <AlertDialogDescription>
            Are you sure you want to remove this role from the API key? This
            disable all actions associated with this role for this API key.
          </AlertDialogDescription>
        </AlertDialogHeader>

        <AlertDialogFooter>
          <AlertDialogCancel>Cancel</AlertDialogCancel>
          <Button variant="destructive" onClick={() => handleDelete()}>
            Remove
          </Button>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
};

const DangerZoneCard = () => {
  const navigate = useNavigate();
  const { organizationId, apiKeyId } = useParams();
  const { data: getAPIKeyResponse, refetch } = useQuery(getAPIKey, {
    id: apiKeyId,
  });
  const deleteAPIKeyMutation = useMutation(deleteAPIKey);
  const revokeAPIKeyMutation = useMutation(revokeAPIKey);

  const [revokeOpen, setRevokeOpen] = useState(false);

  const handleDelete = async () => {
    await deleteAPIKeyMutation.mutateAsync({
      id: apiKeyId,
    });

    toast.success('API key deleted successfully');

    navigate(`/organizations/${organizationId}/api-keys`);
  };

  const handleRevoke = async () => {
    await revokeAPIKeyMutation.mutateAsync({
      id: apiKeyId,
    });

    toast.success('API key revoked successfully');

    await refetch();
  };

  return (
    <ConsoleCard className="border-destructive">
      <ConsoleCardHeader>
        <ConsoleCardTitle>Danger Zone</ConsoleCardTitle>
        <ConsoleCardDescription>
          Actions in this section cannot be undone. Please proceed with caution.
        </ConsoleCardDescription>
      </ConsoleCardHeader>
      <ConsoleCardContent>
        <div className="space-y-8">
          {!getAPIKeyResponse?.apiKey?.revoked ? (
            <div className="flex justify-between items-center space-y-2">
              <div>
                <div className="font-semibold text-sm">Revoke this API Key</div>
                <div className="text-sm text-muted-foreground">
                  This action cannot be undone. The{' '}
                  <b>{getAPIKeyResponse?.apiKey?.displayName}</b> API key will
                  no longer be usable.
                </div>
              </div>
              <AlertDialog open={revokeOpen} onOpenChange={setRevokeOpen}>
                <AlertDialogTrigger asChild>
                  <Button variant="destructive">Revoke API Key</Button>
                </AlertDialogTrigger>
                <AlertDialogContent>
                  <AlertDialogHeader>
                    <AlertDialogTitle>Are you sure?</AlertDialogTitle>
                    <p>
                      This action cannot be undone. The{' '}
                      <b>{getAPIKeyResponse?.apiKey?.displayName}</b> API key
                      will no longer be usable.
                    </p>
                  </AlertDialogHeader>
                  <AlertDialogFooter>
                    <AlertDialogCancel>Cancel</AlertDialogCancel>
                    <Button variant="destructive" onClick={handleRevoke}>
                      Revoke API Key
                    </Button>
                  </AlertDialogFooter>
                </AlertDialogContent>
              </AlertDialog>
            </div>
          ) : (
            <div className="flex justify-between items-center space-y-2">
              <div>
                <div className="font-semibold text-sm">Delete this API Key</div>
                <div className="text-sm text-muted-foreground">
                  This action cannot be undone. The{' '}
                  <b>{getAPIKeyResponse?.apiKey?.displayName}</b> API key will
                  be permanently deleted.
                </div>
              </div>
              <AlertDialog>
                <AlertDialogTrigger asChild>
                  <Button variant="destructive">Delete API Key</Button>
                </AlertDialogTrigger>
                <AlertDialogContent>
                  <AlertDialogHeader>
                    <AlertDialogTitle>Are you sure?</AlertDialogTitle>
                    <p>
                      This action cannot be undone.{' '}
                      <b>{getAPIKeyResponse?.apiKey?.displayName}</b> will be
                      permanently deleted.
                    </p>
                  </AlertDialogHeader>
                  <AlertDialogFooter>
                    <AlertDialogCancel>Cancel</AlertDialogCancel>
                    <Button variant="destructive" onClick={handleDelete}>
                      Delete API Key
                    </Button>
                  </AlertDialogFooter>
                </AlertDialogContent>
              </AlertDialog>
            </div>
          )}
        </div>
      </ConsoleCardContent>
    </ConsoleCard>
  );
};
