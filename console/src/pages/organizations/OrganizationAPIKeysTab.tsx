import React, { useEffect, useState } from 'react';
import {
  ConsoleCard,
  ConsoleCardContent,
  ConsoleCardDescription,
  ConsoleCardDetails,
  ConsoleCardHeader,
  ConsoleCardTableContent,
  ConsoleCardTitle,
} from '@/components/ui/console-card';
import { useParams } from 'react-router';
import {
  createStripeCheckoutLink,
  getOrganization,
  getProject,
  getProjectEntitlements,
  listAPIKeys,
  updateOrganization,
} from '@/gen/tesseral/backend/v1/backend-BackendService_connectquery';
import {
  useInfiniteQuery,
  useMutation,
  useQuery,
} from '@connectrpc/connect-query';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { DateTime } from 'luxon';
import { timestampDate } from '@bufbuild/protobuf/wkt';
import {
  AlertDialog,
  AlertDialogContent,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTrigger,
  AlertDialogCancel,
  AlertDialogTitle,
  AlertDialogDescription,
} from '@/components/ui/alert-dialog';
import { Button } from '@/components/ui/button';
import { LoaderCircle } from 'lucide-react';
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
} from '@/components/ui/form';
import { z } from 'zod';
import { useForm } from 'react-hook-form';
import { toast } from 'sonner';
import { Link } from 'react-router-dom';
import { zodResolver } from '@hookform/resolvers/zod';
import { Switch } from '@/components/ui/switch';
import { CreateAPIKeyButton } from './api-keys/CreateAPIKeyButton';
import {
  DetailsGrid,
  DetailsGridColumn,
  DetailsGridEntry,
  DetailsGridKey,
  DetailsGridValue,
} from '@/components/details-grid';

export const OrganizationAPIKeysTab = () => {
  const { organizationId } = useParams();
  const {
    data: listApiKeysResponses,
    fetchNextPage,
    hasNextPage,
    isFetchingNextPage,
  } = useInfiniteQuery(
    listAPIKeys,
    {
      organizationId,
      pageToken: '',
    },
    {
      pageParamKey: 'pageToken',
      getNextPageParam: (page) => page.nextPageToken || undefined,
    },
  );
  const { data: getProjectResponse } = useQuery(getProject);
  const { data: getOrganizationResponse } = useQuery(getOrganization, {
    id: organizationId,
  });
  const { data: getProjectEntitlementsResponse } = useQuery(
    getProjectEntitlements,
  );
  const createStripeCheckoutLinkMutation = useMutation(
    createStripeCheckoutLink,
  );

  const apiKeys = listApiKeysResponses?.pages?.flatMap((page) => page.apiKeys);

  async function handleUpgrade() {
    const { url } = await createStripeCheckoutLinkMutation.mutateAsync({});
    window.location.href = url;
  }

  return (
    <div className="space-y-8">
      {!getProjectEntitlementsResponse?.entitledBackendApiKeys ? (
        <ConsoleCard>
          <ConsoleCardHeader>
            <ConsoleCardDetails>
              <ConsoleCardTitle>API Key Management</ConsoleCardTitle>
              <ConsoleCardDescription>
                API keys are used to authenticate requests to your service. You
                can create and manage API keys for this Organization.
              </ConsoleCardDescription>
            </ConsoleCardDetails>
          </ConsoleCardHeader>
          <ConsoleCardContent>
            <div className="text-sm my-8 w-full flex flex-col items-center justify-center space-y-6">
              <div className="font-medium">
                API Keys are available on the Growth Tier.
              </div>

              <div className="flex items-center gap-x-4">
                <Button onClick={handleUpgrade}>Upgrade to Growth Tier</Button>
                <span>
                  or{' '}
                  <a
                    href="https://cal.com/ned-o-leary-j8ydyi/30min"
                    className="font-medium underline underline-offset-2 decoration-muted-foreground/40"
                  >
                    meet an expert
                  </a>
                  .
                </span>
              </div>
            </div>
          </ConsoleCardContent>
        </ConsoleCard>
      ) : (
        <>
          <ConsoleCard>
            <ConsoleCardHeader className="py-4 flex flex-row items-center justify-between">
              <div className="flex flex-col space-y-1 5">
                <ConsoleCardTitle>API Key Management</ConsoleCardTitle>
                <ConsoleCardDescription>
                  API keys are used to authenticate requests to your service.
                  You can create and manage API keys for this Organization.
                </ConsoleCardDescription>
              </div>
              <EditAPIKeySettingsButton />
            </ConsoleCardHeader>
            <ConsoleCardContent>
              <DetailsGrid>
                <DetailsGridColumn>
                  <DetailsGridEntry>
                    <DetailsGridKey>Status</DetailsGridKey>
                    <DetailsGridValue>
                      {getOrganizationResponse?.organization?.apiKeysEnabled
                        ? 'Enabled'
                        : 'Disabled'}
                    </DetailsGridValue>
                  </DetailsGridEntry>
                </DetailsGridColumn>
              </DetailsGrid>
            </ConsoleCardContent>
          </ConsoleCard>
          <ConsoleCard>
            <ConsoleCardHeader className="py-4 flex flex-row items-center justify-between">
              <div className="flex flex-col space-y-1 5">
                <ConsoleCardTitle>API Keys</ConsoleCardTitle>
                <ConsoleCardDescription>
                  Manage the API keys for this organization.
                </ConsoleCardDescription>
              </div>
              <CreateAPIKeyButton />
            </ConsoleCardHeader>
            <ConsoleCardTableContent>
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Display Name</TableHead>
                    <TableHead>ID</TableHead>
                    <TableHead>Value</TableHead>
                    <TableHead>Status</TableHead>
                    <TableHead>Expires</TableHead>
                    <TableHead>Created At</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {apiKeys &&
                    apiKeys.map((apiKey) => (
                      <TableRow key={apiKey.id}>
                        <TableCell>
                          <Link
                            to={`/organizations/${organizationId}/api-keys/${apiKey.id}`}
                          >
                            {apiKey.displayName}
                          </Link>
                        </TableCell>
                        <TableCell>
                          <Link
                            to={`/organizations/${organizationId}/api-keys/${apiKey.id}`}
                          >
                            {apiKey.id}
                          </Link>
                        </TableCell>
                        <TableCell>
                          {apiKey.secretTokenSuffix ? (
                            <span className="font-mono text-sm">
                              {getProjectResponse?.project
                                ?.apiKeySecretTokenPrefix || 'api_key_'}
                              ...{apiKey.secretTokenSuffix}
                            </span>
                          ) : (
                            '—'
                          )}
                        </TableCell>
                        <TableCell>
                          {apiKey.revoked ? (
                            <span>Active</span>
                          ) : (
                            <span>Revoked</span>
                          )}
                        </TableCell>
                        <TableCell>
                          {apiKey.expireTime
                            ? DateTime.fromJSDate(
                                timestampDate(apiKey.expireTime),
                              ).toRelative()
                            : 'Never'}
                        </TableCell>
                        <TableCell>
                          {apiKey.createTime &&
                            DateTime.fromJSDate(
                              timestampDate(apiKey.createTime),
                            ).toRelative()}
                        </TableCell>
                      </TableRow>
                    ))}
                </TableBody>
              </Table>

              {hasNextPage && (
                <div className="flex justify-center mt-8">
                  <Button
                    className="mt-4"
                    variant="outline"
                    onClick={() => fetchNextPage()}
                  >
                    {isFetchingNextPage && (
                      <LoaderCircle className="h-4 w-4 animate-spin" />
                    )}
                    Load more
                  </Button>
                </div>
              )}
            </ConsoleCardTableContent>
          </ConsoleCard>
        </>
      )}
    </div>
  );
};

const schema = z.object({
  apiKeysEnabled: z.boolean(),
});

function EditAPIKeySettingsButton() {
  const { organizationId } = useParams();
  const [open, setOpen] = useState(false);

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      apiKeysEnabled: false,
    },
  });

  const { data: getOrganizationResponse, refetch } = useQuery(getOrganization, {
    id: organizationId,
  });
  const updateOrganizationMutation = useMutation(updateOrganization);

  async function handleSubmit(data: z.infer<typeof schema>) {
    await updateOrganizationMutation.mutateAsync({
      id: organizationId,
      organization: {
        apiKeysEnabled: data.apiKeysEnabled,
      },
    });

    await refetch();

    toast.success('API key settings updated successfully');
    setOpen(false);
  }

  useEffect(() => {
    if (getOrganizationResponse?.organization) {
      form.reset({
        apiKeysEnabled:
          getOrganizationResponse?.organization?.apiKeysEnabled || false,
      });
    }
  }, [getOrganizationResponse]);

  return (
    <AlertDialog open={open} onOpenChange={setOpen}>
      <AlertDialogTrigger asChild>
        <Button variant="outline">Edit API Key Settings</Button>
      </AlertDialogTrigger>

      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Edit API Key Settings</AlertDialogTitle>
          <AlertDialogDescription>
            Manage the API key settings for this Organization.
          </AlertDialogDescription>
        </AlertDialogHeader>

        <Form {...form}>
          <form
            onSubmit={form.handleSubmit(handleSubmit)}
            className="space-y-4"
          >
            <FormField
              control={form.control}
              name="apiKeysEnabled"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>API Keys Enabled</FormLabel>
                  <FormControl>
                    <Switch
                      className="block"
                      checked={field.value}
                      onCheckedChange={field.onChange}
                    />
                  </FormControl>
                  <FormDescription>
                    Whether this Organization can use API keys to authenticate
                    to your service.
                  </FormDescription>
                </FormItem>
              )}
            />

            <AlertDialogFooter>
              <AlertDialogCancel onClick={() => setOpen(false)}>
                Cancel
              </AlertDialogCancel>
              <Button type="submit">Save</Button>
            </AlertDialogFooter>
          </form>
        </Form>
      </AlertDialogContent>
    </AlertDialog>
  );
}
