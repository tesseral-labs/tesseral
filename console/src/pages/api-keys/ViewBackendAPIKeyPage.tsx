import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from '@/components/ui/breadcrumb';
import { Link } from 'react-router-dom';
import React, { useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router';
import { useMutation, useQuery } from '@connectrpc/connect-query';
import {
  deleteBackendAPIKey,
  getBackendAPIKey,
  revokeBackendAPIKey,
  updateBackendAPIKey,
} from '@/gen/tesseral/backend/v1/backend-BackendService_connectquery';
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import { DateTime } from 'luxon';
import { timestampDate } from '@bufbuild/protobuf/wkt';
import { toast } from 'sonner';
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
import { z } from 'zod';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form';
import { Input } from '@/components/ui/input';
import {
  PageCodeSubtitle,
  PageContent,
  PageDescription,
  PageHeader,
  PageTitle,
} from '@/components/page';
import { ConsoleCardDetails } from '@/components/ui/console-card';

export const ViewBackendAPIKeyPage = () => {
  const { backendApiKeyId } = useParams();
  const { data: getBackendApiKeyResponse } = useQuery(getBackendAPIKey, {
    id: backendApiKeyId,
  });
  return (
    <>
      <PageHeader>
        <PageTitle>
          {getBackendApiKeyResponse?.backendApiKey?.displayName}
        </PageTitle>
        <PageCodeSubtitle>{backendApiKeyId}</PageCodeSubtitle>
        <PageDescription>
          Backend API keys are how your backend can automate operations in
          Tesseral using the Tesseral Backend API.
        </PageDescription>
      </PageHeader>

      <PageContent>
        <Card className="my-8">
          <CardHeader className="flex-row justify-between items-center">
            <ConsoleCardDetails>
              <CardTitle>Configuration</CardTitle>
              <CardDescription>Details about your Backend API.</CardDescription>
            </ConsoleCardDetails>
            <EditBackendAPIKeyButton />
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-3 gap-x-2 text-sm">
              <div className="border-r border-gray-200 pr-8 flex flex-col gap-4">
                <div>
                  <div className="font-semibold">Display Name</div>
                  <div className="truncate">
                    {getBackendApiKeyResponse?.backendApiKey?.displayName}
                  </div>
                </div>
                <div>
                  <div className="font-semibold">Revoked</div>
                  <div className="truncate">
                    {getBackendApiKeyResponse?.backendApiKey?.revoked
                      ? 'Yes'
                      : 'No'}
                  </div>
                </div>
              </div>
              <div className="border-r border-gray-200 pr-8 pl-8 flex flex-col gap-4">
                <div>
                  <div className="font-semibold">Created</div>
                  <div className="truncate">
                    {getBackendApiKeyResponse?.backendApiKey?.createTime &&
                      DateTime.fromJSDate(
                        timestampDate(
                          getBackendApiKeyResponse?.backendApiKey?.createTime,
                        ),
                      ).toRelative()}
                  </div>
                </div>
              </div>
              <div className="border-gray-200 pl-8 flex flex-col gap-4">
                <div>
                  <div className="font-semibold">Updated</div>
                  <div className="truncate">
                    {getBackendApiKeyResponse?.backendApiKey?.updateTime &&
                      DateTime.fromJSDate(
                        timestampDate(
                          getBackendApiKeyResponse?.backendApiKey?.updateTime,
                        ),
                      ).toRelative()}
                  </div>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>

        <DangerZoneCard />
      </PageContent>
    </>
  );
};

const schema = z.object({
  displayName: z.string(),
});

const EditBackendAPIKeyButton = () => {
  const { backendApiKeyId } = useParams();
  const { data: getBackendAPIKeyResponse, refetch } = useQuery(
    getBackendAPIKey,
    {
      id: backendApiKeyId,
    },
  );
  const updateBackendAPIKeyMutation = useMutation(updateBackendAPIKey);

  // Currently there's an issue with the types of react-hook-form and zod
  // preventing the compiler from inferring the correct types.
  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      displayName: '',
    },
  });
  useEffect(() => {
    if (getBackendAPIKeyResponse?.backendApiKey) {
      form.reset({
        displayName: getBackendAPIKeyResponse.backendApiKey.displayName,
      });
    }
  }, [getBackendAPIKeyResponse]);

  const [open, setOpen] = useState(false);

  const handleSubmit = async (values: z.infer<typeof schema>) => {
    await updateBackendAPIKeyMutation.mutateAsync({
      id: backendApiKeyId,
      backendApiKey: {
        displayName: values.displayName,
      },
    });
    await refetch();
    setOpen(false);
  };

  return (
    <AlertDialog open={open} onOpenChange={setOpen}>
      <AlertDialogTrigger>
        <Button variant="outline">Edit</Button>
      </AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Edit Backend API Key</AlertDialogTitle>
          <AlertDialogDescription>
            Edit Backend API Key settings.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <Form {...form}>
          {}
          {/** Currently there's an issue with the types of react-hook-form and zod
          preventing the compiler from inferring the correct types.*/}
          <form onSubmit={form.handleSubmit(handleSubmit)}>
            {}
            <FormField
              control={form.control}
              name="displayName"
              render={({ field }: { field: any }) => (
                <FormItem>
                  <FormLabel>Display Name</FormLabel>
                  <FormDescription>
                    A human-friendly name for the Backend API Key.
                  </FormDescription>
                  <FormControl>
                    <Input className="max-w-96" {...field} />
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

const DangerZoneCard = () => {
  const { backendApiKeyId } = useParams();
  const { data: getBackendApiKeyResponse, refetch } = useQuery(
    getBackendAPIKey,
    {
      id: backendApiKeyId,
    },
  );

  const [confirmRevokeOpen, setConfirmRevokeOpen] = useState(false);
  const handleRevoke = () => {
    setConfirmRevokeOpen(true);
  };

  const revokeBackendApiKeyMutation = useMutation(revokeBackendAPIKey);
  const handleConfirmRevoke = async () => {
    await revokeBackendApiKeyMutation.mutateAsync({
      id: backendApiKeyId,
    });

    await refetch();
    toast.success('Backend API Key revoked');
    setConfirmRevokeOpen(false);
  };

  const [confirmDeleteOpen, setConfirmDeleteOpen] = useState(false);

  const handleDelete = () => {
    setConfirmDeleteOpen(true);
  };

  const deleteBackendApiKeyMutation = useMutation(deleteBackendAPIKey);
  const navigate = useNavigate();
  const handleConfirmDelete = async () => {
    await deleteBackendApiKeyMutation.mutateAsync({
      id: backendApiKeyId,
    });

    toast.success('Backend API Key deleted');
    navigate(`/project-settings/api-keys`);
  };

  return (
    <>
      <AlertDialog open={confirmRevokeOpen} onOpenChange={setConfirmRevokeOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>
              Revoke {getBackendApiKeyResponse?.backendApiKey?.displayName}?
            </AlertDialogTitle>
            <AlertDialogDescription>
              Backend API calls from
              {getBackendApiKeyResponse?.backendApiKey?.displayName} will stop
              working. This cannot be undone.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <Button variant="destructive" onClick={handleConfirmRevoke}>
              Revoke Backend API Key
            </Button>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      <AlertDialog open={confirmDeleteOpen} onOpenChange={setConfirmDeleteOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>
              Delete {getBackendApiKeyResponse?.backendApiKey?.displayName}?
            </AlertDialogTitle>
            <AlertDialogDescription>
              Deleting a Backend API Key cannot be undone.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <Button variant="destructive" onClick={handleConfirmDelete}>
              Delete Backend API Key
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
                Revoke Backend API Key
              </div>
              <p className="text-sm">
                Revoke this Backend API Key. Backend API calls from this key
                will stop working. This cannot be undone.
              </p>
            </div>

            <Button
              variant="destructive"
              disabled={getBackendApiKeyResponse?.backendApiKey?.revoked}
              onClick={handleRevoke}
            >
              Revoke Backend API Key
            </Button>
          </div>

          <div className="flex justify-between items-center">
            <div>
              <div className="text-sm font-semibold">
                Delete Backend API Key
              </div>
              <p className="text-sm">
                Delete this Backend API Key. You must revoke this Backend API
                Key first.
              </p>
            </div>

            <Button
              variant="destructive"
              disabled={!getBackendApiKeyResponse?.backendApiKey?.revoked}
              onClick={handleDelete}
            >
              Delete Backend API Key
            </Button>
          </div>
        </CardContent>
      </Card>
    </>
  );
};
