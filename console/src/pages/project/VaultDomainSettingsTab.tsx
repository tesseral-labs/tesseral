import React, { useEffect, useState } from 'react';
import {
  getProject,
  getVaultDomainSettings,
  updateVaultDomainSettings,
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
  TableHeader,
  TableRow,
  TableCell,
  TableHead,
  TableBody,
} from '@/components/ui/table';
import { useMutation, useQuery } from '@connectrpc/connect-query';
import {
  DetailsGrid,
  DetailsGridColumn,
  DetailsGridEntry,
  DetailsGridKey,
  DetailsGridValue,
} from '@/components/details-grid';
import { Badge } from '@/components/ui/badge';
import { CheckIcon, XIcon } from 'lucide-react';
import { VaultDomainSettingsDNSRecord } from '@/gen/tesseral/backend/v1/models_pb';
import { Button } from '@/components/ui/button';
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from '@/components/ui/tooltip';
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
import { Input } from '@/components/ui/input';
import Loader from '@/components/ui/loader';

export const VaultDomainSettingsTab = () => {
  const { data: getProjectResponse } = useQuery(getProject, {});
  const { data: getVaultDomainSettingsResponse } = useQuery(
    getVaultDomainSettings,
  );

  return (
    <div className="space-y-8">
      <Card>
        <CardHeader className="flex-row justify-between items-center">
          <div className="flex flex-col space-y-1 5">
            <CardTitle>Vault Domain Settings</CardTitle>
            <CardDescription>
              Configure a custom domain for your Vault.
            </CardDescription>
          </div>
          <EditCustomAuthDomainButton />
        </CardHeader>
        <CardContent>
          <DetailsGrid>
            <DetailsGridColumn>
              <DetailsGridEntry>
                <DetailsGridKey>Current Vault Domain</DetailsGridKey>
                <DetailsGridValue>
                  {getProjectResponse?.project?.vaultDomain}
                </DetailsGridValue>
              </DetailsGridEntry>
              <DetailsGridEntry>
                <DetailsGridKey>Current Email Send-From Domain</DetailsGridKey>
                <DetailsGridValue>
                  {getProjectResponse?.project?.emailSendFromDomain}
                </DetailsGridValue>
              </DetailsGridEntry>
            </DetailsGridColumn>
            <DetailsGridColumn>
              <DetailsGridEntry>
                <DetailsGridKey>Pending Custom Domain</DetailsGridKey>
                <DetailsGridValue>
                  {getVaultDomainSettingsResponse?.vaultDomainSettings
                    ?.pendingDomain || '-'}
                </DetailsGridValue>
              </DetailsGridEntry>
            </DetailsGridColumn>
          </DetailsGrid>
        </CardContent>
      </Card>

      {getVaultDomainSettingsResponse?.vaultDomainSettings?.pendingDomain && (
        <>
          <Card>
            <CardHeader className="flex-row justify-between items-center">
              <div className="flex flex-col space-y-1 5">
                <CardTitle>Vault Domain Records</CardTitle>
                <CardDescription>
                  You need to add the following DNS records before you can use{' '}
                  <span className="font-medium">
                    {
                      getVaultDomainSettingsResponse?.vaultDomainSettings
                        ?.pendingDomain
                    }
                  </span>{' '}
                  as your Vault domain.
                </CardDescription>
              </div>

              {getVaultDomainSettingsResponse?.vaultDomainSettings
                ?.pendingSendFromDomainReady ? (
                <Button variant="outline">Enable Custom Vault Domain</Button>
              ) : (
                <TooltipProvider>
                  <Tooltip>
                    <TooltipTrigger asChild>
                      <Button
                        className="disabled:pointer-events-auto"
                        variant="outline"
                        disabled
                      >
                        Enable Custom Vault Domain
                      </Button>
                    </TooltipTrigger>
                    <TooltipContent className="max-w-96">
                      Your DNS records need to be correct and widely propagated
                      before you can enable this.
                    </TooltipContent>
                  </Tooltip>
                </TooltipProvider>
              )}
            </CardHeader>
            <CardContent>
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Type</TableHead>
                    <TableHead>Name</TableHead>
                    <TableHead>Value</TableHead>
                    <TableHead>Status</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {getVaultDomainSettingsResponse?.vaultDomainSettings?.vaultDomainRecords?.map(
                    (record, i) => <DNSRecordRows key={i} record={record} />,
                  )}
                </TableBody>
              </Table>
            </CardContent>
          </Card>
          <Card className="mt-8">
            <CardHeader className="flex-row justify-between items-center">
              <div className="flex flex-col space-y-1 5">
                <CardTitle>
                  Email Send-From Records
                  <Badge className="ml-4" variant="outline">
                    Optional
                  </Badge>
                </CardTitle>
                <CardDescription>
                  You can optionally configure the emails Tesseral sends to your
                  end users to come from{' '}
                  <span className="font-medium">{`noreply@mail.${getVaultDomainSettingsResponse?.vaultDomainSettings?.pendingDomain}`}</span>
                  , instead of the Tesseral-provided{' '}
                  <span className="font-medium">noreply@mail.tesseral.app</span>
                  .
                </CardDescription>
              </div>

              {getVaultDomainSettingsResponse?.vaultDomainSettings
                ?.pendingSendFromDomainReady ? (
                <Button variant="outline">
                  Enable Custom Email Send-From Domain
                </Button>
              ) : (
                <TooltipProvider>
                  <Tooltip>
                    <TooltipTrigger asChild>
                      <Button
                        className="disabled:pointer-events-auto"
                        variant="outline"
                        disabled
                      >
                        Enable Custom Email Send-From Domain
                      </Button>
                    </TooltipTrigger>
                    <TooltipContent className="max-w-96">
                      Your DNS records need to be correct and widely propagated
                      before you can enable this.
                    </TooltipContent>
                  </Tooltip>
                </TooltipProvider>
              )}
            </CardHeader>
            <CardContent>
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Type</TableHead>
                    <TableHead>Name</TableHead>
                    <TableHead>Value</TableHead>
                    <TableHead>Status</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {getVaultDomainSettingsResponse?.vaultDomainSettings?.emailSendFromRecords?.map(
                    (record, i) => <DNSRecordRows key={i} record={record} />,
                  )}
                </TableBody>
              </Table>
            </CardContent>
          </Card>
        </>
      )}
    </div>
  );
};

const DNSRecordRows = ({
  record,
}: {
  record: VaultDomainSettingsDNSRecord;
}) => {
  const noValue = (record.actualValues ?? []).length === 0;
  const tooManyValues = record.actualValues?.length > 1;
  const incorrectValue =
    record.actualValues?.length === 1 &&
    record.actualValues[0] !== record.wantValue;

  return (
    <>
      <TableRow>
        <TableCell>{record.type}</TableCell>
        <TableCell>{record.name}</TableCell>
        <TableCell>{record.wantValue}</TableCell>
        <TableCell>
          {record.correct ? (
            <Badge variant="default">
              <CheckIcon className="w-4 h-4" />
            </Badge>
          ) : (
            <Badge variant="destructive">
              <XIcon className="w-4 h-4" />
            </Badge>
          )}
        </TableCell>
      </TableRow>

      {/*{!record.correct && noValue && (*/}
      {/*  <TableRow>*/}
      {/*    <TableCell colSpan={3} className="bg-red-100 text-red-500 text-xs">*/}
      {/*      <div className="mx-4">*/}
      {/*        You haven't configured a{' '}*/}
      {/*        <span className="font-medium">{record.type}</span> record with the*/}
      {/*        name <span className="font-medium">{record.name}</span>. If you*/}
      {/*        recently created that record, it may still be propagating.*/}
      {/*      </div>*/}
      {/*    </TableCell>*/}
      {/*  </TableRow>*/}
      {/*)}*/}

      {/*{!record.correct && incorrectValue && (*/}
      {/*  <TableRow>*/}
      {/*    <TableCell colSpan={3} className="bg-red-100 text-red-500 text-xs">*/}
      {/*      <div className="mx-4">*/}
      {/*        <p>You created this record, but it has the wrong value.</p>*/}
      {/*        <p>*/}
      {/*          Your record has the value: <span className="font-medium">{record.actualValues[0]}</span>*/}
      {/*        </p>*/}
      {/*        <p>*/}
      {/*          But the correct value is: <span className="font-medium">{record.wantValue}</span>*/}
      {/*        </p>*/}
      {/*        <p>*/}
      {/*          Once you fix this, it will take at least{" "}*/}
      {/*          {record.actualTtlSeconds} seconds for the change to propagate,*/}
      {/*          because that's the time-to-live (TTL) you configured on the*/}
      {/*          incorrect record.*/}
      {/*        </p>*/}
      {/*      </div>*/}
      {/*    </TableCell>*/}
      {/*  </TableRow>*/}
      {/*)}*/}
    </>
  );
};

const schema = z.object({
  pendingDomain: z.string(),
});

const EditCustomAuthDomainButton = () => {
  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      pendingDomain: '',
    },
  });

  const { data: getVaultDomainSettingsResponse, refetch } = useQuery(
    getVaultDomainSettings,
  );
  useEffect(() => {
    if (getVaultDomainSettingsResponse?.vaultDomainSettings) {
      form.reset({
        pendingDomain:
          getVaultDomainSettingsResponse?.vaultDomainSettings?.pendingDomain,
      });
    }
  }, [getVaultDomainSettingsResponse]);

  const updateVaultDomainSettingsMutation = useMutation(
    updateVaultDomainSettings,
  );
  const [open, setOpen] = useState(false);
  const handleSubmit = async (values: z.infer<typeof schema>) => {
    await updateVaultDomainSettingsMutation.mutateAsync({
      vaultDomainSettings: {
        pendingDomain: values.pendingDomain,
      },
    });
    await refetch();
    setOpen(false);
  };

  return (
    <AlertDialog open={open} onOpenChange={setOpen}>
      <AlertDialogTrigger>
        <Button variant="outline">Edit</Button>
      </AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Edit Vault Domain</AlertDialogTitle>
          <AlertDialogDescription>
            Configure a custom domain for your Vault.
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
              name="pendingDomain"
              render={({ field }: { field: any }) => (
                <FormItem>
                  <FormLabel>Custom Vault Domain</FormLabel>
                  <FormControl>
                    <Input
                      className="max-w-96"
                      placeholder="vault.company.com"
                      {...field}
                    />
                  </FormControl>
                  <FormDescription>
                    A custom domain for your Vault. Typically, you'll use
                    "vault.company.com", where "company.com" is your company
                    domain.
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />
            <AlertDialogFooter className="mt-8">
              <AlertDialogCancel>Cancel</AlertDialogCancel>
              <Button
                type="submit"
                disabled={updateVaultDomainSettingsMutation.isPending}
              >
                {updateVaultDomainSettingsMutation.isPending && <Loader />}
                Save
              </Button>
            </AlertDialogFooter>
          </form>
        </Form>
      </AlertDialogContent>
    </AlertDialog>
  );
};
