import React, { useState } from 'react';
import { useNavigate, useParams } from 'react-router';
import { useMutation, useQuery } from '@connectrpc/connect-query';
import {
  createSCIMAPIKey,
  listSCIMAPIKeys,
} from '@/gen/tesseral/backend/v1/backend-BackendService_connectquery';
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { Link } from 'react-router-dom';
import { DateTime } from 'luxon';
import { timestampDate } from '@bufbuild/protobuf/wkt';
import { Button } from '@/components/ui/button';
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
import { z } from 'zod';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
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
import { SecretCopier } from '@/components/SecretCopier';

export const OrganizationSCIMAPIKeysTab = () => {
  const { organizationId } = useParams();
  const { data: listSCIMAPIKeysResponse } = useQuery(listSCIMAPIKeys, {
    organizationId,
  });

  return (
    <Card>
      <CardHeader className="flex-row justify-between items-center">
        <div className="flex flex-col space-y-1 5">
          <CardTitle>SCIM API Keys</CardTitle>
          <CardDescription>
            A SCIM API key lets your customer do enterprise directory syncing.
            Lorem ipsum dolor.
          </CardDescription>
        </div>
        <CreateSCIMAPIKeyButton />
      </CardHeader>
      <CardContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Display Name</TableHead>
              <TableHead>ID</TableHead>
              <TableHead>Created At</TableHead>
              <TableHead>Updated At</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {listSCIMAPIKeysResponse?.scimApiKeys?.map((scimAPIKey) => (
              <TableRow key={scimAPIKey.id}>
                <TableCell>
                  <Link
                    className="font-medium underline underline-offset-2 decoration-muted-foreground/40"
                    to={`/organizations/${organizationId}/scim-api-keys/${scimAPIKey.id}`}
                  >
                    {scimAPIKey.displayName}
                  </Link>
                </TableCell>
                <TableCell className="font-mono">{scimAPIKey.id}</TableCell>
                <TableCell>
                  {scimAPIKey.createTime &&
                    DateTime.fromJSDate(
                      timestampDate(scimAPIKey.createTime),
                    ).toRelative()}
                </TableCell>
                <TableCell>
                  {scimAPIKey.updateTime &&
                    DateTime.fromJSDate(
                      timestampDate(scimAPIKey.updateTime),
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
});

const CreateSCIMAPIKeyButton = () => {
  const { organizationId } = useParams();
  const createSCIMAPIKeyMutation = useMutation(createSCIMAPIKey);

  /* eslint-disable @typescript-eslint/no-unsafe-call */
  // Currently there's an issue with the types of react-hook-form and zod
  // preventing the compiler from inferring the correct types.
  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      displayName: '',
    },
  });
  /* eslint-enable @typescript-eslint/no-unsafe-call */

  const navigate = useNavigate();
  const [createOpen, setCreateOpen] = useState(false);
  const [scimAPIKeyID, setScimAPIKeyID] = useState('');
  const [secretToken, setSecretToken] = useState('');

  const handleSubmit = async (values: z.infer<typeof schema>) => {
    const { scimApiKey } = await createSCIMAPIKeyMutation.mutateAsync({
      scimApiKey: {
        organizationId,
        displayName: values.displayName,
      },
    });

    setCreateOpen(false);
    if (scimApiKey?.id) {
      setScimAPIKeyID(scimApiKey.id);
    }
    if (scimApiKey?.secretToken) {
      setSecretToken(scimApiKey.secretToken);
    }
  };

  const handleClose = () => {
    navigate(`/organizations/${organizationId}/scim-api-keys/${scimAPIKeyID}`);
  };

  return (
    <>
      <AlertDialog open={!!secretToken}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>SCIM API Key Created</AlertDialogTitle>
            <AlertDialogDescription>
              SCIM API Key was created successfully.
            </AlertDialogDescription>
          </AlertDialogHeader>

          <div className="text-sm font-medium leading-none">
            SCIM Secret Bearer Token
          </div>

          <SecretCopier
            placeholder="tesseral_secret_scim_api_key_•••••••••••••••••••••••••"
            secret={secretToken}
          />

          <div className="text-sm text-muted-foreground">
            Give this secret to your customer's IT admin. They will input it
            into their Identity Provider. You will not be able to see this
            secret token again later.
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
            <AlertDialogTitle>Create SCIM API Key</AlertDialogTitle>
            <AlertDialogDescription>
              A SCIM API key lets your customer do enterprise directory syncing.
              Lorem ipsum dolor.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <Form {...form}>
            {/* eslint-disable @typescript-eslint/no-unsafe-call */}
            {/** Currently
            there's an issue with the types of react-hook-form and zod
            preventing the compiler from inferring the correct types.
            */}
            <form onSubmit={form.handleSubmit(handleSubmit)}>
              {/** eslint-enable @typescript-eslint/no-unsafe-call */}
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
                      A human-friendly name for the SCIM API Key.
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <AlertDialogFooter className="mt-8">
                <AlertDialogCancel>Cancel</AlertDialogCancel>
                <Button type="submit">Create SCIM API Key</Button>
              </AlertDialogFooter>
            </form>
          </Form>
        </AlertDialogContent>
      </AlertDialog>
    </>
  );
};
