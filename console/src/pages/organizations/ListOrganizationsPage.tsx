import React, { FC, useState } from 'react';
import { useInfiniteQuery, useMutation } from '@connectrpc/connect-query';
import {
  createOrganization,
  listOrganizations,
} from '@/gen/tesseral/backend/v1/backend-BackendService_connectquery';
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
import { Link, useNavigate } from 'react-router-dom';
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from '@/components/ui/breadcrumb';
import { PageDescription, PageTitle } from '@/components/page';
import { Button } from '@/components/ui/button';
import { CirclePlus, LoaderCircleIcon } from 'lucide-react';
import { z } from 'zod';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import {
  Form,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
} from '@/components/ui/form';
import { Input } from '@/components/ui/input';
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
import { toast } from 'sonner';

export const ListOrganizationsPage = () => {
  const {
    data: listOrganizationsResponses,
    fetchNextPage,
    hasNextPage,
    isFetchingNextPage,
    refetch,
  } = useInfiniteQuery(
    listOrganizations,
    {
      pageToken: '',
    },
    {
      pageParamKey: 'pageToken',
      getNextPageParam: (page) => page.nextPageToken || undefined,
    },
  );

  const organizations = listOrganizationsResponses?.pages?.flatMap(
    (page) => page.organizations,
  );

  return (
    <div>
      <Breadcrumb>
        <BreadcrumbList>
          <BreadcrumbItem>
            <BreadcrumbLink asChild>
              <Link to="/">Home</Link>
            </BreadcrumbLink>
          </BreadcrumbItem>
          <BreadcrumbSeparator />
          <BreadcrumbItem>
            <BreadcrumbPage>Organizations</BreadcrumbPage>
          </BreadcrumbItem>
        </BreadcrumbList>
      </Breadcrumb>

      <PageTitle>Organizations</PageTitle>
      <PageDescription>
        An Organization represents one of your business customers.
      </PageDescription>

      <Card className="mt-8 overflow-hidden">
        <CardHeader className="flex flex-row items-center space-x-4 justify-between">
          <div className="flex flex-col space-y-1 5">
            <CardTitle>Organizations list</CardTitle>
            <CardDescription>
              This is a list of all Organizations in your project. You can
              create and edit these Organizations manually.
            </CardDescription>
          </div>
          <CreateOrganizationButton />
        </CardHeader>
        <CardContent className="-m-6 mt-0">
          <Table>
            <TableHeader className="bg-gray-50">
              <TableRow>
                <TableHead>Display Name</TableHead>
                <TableHead>ID</TableHead>
                <TableHead>Created</TableHead>
                <TableHead>Updated</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {organizations?.map((org) => (
                <TableRow key={org.id}>
                  <TableCell className="font-medium">
                    <Link
                      className="underline underline-offset-2 decoration-muted-foreground/40"
                      to={`/organizations/${org.id}`}
                    >
                      {org.displayName}
                    </Link>
                  </TableCell>
                  <TableCell className="font-mono">{org.id}</TableCell>
                  <TableCell>
                    {org.createTime &&
                      DateTime.fromJSDate(
                        timestampDate(org.createTime),
                      ).toRelative()}
                  </TableCell>
                  <TableCell>
                    {org.updateTime &&
                      DateTime.fromJSDate(
                        timestampDate(org.updateTime),
                      ).toRelative()}
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </CardContent>
      </Card>

      {hasNextPage && (
        <Button
          className="mt-4"
          variant="outline"
          onClick={() => fetchNextPage()}
        >
          {isFetchingNextPage && (
            <LoaderCircleIcon className="h-4 w-4 animate-spin" />
          )}
          Load more
        </Button>
      )}
    </div>
  );
};

const CreateOrganizationButton: FC = () => {
  const navigate = useNavigate();

  const form = useForm<z.infer<typeof organizationSchema>>({
    resolver: zodResolver(organizationSchema),
    defaultValues: {
      displayName: '',
    },
  });

  const createOrganizationMutation = useMutation(createOrganization);

  const [open, setOpen] = useState(false);

  const handleSubmit = async (values: z.infer<typeof organizationSchema>) => {
    const createOrganizationResponse =
      await createOrganizationMutation.mutateAsync({
        organization: {
          ...values,
        },
      });
    toast.success('Organization created successfully');

    setOpen(false);

    navigate(`/organizations/${createOrganizationResponse.organization?.id}`);
  };

  return (
    <AlertDialog open={open} onOpenChange={setOpen}>
      <AlertDialogTrigger asChild>
        <Button variant="outline">
          <CirclePlus />
          Create Organization
        </Button>
      </AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Create Organization</AlertDialogTitle>
          <AlertDialogDescription>
            Create a new Organization.
          </AlertDialogDescription>
        </AlertDialogHeader>

        <Form {...form}>
          <form
            className="space-y-4"
            onSubmit={form.handleSubmit(handleSubmit)}
          >
            <FormField
              control={form.control}
              name="displayName"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Display Name</FormLabel>
                  <Input placeholder="ACME Corp" {...field} />
                  <FormDescription>
                    The display name of the Organization. This will be displayed
                    to users during the login process.
                  </FormDescription>
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

const organizationSchema = z.object({
  displayName: z.string(),
});
