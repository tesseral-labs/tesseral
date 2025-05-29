import { useMutation, useQuery } from '@connectrpc/connect-query';
import {
  createBackendAPIKey,
  createPublishableKey,
  createStripeCheckoutLink,
  getProjectEntitlements,
  listBackendAPIKeys,
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
  PageContent,
  PageDescription,
  PageHeader,
  PageTitle,
} from '@/components/page';
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
import { Switch } from '@/components/ui/switch';
import {
  ConsoleCardDetails,
  ConsoleCardTableContent,
} from '@/components/ui/console-card';

export const ListAPIKeysTab = () => {
  const { data: getProjectEntitlementsResponse } = useQuery(
    getProjectEntitlements,
    {},
  );
  const createStripeCheckoutLinkMutation = useMutation(
    createStripeCheckoutLink,
  );

  const handleUpgrade = async () => {
    const { url } = await createStripeCheckoutLinkMutation.mutateAsync({});
    window.location.href = url;
  };

  const { data: listPublishableKeysResponse } = useQuery(
    listPublishableKeys,
    {},
  );
  const { data: listBackendAPIKeysResponse } = useQuery(listBackendAPIKeys, {});

  return (
    <div className="mt-8 space-y-8">
      <Card>
        <CardHeader className="flex-row justify-between items-center">
          <ConsoleCardDetails>
            <CardTitle>Publishable Keys</CardTitle>
            <CardDescription>
              Tesseral's client-side SDKs require a publishable key. Publishable
              keys can be publicly accessible in your web or mobile app's
              client-side code.
            </CardDescription>
          </ConsoleCardDetails>
          <CreatePublishableKeyButton />
        </CardHeader>
        <ConsoleCardTableContent>
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
        </ConsoleCardTableContent>
      </Card>

      <Card>
        <CardHeader className="flex-row justify-between items-center">
          <ConsoleCardDetails>
            <CardTitle>Backend API Keys</CardTitle>
            <CardDescription>
              Backend API keys are how your backend can automate operations in
              Tesseral using the Tesseral Backend API.
            </CardDescription>
          </ConsoleCardDetails>
          <CreateBackendAPIKeyButton />
        </CardHeader>
        <ConsoleCardTableContent>
          {/* do not treat undefined as unentitled, to avoid flickering here */}
          {getProjectEntitlementsResponse?.entitledBackendApiKeys === false ? (
            <div className="text-sm my-8 w-full flex flex-col items-center justify-center space-y-6">
              <div className="font-medium">
                Backend API Keys are available on the Growth Tier.
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
          ) : (
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
                {listBackendAPIKeysResponse?.backendApiKeys?.map(
                  (backendApiKey) => (
                    <TableRow key={backendApiKey.id}>
                      <TableCell className="font-medium">
                        <Link
                          className="font-medium underline underline-offset-2 decoration-muted-foreground/40"
                          to={`/project-settings/api-keys/backend-api-keys/${backendApiKey.id}`}
                        >
                          {backendApiKey.displayName}
                        </Link>
                      </TableCell>
                      <TableCell className="font-mono">
                        {backendApiKey.id}
                      </TableCell>
                      <TableCell>
                        {backendApiKey?.revoked ? 'Revoked' : 'Active'}
                      </TableCell>
                      <TableCell>
                        {backendApiKey.createTime &&
                          DateTime.fromJSDate(
                            timestampDate(backendApiKey.createTime),
                          ).toRelative()}
                      </TableCell>
                      <TableCell>
                        {backendApiKey.updateTime &&
                          DateTime.fromJSDate(
                            timestampDate(backendApiKey.updateTime),
                          ).toRelative()}
                      </TableCell>
                    </TableRow>
                  ),
                )}
              </TableBody>
            </Table>
          )}
        </ConsoleCardTableContent>
      </Card>
    </div>
  );
};

const publishableKeySchema = z.object({
  displayName: z.string().nonempty(),
  devMode: z.boolean(),
});

const CreatePublishableKeyButton = () => {
  const createPublishableKeyMutation = useMutation(createPublishableKey);

  // Currently there's an issue with the types of react-hook-form and zod
  // preventing the compiler from inferring the correct types.
  const form = useForm<z.infer<typeof publishableKeySchema>>({
    resolver: zodResolver(publishableKeySchema),
    defaultValues: {
      displayName: '',
      devMode: false,
    },
  });

  const navigate = useNavigate();
  const [open, setOpen] = useState(false);

  const handleSubmit = async (values: z.infer<typeof publishableKeySchema>) => {
    const { publishableKey } = await createPublishableKeyMutation.mutateAsync({
      publishableKey: {
        displayName: values.displayName,
        devMode: values.devMode,
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
              client-side code.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <Form {...form}>
            {}
            {/** Currently there's an issue with the types of react-hook-form and zod
             preventing the compiler from inferring the correct types.*/}
            <form
              onSubmit={form.handleSubmit(handleSubmit)}
              className="space-y-8"
            >
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

              <FormField
                control={form.control}
                name="devMode"
                render={({ field }: { field: any }) => (
                  <FormItem>
                    <FormLabel>Dev Mode</FormLabel>
                    <FormControl>
                      <Switch
                        className="block"
                        checked={field.value}
                        onCheckedChange={field.onChange}
                      />
                    </FormControl>
                    <FormDescription>
                      Enable this if you want to use this publishable key from
                      localhost.
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

const backendAPIKeySchema = z.object({
  displayName: z.string(),
});

const CreateBackendAPIKeyButton = () => {
  const { data: getProjectEntitlementsResponse } = useQuery(
    getProjectEntitlements,
    {},
  );

  const createBackendAPIKeyMutation = useMutation(createBackendAPIKey);

  // Currently there's an issue with the types of react-hook-form and zod
  // preventing the compiler from inferring the correct types.
  const form = useForm<z.infer<typeof backendAPIKeySchema>>({
    resolver: zodResolver(backendAPIKeySchema),
    defaultValues: {
      displayName: '',
    },
  });

  const navigate = useNavigate();
  const [createOpen, setCreateOpen] = useState(false);
  const [backendAPIKeyID, setBackendAPIKeyID] = useState('');
  const [secretToken, setSecretToken] = useState('');

  const handleSubmit = async (values: z.infer<typeof backendAPIKeySchema>) => {
    const { backendApiKey } = await createBackendAPIKeyMutation.mutateAsync({
      backendApiKey: {
        displayName: values.displayName,
      },
    });

    setCreateOpen(false);
    setBackendAPIKeyID(backendApiKey!.id);
    setSecretToken(backendApiKey!.secretToken);
  };

  const handleClose = () => {
    navigate(`/project-settings/api-keys/backend-api-keys/${backendAPIKeyID}`);
  };

  return (
    <>
      <AlertDialog open={!!secretToken}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Backend API Key Created</AlertDialogTitle>
            <AlertDialogDescription>
              Backend API Key was created successfully.
            </AlertDialogDescription>
          </AlertDialogHeader>

          <div className="text-sm font-medium leading-none">
            Backend API Key Secret Token
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
        <AlertDialogTrigger
          disabled={!getProjectEntitlementsResponse?.entitledBackendApiKeys}
        >
          <Button
            variant="outline"
            disabled={!getProjectEntitlementsResponse?.entitledBackendApiKeys}
          >
            Create
          </Button>
        </AlertDialogTrigger>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Create Backend API Key</AlertDialogTitle>
            <AlertDialogDescription>
              Backend API keys are how your backend can automate operations in
              Tesseral using the Tesseral Backend API.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <Form {...form}>
            {}
            {/** Currently there's an issue with the types of react-hook-form and zod
             preventing the compiler from inferring the correct types.*/}
            <form onSubmit={form.handleSubmit(handleSubmit)}>
              {}
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
                      A human-friendly name for the Backend API Key.
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <AlertDialogFooter className="mt-8">
                <AlertDialogCancel>Cancel</AlertDialogCancel>
                <Button type="submit">Create Backend API Key</Button>
              </AlertDialogFooter>
            </form>
          </Form>
        </AlertDialogContent>
      </AlertDialog>
    </>
  );
};
