import { useMutation, useQuery } from '@connectrpc/connect-query';
import {
  createProjectAPIKey,
  createPublishableKey,
  listProjectAPIKeys,
  listPublishableKeys,
} from '@/gen/tesseral/backend/v1/backend-BackendService_connectquery';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { Link } from 'react-router-dom';
import React, { useState } from 'react';
import { DateTime } from 'luxon';
import { timestampDate } from '@bufbuild/protobuf/wkt';
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from '@/components/ui/breadcrumb';
import { PageDescription, PageTitle } from '@/components/page';
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { z } from 'zod';
import { useNavigate } from 'react-router';
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
import { SecretCopier } from '@/components/SecretCopier';
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

export const ListAPIKeysPage = () => {
  const { data: listPublishableKeysResponse } = useQuery(
    listPublishableKeys,
    {},
  );
  const { data: listProjectAPIKeysResponse } = useQuery(listProjectAPIKeys, {});

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
            <BreadcrumbLink asChild>
              <Link to="/project-settings">Project Settings</Link>
            </BreadcrumbLink>
          </BreadcrumbItem>
          <BreadcrumbSeparator />
          <BreadcrumbItem>
            <BreadcrumbPage>API Keys</BreadcrumbPage>
          </BreadcrumbItem>
        </BreadcrumbList>
      </Breadcrumb>

      <PageTitle>API Keys</PageTitle>
      <PageDescription className="mt-2">Lorem ipsum dolor.</PageDescription>

      <div className="mt-8 space-y-8">
        <Card>
          <CardHeader className="flex-row justify-between items-center">
            <div className="flex flex-col space-y-1 5">
              <CardTitle>Publishable Keys</CardTitle>
              <CardDescription>
                Tesseral's client-side SDKs require a publishable key.
                Publishable keys can be publicly accessible in your web or
                mobile app's client-side code. Lorem ipsum dolor.
              </CardDescription>
            </div>
            <CreatePublishableKeyButton />
          </CardHeader>
          <CardContent>
            <Table>
              <TableHeader>
                <TableRow>
                  <TableCell>Display Name</TableCell>
                  <TableHead>ID</TableHead>
                  <TableHead>Created</TableHead>
                  <TableHead>Updated</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {listPublishableKeysResponse?.publishableKeys?.map(
                  (publishableKey) => (
                    <TableRow key={publishableKey.id}>
                      <TableCell className="font-medium">
                        <Link
                          className="font-medium underline underline-offset-2 decoration-muted-foreground/40"
                          to={`/project-settings/api-keys/publishable-keys/${publishableKey.id}`}
                        >
                          {publishableKey.displayName}
                        </Link>
                      </TableCell>
                      <TableCell className="font-mono">
                        {publishableKey.id}
                      </TableCell>
                      <TableCell>
                        {publishableKey.createTime &&
                          DateTime.fromJSDate(
                            timestampDate(publishableKey.createTime),
                          ).toRelative()}
                      </TableCell>
                      <TableCell>
                        {publishableKey.updateTime &&
                          DateTime.fromJSDate(
                            timestampDate(publishableKey.updateTime),
                          ).toRelative()}
                      </TableCell>
                    </TableRow>
                  ),
                )}
              </TableBody>
            </Table>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex-row justify-between items-center">
            <div className="flex flex-col space-y-1 5">
              <CardTitle>Project API Keys</CardTitle>
              <CardDescription>
                Project API keys are how your backend can automate operations in
                Tesseral. Lorem ipsum dolor.
              </CardDescription>
            </div>
            <CreateProjectAPIKeyButton />
          </CardHeader>
          <CardContent>
            <Table>
              <TableHeader>
                <TableRow>
                  <TableCell>Display Name</TableCell>
                  <TableHead>ID</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead>Created At</TableHead>
                  <TableHead>Updated At</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {listProjectAPIKeysResponse?.projectApiKeys?.map(
                  (projectAPIKey) => (
                    <TableRow key={projectAPIKey.id}>
                      <TableCell className="font-medium">
                        <Link
                          className="font-medium underline underline-offset-2 decoration-muted-foreground/40"
                          to={`/project-settings/api-keys/project-api-keys/${projectAPIKey.id}`}
                        >
                          {projectAPIKey.displayName}
                        </Link>
                      </TableCell>
                      <TableCell className="font-mono">
                        {projectAPIKey.id}
                      </TableCell>
                      <TableCell>
                        {projectAPIKey?.revoked ? 'Revoked' : 'Active'}
                      </TableCell>
                      <TableCell>
                        {projectAPIKey.createTime &&
                          DateTime.fromJSDate(
                            timestampDate(projectAPIKey.createTime),
                          ).toRelative()}
                      </TableCell>
                      <TableCell>
                        {projectAPIKey.updateTime &&
                          DateTime.fromJSDate(
                            timestampDate(projectAPIKey.updateTime),
                          ).toRelative()}
                      </TableCell>
                    </TableRow>
                  ),
                )}
              </TableBody>
            </Table>
          </CardContent>
        </Card>
      </div>
    </div>
  );
};

const publishableKeySchema = z.object({
  displayName: z.string(),
});

const CreatePublishableKeyButton = () => {
  const createPublishableKeyMutation = useMutation(createPublishableKey);
  /* eslint-disable @typescript-eslint/no-unsafe-call */
  // Currently there's an issue with the types of react-hook-form and zod
  // preventing the compiler from inferring the correct types.
  const form = useForm<z.infer<typeof publishableKeySchema>>({
    resolver: zodResolver(publishableKeySchema),
    defaultValues: {
      displayName: '',
    },
  });
  /* eslint-enable @typescript-eslint/no-unsafe-call */
  const navigate = useNavigate();
  const [open, setOpen] = useState(false);

  const handleSubmit = async (values: z.infer<typeof projectApiKeySchema>) => {
    const { publishableKey } = await createPublishableKeyMutation.mutateAsync({
      publishableKey: {
        displayName: values.displayName,
      },
    });

    setOpen(false);
    navigate(
      `/project-settings/api-keys/publishable-keys/${publishableKey?.id}`,
    );
  };

  return (
    <>
      <AlertDialog open={open} onOpenChange={setOpen}>
        <AlertDialogTrigger>
          <Button variant="outline">Create</Button>
        </AlertDialogTrigger>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Create Publishable Key</AlertDialogTitle>
            <AlertDialogDescription>
              Tesseral's client-side SDKs require a publishable key. Publishable
              keys can be publicly accessible in your web or mobile app's
              client-side code. Lorem ipsum dolor.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <Form {...form}>
            {/* eslint-disable @typescript-eslint/no-unsafe-call */}
            {/** Currently there's an issue with the types of react-hook-form and zod
            preventing the compiler from inferring the correct types.*/}
            <form onSubmit={form.handleSubmit(handleSubmit)}>
              {/* eslint-enable @typescript-eslint/no-unsafe-call */}
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
                      A human-friendly name for the Publishable Key. You can
                      edit this later.
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <AlertDialogFooter className="mt-8">
                <AlertDialogCancel>Cancel</AlertDialogCancel>
                <Button type="submit">Create Publishable Key</Button>
              </AlertDialogFooter>
            </form>
          </Form>
        </AlertDialogContent>
      </AlertDialog>
    </>
  );
};

const projectApiKeySchema = z.object({
  displayName: z.string(),
});

const CreateProjectAPIKeyButton = () => {
  const createProjectAPIKeyMutation = useMutation(createProjectAPIKey);

  /* eslint-disable @typescript-eslint/no-unsafe-call */
  // Currently there's an issue with the types of react-hook-form and zod
  // preventing the compiler from inferring the correct types.
  const form = useForm<z.infer<typeof projectApiKeySchema>>({
    resolver: zodResolver(projectApiKeySchema),
    defaultValues: {
      displayName: '',
    },
  });
  /* eslint-enable @typescript-eslint/no-unsafe-call */
  const navigate = useNavigate();
  const [createOpen, setCreateOpen] = useState(false);
  const [projectAPIKeyID, setProjectAPIKeyID] = useState('');
  const [secretToken, setSecretToken] = useState('');

  const handleSubmit = async (values: z.infer<typeof projectApiKeySchema>) => {
    const { projectApiKey } = await createProjectAPIKeyMutation.mutateAsync({
      projectApiKey: {
        displayName: values.displayName,
      },
    });

    setCreateOpen(false);
    if (projectApiKey?.id) {
      setProjectAPIKeyID(projectApiKey.id);
    }
    if (projectApiKey?.secretToken) {
      setSecretToken(projectApiKey.secretToken);
    }
  };

  const handleClose = () => {
    navigate(`/project-settings/api-keys/project-api-keys/${projectAPIKeyID}`);
  };

  return (
    <>
      <AlertDialog open={!!secretToken}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Project API Key Created</AlertDialogTitle>
            <AlertDialogDescription>
              Project API Key was created successfully.
            </AlertDialogDescription>
          </AlertDialogHeader>

          <div className="text-sm font-medium leading-none">
            Project Secret Token
          </div>

          <SecretCopier
            placeholder="tesseral_secret_key_•••••••••••••••••••••••••"
            secret={secretToken}
          />

          <div className="text-sm text-muted-foreground">
            Store this secret as TESSERAL_API_KEY in your secrets manager. You
            will not be able to see this secret token again later.
          </div>

          <AlertDialogFooter>
            <AlertDialogCancel onClick={handleClose}>Close</AlertDialogCancel>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      <AlertDialog open={createOpen} onOpenChange={setCreateOpen}>
        <AlertDialogTrigger>
          <Button variant="outline">Create</Button>
        </AlertDialogTrigger>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Create Project API Key</AlertDialogTitle>
            <AlertDialogDescription>
              A Project API key is how your backend talks to the Tesseral
              Backend API. Lorem ipsum dolor.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <Form {...form}>
            {/* eslint-disable @typescript-eslint/no-unsafe-call */}
            {/** Currently there's an issue with the types of react-hook-form and zod
            preventing the compiler from inferring the correct types.*/}
            <form onSubmit={form.handleSubmit(handleSubmit)}>
              {/* eslint-enable @typescript-eslint/no-unsafe-call */}
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
                      A human-friendly name for the Project API Key.
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <AlertDialogFooter className="mt-8">
                <AlertDialogCancel>Cancel</AlertDialogCancel>
                <Button type="submit">Create Project API Key</Button>
              </AlertDialogFooter>
            </form>
          </Form>
        </AlertDialogContent>
      </AlertDialog>
    </>
  );
};
