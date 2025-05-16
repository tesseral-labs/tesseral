import React, { useEffect } from 'react';
import {
  PageCodeSubtitle,
  PageContent,
  PageDescription,
  PageHeader,
  PageTitle,
} from '@/components/page';
import {
  getAPIKey,
  listAPIKeyRoleAssignments,
  listRoles,
} from '@/gen/tesseral/backend/v1/backend-BackendService_connectquery';
import { useQuery } from '@connectrpc/connect-query';
import { useParams } from 'react-router';
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import {
  AlertDialog,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from '@/components/ui/alert-dialog';
import { Button } from '@/components/ui/button';
import { Form, useForm } from 'react-hook-form';
import { z } from 'zod';
import {
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
import { CirclePlus } from 'lucide-react';
import {
  Table,
  TableBody,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { Label } from '@/components/ui/label';

export const ViewAPIKeyPage = () => {
  const { organizationId, apiKeyId } = useParams();

  const { data: getAPIKeyResponse } = useQuery(getAPIKey, {
    id: apiKeyId,
    organizationId,
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
          <Card>
            <CardHeader className="py-4 flex flex-row items-center justify-between">
              <div>
                <CardTitle>API Key Details</CardTitle>
                <CardDescription></CardDescription>
              </div>

              <EditAPIKeyButton />
            </CardHeader>
            <CardContent>
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
                        <>{'never'}</>
                      )}
                    </DetailsGridValue>
                  </DetailsGridEntry>
                </DetailsGridColumn>
              </DetailsGrid>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="py-4 flex flex-row items-center justify-between">
              <div>
                <CardTitle>API Key Roles</CardTitle>
                <CardDescription>
                  Manage the roles associated with this API key.
                </CardDescription>
              </div>

              <AddRoleButton />
            </CardHeader>
            <CardContent>
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
                    (role) => (
                      <TableRow key={role.id}>
                        <TableHead>{role.id}</TableHead>
                        <TableHead>
                          {role.role?.actions.map((action) => (
                            <span
                              key={action}
                              className="p-1 text-xs text-mono bg-muted text-muted-foreground rounded"
                            >
                              {action}
                            </span>
                          ))}
                        </TableHead>
                        <TableHead>
                          <RemoveRoleButton />
                        </TableHead>
                      </TableRow>
                    ),
                  )}
                </TableBody>
              </Table>
            </CardContent>
          </Card>
        </div>
      </PageContent>
    </>
  );
};

const schema = z.object({
  displayName: z.string().min(1, { message: 'Display name is required' }),
});

const EditAPIKeyButton = () => {
  const { organizationId, apiKeyId } = useParams();
  const { data: getAPIKeyResponse } = useQuery(getAPIKey, {
    id: apiKeyId,
    organizationId,
  });

  const form = useForm<z.infer<typeof schema>>({
    defaultValues: {
      displayName: '',
    },
  });

  const handleSubmit = async (data: z.infer<typeof schema>) => {};

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

const AddRoleButton = () => {
  const { organizationId } = useParams();
  const { data: listRolesResponse } = useQuery(listRoles, {
    organizationId,
  });

  return (
    <AlertDialog>
      <AlertDialogTrigger asChild>
        <Button variant="outline">
          <CirclePlus className="h-4 w-4" />
          Add Role
        </Button>
      </AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Add Role</AlertDialogTitle>
        </AlertDialogHeader>

        <form>
          {listRolesResponse?.roles?.map((role) => (
            <div key={role.id}>
              <Label>{role.displayName}</Label>
              <Input type="checkbox" value={role.id} />
            </div>
          ))}

          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <Button type="submit">Save</Button>
          </AlertDialogFooter>
        </form>
      </AlertDialogContent>
    </AlertDialog>
  );
};

const RemoveRoleButton = () => {
  return (
    <AlertDialog>
      <AlertDialogTrigger asChild>
        <Button variant="destructive">Remove</Button>
      </AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Remove Role</AlertDialogTitle>
        </AlertDialogHeader>

        <AlertDialogFooter>
          <AlertDialogCancel>Cancel</AlertDialogCancel>
          <Button type="submit">Remove</Button>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
};
