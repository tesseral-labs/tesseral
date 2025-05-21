import { SecretCopier } from '@/components/SecretCopier';
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
import { Calendar } from '@/components/ui/calendar';
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
} from '@/components/ui/form';
import { Input } from '@/components/ui/input';
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@/components/ui/popover';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import {
  createAPIKey,
  getProject,
  listAPIKeys,
} from '@/gen/tesseral/backend/v1/backend-BackendService_connectquery';
import { APIKey } from '@/gen/tesseral/backend/v1/models_pb';
import { cn } from '@/lib/utils';
import { timestampFromDate } from '@bufbuild/protobuf/wkt';
import {
  useInfiniteQuery,
  useMutation,
  useQuery,
} from '@connectrpc/connect-query';
import { format } from 'date-fns';
import { CalendarIcon, CirclePlus } from 'lucide-react';
import React, { useState } from 'react';
import { useForm } from 'react-hook-form';
import { useNavigate, useParams } from 'react-router';
import { Link } from 'react-router-dom';
import { toast } from 'sonner';
import { z } from 'zod';

const schema = z.object({
  displayName: z.string(),
  expireTime: z.string(),
});

export function CreateAPIKeyButton() {
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
}
