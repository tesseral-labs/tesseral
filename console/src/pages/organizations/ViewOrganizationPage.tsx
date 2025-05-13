import { useMutation, useQuery } from '@connectrpc/connect-query';
import {
  getOrganization,
  updateOrganization,
} from '@/gen/tesseral/backend/v1/backend-BackendService_connectquery';
import { Outlet, useLocation, useParams } from 'react-router';
import React, { FC, useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { DateTime } from 'luxon';
import { timestampDate } from '@bufbuild/protobuf/wkt';
import { clsx } from 'clsx';
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from '@/components/ui/breadcrumb';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import {
  PageCodeSubtitle,
  PageContent,
  PageDescription,
  PageHeader,
  PageTitle,
} from '@/components/page';
import { z } from 'zod';
import { useForm } from 'react-hook-form';
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
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
} from '@/components/ui/form';
import { Input } from '@/components/ui/input';
import { toast } from 'sonner';
import { Organization } from '@/gen/tesseral/backend/v1/models_pb';
import { TabBar, TabBarLink } from '@/components/ui/tab-bar';

export const ViewOrganizationPage = () => {
  const { organizationId } = useParams();
  const { data: getOrganizationResponse, refetch } = useQuery(getOrganization, {
    id: organizationId,
  });
  const { pathname } = useLocation();

  const tabs = [
    {
      root: true,
      name: 'Details',
      url: `/organizations/${organizationId}`,
    },
    {
      name: 'Users',
      url: `/organizations/${organizationId}/users`,
    },
    {
      name: 'User Invites',
      url: `/organizations/${organizationId}/user-invites`,
    },
    {
      name: 'Roles',
      url: `/organizations/${organizationId}/roles`,
    },
    {
      name: 'SAML Connections',
      url: `/organizations/${organizationId}/saml-connections`,
    },
    {
      name: 'SCIM API Keys',
      url: `/organizations/${organizationId}/scim-api-keys`,
    },
  ];

  const currentTab = tabs.find((tab) => tab.url === pathname);

  return (
    // TODO remove padding when app shell in place
    <>
      <TabBar>
        {tabs.map((tab) => (
          <TabBarLink
            key={tab.name}
            active={tab.url === currentTab?.url}
            url={tab.url}
            label={tab.name}
          />
        ))}
      </TabBar>
      <PageHeader>
        <PageTitle>
          {getOrganizationResponse?.organization?.displayName}
        </PageTitle>
        <PageCodeSubtitle>{organizationId}</PageCodeSubtitle>
        <PageDescription>
          An Organization represents one of your business customers.
        </PageDescription>
      </PageHeader>
      <PageContent>
        {currentTab?.root && (
          <Card className="my-8">
            <CardHeader className="py-4 flex flex-row items-center justify-between">
              <div>
                <CardTitle className="text-xl">General configuration</CardTitle>
              </div>
              <EditOrganizationButton />
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-3 gap-x-2 text-sm">
                <div className="border-r border-gray-200 pr-8">
                  <div className="font-semibold">Display Name</div>
                  <div>
                    {getOrganizationResponse?.organization?.displayName}
                  </div>
                </div>
                <div className="border-r border-gray-200 pl-8 pr-8">
                  <div className="font-semibold">Created</div>
                  <div>
                    {getOrganizationResponse?.organization?.createTime &&
                      DateTime.fromJSDate(
                        timestampDate(
                          getOrganizationResponse.organization.createTime,
                        ),
                      ).toRelative()}
                  </div>
                </div>
                <div className="px-8">
                  <div className="font-semibold">Last updated</div>
                  <div>
                    {getOrganizationResponse?.organization?.updateTime &&
                      DateTime.fromJSDate(
                        timestampDate(
                          getOrganizationResponse.organization.updateTime,
                        ),
                      ).toRelative()}
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>
        )}

        <div className="mt-4">
          <Outlet />
        </div>
      </PageContent>
    </>
  );
};

const schema = z.object({
  displayName: z.string(),
});

const EditOrganizationButton: FC = () => {
  const { organizationId } = useParams();

  const form = useForm<z.infer<typeof schema>>({
    defaultValues: {
      displayName: '',
    },
  });

  const { data: getOrganizationResponse, refetch } = useQuery(getOrganization, {
    id: organizationId,
  });

  const updateOrganizationMutation = useMutation(updateOrganization);

  const [open, setOpen] = useState(false);

  const handleSubmit = async (data: z.infer<typeof schema>) => {
    const updatedOrganization: Partial<Organization> = {
      displayName: data.displayName,
    };

    await updateOrganizationMutation.mutateAsync({
      id: organizationId,
      organization: updatedOrganization as Organization,
    });

    await refetch();

    setOpen(false);

    toast.success('Organization updated successfully');
  };

  useEffect(() => {
    if (getOrganizationResponse?.organization) {
      form.reset({
        displayName: getOrganizationResponse.organization.displayName,
      });
    }
  }, [getOrganizationResponse]);

  return (
    <AlertDialog open={open} onOpenChange={setOpen}>
      <AlertDialogTrigger asChild>
        <Button variant="outline">Edit</Button>
      </AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Edit Organization</AlertDialogTitle>
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
