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
  deleteSCIMAPIKey,
  getOrganization,
  getProject,
  getSCIMAPIKey,
  revokeSCIMAPIKey,
  updateSCIMAPIKey,
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

export const ViewSCIMAPIKeyPage = () => {
  const { organizationId, scimApiKeyId } = useParams();
  const { data: getProjectResponse } = useQuery(getProject, {});
  const { data: getOrganizationResponse } = useQuery(getOrganization, {
    id: organizationId,
  });
  const { data: getScimApiKeyResponse } = useQuery(getSCIMAPIKey, {
    id: scimApiKeyId,
  });
  return (
    // TODO remove padding when app shell in place
    <>
      <PageHeader>
        <PageTitle>{getScimApiKeyResponse?.scimApiKey?.displayName}</PageTitle>
        <PageCodeSubtitle>{scimApiKeyId}</PageCodeSubtitle>
        <PageDescription>
          A SCIM API key lets your customer do enterprise directory syncing.
        </PageDescription>
      </PageHeader>

      <PageContent>
        <Card className="my-8">
          <CardHeader className="flex-row justify-between items-center">
            <div className="flex flex-col space-y-1 5">
              <CardTitle>Configuration</CardTitle>
              <CardDescription>
                Details about this SCIM API Key.
              </CardDescription>
            </div>
            <EditSCIMAPIKeyButton />
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-2 gap-x-2 text-sm">
              <div className="border-r border-gray-200 pr-8 flex flex-col gap-4">
                <div>
                  <div className="font-semibold">Display Name</div>
                  <div className="truncate">
                    {getScimApiKeyResponse?.scimApiKey?.displayName}
                  </div>
                </div>
                <div>
                  <div className="font-semibold">SCIM Base URL</div>
                  <div className="truncate">
                    {`https://${getProjectResponse?.project?.vaultDomain}/api/scim/v1`}
                  </div>
                </div>
                <div>
                  <div className="font-semibold">Revoked</div>
                  <div className="truncate">
                    {getScimApiKeyResponse?.scimApiKey?.revoked ? 'Yes' : 'No'}
                  </div>
                </div>
              </div>
              <div className="border-gray-200 pl-8 flex flex-col gap-4">
                <div>
                  <div className="font-semibold">Created</div>
                  <div className="truncate">
                    {getScimApiKeyResponse?.scimApiKey?.createTime &&
                      DateTime.fromJSDate(
                        timestampDate(
                          getScimApiKeyResponse?.scimApiKey?.createTime,
                        ),
                      ).toRelative()}
                  </div>
                </div>
                <div>
                  <div className="font-semibold">Updated</div>
                  <div className="truncate">
                    {getScimApiKeyResponse?.scimApiKey?.updateTime &&
                      DateTime.fromJSDate(
                        timestampDate(
                          getScimApiKeyResponse?.scimApiKey?.updateTime,
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

const EditSCIMAPIKeyButton = () => {
  const { scimApiKeyId } = useParams();
  const { data: getSCIMAPIKeyResponse, refetch } = useQuery(getSCIMAPIKey, {
    id: scimApiKeyId,
  });
  const updateSCIMAPIKeyMutation = useMutation(updateSCIMAPIKey);

  // Currently there's an issue with the types of react-hook-form and zod
  // preventing the compiler from inferring the correct types.
  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      displayName: '',
    },
  });
  useEffect(() => {
    if (getSCIMAPIKeyResponse?.scimApiKey) {
      form.reset({
        displayName: getSCIMAPIKeyResponse.scimApiKey.displayName,
      });
    }
  }, [getSCIMAPIKeyResponse]);

  const [open, setOpen] = useState(false);

  const handleSubmit = async (values: z.infer<typeof schema>) => {
    await updateSCIMAPIKeyMutation.mutateAsync({
      id: scimApiKeyId,
      scimApiKey: {
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
          <AlertDialogTitle>Edit SCIM API Key</AlertDialogTitle>
          <AlertDialogDescription>
            Edit SCIM API Key settings.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <Form {...form}>
          {/** Currently there's an issue with the types of react-hook-form and zod
           preventing the compiler from inferring the correct types.*/}
          <form onSubmit={form.handleSubmit(handleSubmit)}>
            <FormField
              control={form.control}
              name="displayName"
              render={({ field }: { field: any }) => (
                <FormItem>
                  <FormLabel>Display Name</FormLabel>
                  <FormControl>
                    <Input className="max-w-96" {...field} />
                  </FormControl>
                  <FormDescription>
                    A human-friendly name for the SCIM API Key.
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
  );
};

const DangerZoneCard = () => {
  const { organizationId, scimApiKeyId } = useParams();
  const { data: getScimApiKeyResponse, refetch } = useQuery(getSCIMAPIKey, {
    id: scimApiKeyId,
  });

  const [confirmRevokeOpen, setConfirmRevokeOpen] = useState(false);
  const handleRevoke = () => {
    setConfirmRevokeOpen(true);
  };

  const revokeScimApiKeyMutation = useMutation(revokeSCIMAPIKey);
  const handleConfirmRevoke = async () => {
    await revokeScimApiKeyMutation.mutateAsync({
      id: scimApiKeyId,
    });

    await refetch();
    toast.success('SCIM API Key revoked');
    setConfirmRevokeOpen(false);
  };

  const [confirmDeleteOpen, setConfirmDeleteOpen] = useState(false);

  const handleDelete = () => {
    setConfirmDeleteOpen(true);
  };

  const deleteScimApiKeyMutation = useMutation(deleteSCIMAPIKey);
  const navigate = useNavigate();
  const handleConfirmDelete = async () => {
    await deleteScimApiKeyMutation.mutateAsync({
      id: scimApiKeyId,
    });

    toast.success('SCIM API Key deleted');
    navigate(`/organizations/${organizationId}/scim-api-keys`);
  };

  return (
    <>
      <AlertDialog open={confirmRevokeOpen} onOpenChange={setConfirmRevokeOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>
              Revoke {getScimApiKeyResponse?.scimApiKey?.displayName}?
            </AlertDialogTitle>
            <AlertDialogDescription>
              Revoking a SCIM API Key cannot be undone. SCIM API calls from{' '}
              {getScimApiKeyResponse?.scimApiKey?.displayName} will stop
              working. This cannot be undone.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <Button variant="destructive" onClick={handleConfirmRevoke}>
              Revoke SCIM API Key
            </Button>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      <AlertDialog open={confirmDeleteOpen} onOpenChange={setConfirmDeleteOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>
              Delete {getScimApiKeyResponse?.scimApiKey?.displayName}?
            </AlertDialogTitle>
            <AlertDialogDescription>
              Deleting a SCIM API Key cannot be undone.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <Button variant="destructive" onClick={handleConfirmDelete}>
              Delete SCIM API Key
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
              <div className="text-sm font-semibold">Revoke SCIM API Key</div>
              <p className="text-sm">
                Revoke this SCIM API Key. SCIM API calls from this key will stop
                working. This cannot be undone.
              </p>
            </div>

            <Button
              variant="destructive"
              disabled={getScimApiKeyResponse?.scimApiKey?.revoked}
              onClick={handleRevoke}
            >
              Revoke SCIM API Key
            </Button>
          </div>

          <div className="flex justify-between items-center">
            <div>
              <div className="text-sm font-semibold">Delete SCIM API Key</div>
              <p className="text-sm">
                Delete this SCIM API Key. You must revoke this SCIM API Key
                first.
              </p>
            </div>

            <Button
              variant="destructive"
              disabled={!getScimApiKeyResponse?.scimApiKey?.revoked}
              onClick={handleDelete}
            >
              Delete SCIM API Key
            </Button>
          </div>
        </CardContent>
      </Card>
    </>
  );
};
