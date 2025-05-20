import React, { useEffect, useState } from 'react';
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import { useNavigate, useParams } from 'react-router';
import {
  createAPIKey,
  getProject,
  listAPIKeys,
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
import { timestampDate, timestampFromDate } from '@bufbuild/protobuf/wkt';
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
import { CalendarIcon, CirclePlus, Copy, LoaderCircle } from 'lucide-react';
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
} from '@/components/ui/form';
import { z } from 'zod';
import { useForm } from 'react-hook-form';
import { toast } from 'sonner';
import { Input } from '@/components/ui/input';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@/components/ui/popover';
import { cn } from '@/lib/utils';
import { Calendar } from '@/components/ui/calendar';
import { format, set } from 'date-fns';
import { Link } from 'react-router-dom';
import { APIKey } from '@/gen/tesseral/backend/v1/models_pb';
import { SecretCopier } from '@/components/SecretCopier';

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

  const apiKeys = listApiKeysResponses?.pages?.flatMap((page) => page.apiKeys);

  return (
    <Card>
      <CardHeader className="py-4 flex flex-row items-center justify-between">
        <div>
          <CardTitle>API Keys</CardTitle>
          <CardDescription>
            Manage the API keys for this organization.
          </CardDescription>
        </div>
        <CreateAPIKeyButton />
      </CardHeader>
      <CardContent>
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
                        {getProjectResponse?.project?.apiKeySecretTokenPrefix ||
                          'api_key_'}
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
      </CardContent>
    </Card>
  );
};

const schema = z.object({
  displayName: z.string(),
  expireTime: z.string(),
});

const CreateAPIKeyButton = () => {
  const [createOpen, setCreateOpen] = useState(false);
  const [secretOpen, setSecretOpen] = useState(false);
  const [apiKey, setApiKey] = useState<APIKey>();
  const navigate = useNavigate();

  const [customDate, setCustomDate] = useState<Date>();

  const { organizationId } = useParams();
  const { data: getProjectResponse } = useQuery(getProject);
  const { refetch } = useInfiniteQuery(
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

  const createApiKeyMutation = useMutation(createAPIKey);

  const form = useForm<z.infer<typeof schema>>({
    defaultValues: {
      displayName: '',
      expireTime: '1 day',
    },
  });

  const handleSubmit = async (data: z.infer<typeof schema>) => {
    const createParams: Record<string, any> = {
      organizationId: organizationId!,
      displayName: data.displayName,
    };

    switch (data.expireTime) {
      case '1 day':
        createParams.expireTime = timestampFromDate(
          new Date(Date.now() + 24 * 60 * 60 * 1000),
        );
        break;
      case '7 days':
        createParams.expireTime = timestampFromDate(
          new Date(Date.now() + 7 * 24 * 60 * 60 * 1000),
        );
        break;
      case '30 days':
        createParams.expireTime = timestampFromDate(
          new Date(Date.now() + 30 * 24 * 60 * 60 * 1000),
        );
        break;
      case 'custom':
        if (customDate) {
          createParams.expireTime = timestampFromDate(customDate);
        }
        break;
      case 'noexpire':
        break;
    }

    const { apiKey } = await createApiKeyMutation.mutateAsync({
      apiKey: createParams,
    });

    if (apiKey) {
      setApiKey(apiKey);
      setCreateOpen(false);
      setSecretOpen(true);

      toast.success('API Key created successfully');

      await refetch();
    }
  };

  return (
    <>
      <AlertDialog
        open={!!apiKey?.secretToken && secretOpen}
        onOpenChange={setSecretOpen}
      >
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>API Key Created</AlertDialogTitle>
            <AlertDialogDescription>
              API Key was created successfully.
            </AlertDialogDescription>
          </AlertDialogHeader>

          <div className="text-sm font-medium leading-none">
            API Key Secret Token
          </div>

          {apiKey?.secretToken && (
            <SecretCopier
              placeholder={`${getProjectResponse?.project?.apiKeySecretTokenPrefix}•••••••••••••••••••••••••••••••••••••••••••••••••••••••`}
              secret={apiKey.secretToken}
            />
          )}

          <div className="text-sm text-muted-foreground">
            Store this secret in your secrets manager. You will not be able to
            see this secret token again later.
          </div>

          <AlertDialogFooter>
            <AlertDialogCancel onClick={() => setSecretOpen(false)}>
              Close
            </AlertDialogCancel>
            {!!apiKey?.id && (
              <Link
                to={`/organizations/${organizationId}/api-keys/${apiKey.id}`}
              >
                <Button>View API Key</Button>
              </Link>
            )}
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      <AlertDialog open={createOpen} onOpenChange={setCreateOpen}>
        <AlertDialogTrigger asChild>
          <Button variant="outline">
            <CirclePlus className="h-4 w-4" />
            Create API Key
          </Button>
        </AlertDialogTrigger>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Create API Key</AlertDialogTitle>
          </AlertDialogHeader>

          <Form {...form}>
            <form
              onSubmit={form.handleSubmit(handleSubmit)}
              className="space-y-4"
            >
              <FormField
                control={form.control}
                name="displayName"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Display Name</FormLabel>
                    <FormControl>
                      <Input {...field} />
                    </FormControl>
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="expireTime"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Expire time</FormLabel>
                    <FormControl>
                      <div className="flex flex-row gap-2">
                        <Select
                          {...field}
                          onValueChange={(value) => {
                            field.onChange(value);

                            console.log(value);
                          }}
                        >
                          <SelectTrigger className="w-[180px]">
                            <SelectValue placeholder="Pick a custom date" />
                          </SelectTrigger>
                          <SelectContent>
                            <SelectItem value="1 day">1 day</SelectItem>
                            <SelectItem value="7 days">7 days</SelectItem>
                            <SelectItem value="30 days">30 days</SelectItem>
                            <SelectItem value="custom">Custom</SelectItem>
                            <SelectItem value="noexpire">
                              No expiration
                            </SelectItem>
                          </SelectContent>
                        </Select>

                        {field.value === 'custom' && (
                          <Popover>
                            <PopoverTrigger asChild>
                              <Button
                                variant={'outline'}
                                className={cn(
                                  'w-[270px] justify-start text-left font-normal',
                                  !customDate && 'text-muted-foreground',
                                )}
                              >
                                <CalendarIcon className="mr-2 h-4 w-4" />
                                {customDate ? (
                                  format(customDate, 'PPP')
                                ) : (
                                  <span>Pick a date</span>
                                )}
                              </Button>
                            </PopoverTrigger>
                            <PopoverContent className="w-auto p-0">
                              <Calendar
                                mode="single"
                                selected={customDate}
                                onSelect={setCustomDate}
                                initialFocus
                              />
                            </PopoverContent>
                          </Popover>
                        )}
                      </div>
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
    </>
  );
};
