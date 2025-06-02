import { useMutation, useQuery } from '@connectrpc/connect-query';
import {
  getProject,
  updateProject,
} from '@/gen/tesseral/backend/v1/backend-BackendService_connectquery';
import React, { useEffect, useState } from 'react';
import { DateTime } from 'luxon';
import { timestampDate } from '@bufbuild/protobuf/wkt';
import {
  PageCodeSubtitle,
  PageContent,
  PageDescription,
  PageHeader,
  PageTitle,
} from '@/components/page';
import {
  ConsoleCard,
  ConsoleCardContent,
  ConsoleCardDetails,
  ConsoleCardHeader,
  ConsoleCardTitle,
} from '@/components/ui/console-card';
import {
  DetailsGrid,
  DetailsGridColumn,
  DetailsGridEntry,
  DetailsGridKey,
  DetailsGridValue,
} from '@/components/details-grid';
import { Outlet, useLocation } from 'react-router';
import { TabBar, TabBarLink } from '@/components/ui/tab-bar';
import { Settings2 } from 'lucide-react';
import { z } from 'zod';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
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
import {
  Form,
  FormField,
  FormItem,
  FormLabel,
  FormDescription,
  FormControl,
  FormMessage,
} from '@/components/ui/form';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';

export const ViewProjectSettingsPage = () => {
  const { data: getProjectResponse } = useQuery(getProject, {});
  const { pathname } = useLocation();

  const tabs = [
    {
      root: true,
      name: 'Details',
      url: `/project-settings`,
    },
    {
      name: 'Login Settings',
      url: `/project-settings/login-settings`,
    },
    {
      name: 'Vault UI Settings',
      url: `/project-settings/vault-ui-settings`,
    },
    {
      name: 'Vault Domain Settings',
      url: `/project-settings/vault-domain-settings`,
    },
    {
      name: 'Role-Based Access Control Settings',
      url: `/project-settings/rbac-settings`,
    },
    {
      name: 'API Keys',
      url: `/project-settings/api-keys`,
    },
  ];

  const currentTab = tabs.find((tab) => tab.url === pathname);

  return (
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
        <PageTitle className="flex items-center">
          <Settings2 className="inline mr-2 w-6 h-6" />
          Project settings
        </PageTitle>
        <PageCodeSubtitle>{getProjectResponse?.project?.id}</PageCodeSubtitle>
        <PageDescription>
          Everything you do in Tesseral happens inside a Project.
        </PageDescription>
      </PageHeader>
      <PageContent>
        <ConsoleCard className="my-8">
          <ConsoleCardHeader>
            <ConsoleCardDetails>
              <ConsoleCardTitle>General configuration</ConsoleCardTitle>
            </ConsoleCardDetails>
            <EditButton />
          </ConsoleCardHeader>

          <ConsoleCardContent>
            <DetailsGrid>
              <DetailsGridColumn>
                <DetailsGridEntry>
                  <DetailsGridKey>Display name</DetailsGridKey>
                  <DetailsGridValue>
                    {getProjectResponse?.project?.displayName}
                  </DetailsGridValue>
                </DetailsGridEntry>
              </DetailsGridColumn>
              <DetailsGridColumn>
                <DetailsGridEntry>
                  <DetailsGridKey>Created</DetailsGridKey>
                  <DetailsGridValue>
                    {getProjectResponse?.project?.createTime &&
                      DateTime.fromJSDate(
                        timestampDate(getProjectResponse.project.createTime),
                      ).toRelative()}
                  </DetailsGridValue>
                </DetailsGridEntry>
              </DetailsGridColumn>
              <DetailsGridColumn>
                <DetailsGridEntry>
                  <DetailsGridKey>Updated</DetailsGridKey>
                  <DetailsGridValue>
                    {getProjectResponse?.project?.updateTime &&
                      DateTime.fromJSDate(
                        timestampDate(getProjectResponse.project.updateTime),
                      ).toRelative()}
                  </DetailsGridValue>
                </DetailsGridEntry>
              </DetailsGridColumn>
            </DetailsGrid>
          </ConsoleCardContent>
        </ConsoleCard>

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

const EditButton = () => {
  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
  });

  const { data: getProjectResponse, refetch } = useQuery(getProject);
  useEffect(() => {
    if (getProjectResponse?.project) {
      form.reset({
        displayName: getProjectResponse.project.displayName,
      });
    }
  }, [getProjectResponse]);

  const updateProjectMutation = useMutation(updateProject);
  const [open, setOpen] = useState(false);
  const handleSubmit = async (values: z.infer<typeof schema>) => {
    await updateProjectMutation.mutateAsync({
      project: {
        displayName: values.displayName,
      },
    });
    await refetch();
    toast.success('Project updated');
    setOpen(false);
  };

  return (
    <AlertDialog open={open} onOpenChange={setOpen}>
      <AlertDialogTrigger>
        <Button variant="outline">Edit</Button>
      </AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Edit Project</AlertDialogTitle>
        </AlertDialogHeader>
        <Form {...form}>
          <form
            className="space-y-4"
            onSubmit={form.handleSubmit(handleSubmit)}
          >
            <FormField
              control={form.control}
              name="displayName"
              render={({ field }: { field: any }) => (
                <FormItem>
                  <FormLabel>Display name</FormLabel>
                  <FormDescription>
                    A user-facing, human-readable display name for your project.
                  </FormDescription>
                  <FormControl>
                    <Input {...field} />
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
