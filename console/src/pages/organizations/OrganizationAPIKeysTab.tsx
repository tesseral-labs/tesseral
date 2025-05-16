import React, { useState } from 'react';
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
import { useMutation, useQuery } from '@connectrpc/connect-query';
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
} from '@/components/ui/alert-dialog';
import { Button } from '@/components/ui/button';
import { CalendarIcon, CirclePlus } from 'lucide-react';
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
import { format } from 'date-fns';
import { Link } from 'react-router-dom';

export const OrganizationAPIKeysTab = () => {
  const { organizationId } = useParams();
  const { data: listApiKeysResponse } = useQuery(listAPIKeys, {
    organizationId,
  });
  const { data: getProjectResponse } = useQuery(getProject);

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
              <TableHead>ID</TableHead>
              <TableHead>Value</TableHead>
              <TableHead>Created At</TableHead>
              <TableHead>Updated At</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {listApiKeysResponse?.apiKeys.map((apiKey) => (
              <TableRow key={apiKey.id}>
                <TableCell>
                  <Link
                    to={`/organizations/${organizationId}/api-keys/${apiKey.id}`}
                  >
                    {apiKey.id}
                  </Link>
                </TableCell>
                <TableCell>
                  <span className="font-mono text-sm">
                    {getProjectResponse?.project?.apiKeySecretTokenPrefix ||
                      'api_key_'}
                    ...{apiKey.secretTokenSuffix}
                  </span>
                </TableCell>
                <TableCell>
                  {apiKey.createTime &&
                    DateTime.fromJSDate(
                      timestampDate(apiKey.createTime),
                    ).toRelative()}
                </TableCell>
                <TableCell>
                  {apiKey.updateTime &&
                    DateTime.fromJSDate(
                      timestampDate(apiKey.updateTime),
                    ).toRelative()}
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  );
};

const schema = z.object({
  displayName: z.string(),
  expireTime: z.string(),
});

const CreateAPIKeyButton = () => {
  const navigate = useNavigate();

  const [customDate, setCustomDate] = useState<Date>();

  const { organizationId } = useParams();
  const { refetch } = useQuery(listAPIKeys, {
    organizationId,
  });

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

    // if (data.expireTime === 'custom') {
    //   createParams.expireTime = customDate?.toISOString();
    // } else if (data.expireTime !== 'noexpire') {
    //   createParams.expireTime = data.expireTime;
    // } else {
    // }

    const { apiKey } = await createApiKeyMutation.mutateAsync(createParams);

    if (apiKey) {
      navigate(`/organizations/${organizationId}/api-keys/${apiKey.id}`);

      toast.success('API Key created successfully');
    }
  };

  return (
    <AlertDialog>
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
  );
};
