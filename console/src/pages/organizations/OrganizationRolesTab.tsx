import { useParams } from 'react-router';
import {
  useInfiniteQuery,
  useMutation,
  useQuery,
} from '@connectrpc/connect-query';
import {
  getOrganization,
  listRoles,
  updateOrganization,
} from '@/gen/tesseral/backend/v1/backend-BackendService_connectquery';
import {
  ConsoleCard,
  ConsoleCardDetails,
  ConsoleCardContent,
  ConsoleCardDescription,
  ConsoleCardHeader,
  ConsoleCardTitle,
  ConsoleCardTableContent,
} from '@/components/ui/console-card';
import { Button, buttonVariants } from '@/components/ui/button';
import { Link } from 'react-router-dom';
import {
  DetailsGrid,
  DetailsGridColumn,
  DetailsGridEntry,
  DetailsGridKey,
  DetailsGridValue,
} from '@/components/details-grid';
import React, { useEffect, useState } from 'react';
import { LoaderCircleIcon } from 'lucide-react';
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
import { z } from 'zod';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
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
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form';
import { Switch } from '@/components/ui/switch';
import { toast } from 'sonner';

export function OrganizationRolesTab() {
  const { organizationId } = useParams();
  const { data: getOrganizationResponse } = useQuery(getOrganization, {
    id: organizationId,
  });
  const {
    data: listRolesResponses,
    fetchNextPage,
    hasNextPage,
    isFetchingNextPage,
  } = useInfiniteQuery(
    listRoles,
    {
      organizationId,
      pageToken: '',
    },
    {
      pageParamKey: 'pageToken',
      getNextPageParam: (page) => page.nextPageToken || undefined,
    },
  );

  const roles = listRolesResponses?.pages?.flatMap((page) => page.roles);

  return (
    <div className="space-y-8">
      <ConsoleCard>
        <ConsoleCardHeader className="flex-row justify-between items-center">
          <ConsoleCardDetails>
            <ConsoleCardTitle>Organization Role Settings</ConsoleCardTitle>
            <ConsoleCardDescription>
              Role-related settings for this Organization.
            </ConsoleCardDescription>
          </ConsoleCardDetails>
          <EditRolesSettingsButton />
        </ConsoleCardHeader>
        <ConsoleCardContent>
          <DetailsGrid>
            <DetailsGridColumn>
              <DetailsGridEntry>
                <DetailsGridKey>Custom Roles</DetailsGridKey>
                <DetailsGridValue>
                  {getOrganizationResponse?.organization?.customRolesEnabled
                    ? 'Enabled'
                    : 'Disabled'}
                </DetailsGridValue>
              </DetailsGridEntry>
            </DetailsGridColumn>
          </DetailsGrid>
        </ConsoleCardContent>
      </ConsoleCard>

      <ConsoleCard>
        <ConsoleCardHeader className="flex-row justify-between items-center">
          <ConsoleCardDetails>
            <ConsoleCardTitle>Organization-Specific Roles</ConsoleCardTitle>
            <ConsoleCardDescription>
              Roles that exist only for this Organization.
            </ConsoleCardDescription>
          </ConsoleCardDetails>

          <Link
            to={`/roles/new?organization-id=${organizationId}`}
            className={buttonVariants({ variant: 'outline' })}
          >
            Create Role
          </Link>
        </ConsoleCardHeader>
        <ConsoleCardTableContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Display Name</TableHead>
                <TableHead>Actions</TableHead>
                <TableHead>Created</TableHead>
                <TableHead>Updated</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {roles?.map((role) => (
                <TableRow key={role.id}>
                  <TableCell>
                    <Link
                      to={`/roles/${role.id}`}
                      className="font-medium underline underline-offset-2 decoration-muted-foreground/40"
                    >
                      {role.displayName}
                    </Link>
                  </TableCell>
                  <TableCell className="font-mono">
                    {role.actions.join(' ')}
                  </TableCell>
                  <TableCell>
                    {role?.createTime &&
                      DateTime.fromJSDate(
                        timestampDate(role.createTime),
                      ).toRelative()}
                  </TableCell>
                  <TableCell>
                    {role?.updateTime &&
                      DateTime.fromJSDate(
                        timestampDate(role.updateTime),
                      ).toRelative()}
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>

          {hasNextPage && (
            <Button
              className="mt-4 mb-6"
              variant="outline"
              onClick={() => fetchNextPage()}
            >
              {isFetchingNextPage && (
                <LoaderCircleIcon className="h-4 w-4 animate-spin" />
              )}
              Load more
            </Button>
          )}
        </ConsoleCardTableContent>
      </ConsoleCard>
    </div>
  );
}

const schema = z.object({
  customRolesEnabled: z.boolean(),
});

function EditRolesSettingsButton() {
  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      customRolesEnabled: false,
    },
  });

  const { organizationId } = useParams();
  const { data: getOrganizationResponse, refetch } = useQuery(getOrganization, {
    id: organizationId,
  });

  useEffect(() => {
    if (getOrganizationResponse?.organization) {
      form.reset({
        customRolesEnabled:
          getOrganizationResponse.organization.customRolesEnabled,
      });
    }
  }, [form, getOrganizationResponse]);

  const { mutateAsync: updateOrganizationAsync } =
    useMutation(updateOrganization);
  const handleSubmit = async (values: z.infer<typeof schema>) => {
    await updateOrganizationAsync({
      id: organizationId,
      organization: {
        customRolesEnabled: values.customRolesEnabled,
      },
    });

    await refetch();
    toast.success('Organization Roles Settings updated');
    setOpen(false);
  };

  const [open, setOpen] = useState(false);
  return (
    <AlertDialog open={open} onOpenChange={setOpen}>
      <AlertDialogTrigger>
        <Button variant="outline">Edit</Button>
      </AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Edit Organization Roles Settings</AlertDialogTitle>
          <AlertDialogDescription>
            Edit Roles-related settings for this Organization.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <Form {...form}>
          <form
            className="space-y-4"
            onSubmit={form.handleSubmit(handleSubmit)}
          >
            <FormField
              control={form.control}
              name="customRolesEnabled"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Custom Roles</FormLabel>
                  <FormControl>
                    <Switch
                      className="block"
                      checked={field.value}
                      onCheckedChange={field.onChange}
                    />
                  </FormControl>
                  <FormDescription>
                    Whether this Organization can create Organization-Specific
                    Roles.
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
}
