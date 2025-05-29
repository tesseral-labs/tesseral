import React, { useEffect, useState } from 'react';
import {
  createStripeCheckoutLink,
  enableCustomVaultDomain,
  enableEmailSendFromDomain,
  getProject,
  getProjectEntitlements,
  getVaultDomainSettings,
  updateVaultDomainSettings,
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
import { CircleXIcon } from 'lucide-react';
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
import { StatusIndicator } from '@/components/status-indicator';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { toast } from 'sonner';

export const VaultDomainSettingsTab = () => {
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

  const { data: getProjectResponse } = useQuery(getProject, {});
  const { data: getVaultDomainSettingsResponse } = useQuery(
    getVaultDomainSettings,
  );

  const customVaultDomainActive =
    getProjectResponse?.project?.vaultDomain ===
    getVaultDomainSettingsResponse?.vaultDomainSettings?.pendingDomain;
  const customEmailSendFromDomainActive =
    getProjectResponse?.project?.emailSendFromDomain ===
    `mail.${getVaultDomainSettingsResponse?.vaultDomainSettings?.pendingDomain}`;

  return (
    <div className="space-y-8">
      <ConsoleCard>
        <ConsoleCardHeader>
          <ConsoleCardDetails>
            <ConsoleCardTitle>Vault Domain Settings</ConsoleCardTitle>
            <ConsoleCardDescription>
              Configure a custom domain for your Vault.
            </ConsoleCardDescription>
          </ConsoleCardDetails>
          <EditCustomAuthDomainButton />
        </ConsoleCardHeader>
        <ConsoleCardContent>
          {/* do not treat undefined as unentitled, to avoid flickering here */}
          {getProjectEntitlementsResponse?.entitledBackendApiKeys === false ? (
            <div className="text-sm my-8 w-full flex flex-col items-center justify-center space-y-6">
              <div className="font-medium">
                Custom Vault Domains are available on the Growth Tier.
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
            <DetailsGrid>
              <DetailsGridColumn>
                <DetailsGridEntry>
                  <DetailsGridKey>Current Vault Domain</DetailsGridKey>
                  <DetailsGridValue>
                    {getProjectResponse?.project?.vaultDomain}
                  </DetailsGridValue>
                </DetailsGridEntry>
                <DetailsGridEntry>
                  <DetailsGridKey>
                    Current Email Send-From Domain
                  </DetailsGridKey>
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
          )}
        </ConsoleCardContent>
      </ConsoleCard>

      {getVaultDomainSettingsResponse?.vaultDomainSettings?.pendingDomain && (
        <>
          <ConsoleCard>
            <ConsoleCardHeader>
              <ConsoleCardDetails>
                <ConsoleCardTitle className="flex items-center">
                  <span>Vault Domain Records</span>
                  {customVaultDomainActive && (
                    <Badge
                      variant="outline"
                      className="ml-4 bg-green-50 text-green-700 border-green-200"
                    >
                      <span className="mr-1 h-2 w-2 rounded-full bg-green-500 inline-block" />
                      Live
                    </Badge>
                  )}
                </ConsoleCardTitle>
                <ConsoleCardDescription>
                  {customVaultDomainActive ? (
                    <p>
                      You need to keep these DNS records in place so that{' '}
                      <span className="font-medium">
                        {getProjectResponse?.project?.vaultDomain}
                      </span>{' '}
                      continues to work as your Vault domain.
                    </p>
                  ) : (
                    <p>
                      You need to add the following DNS records before you can
                      use{' '}
                      <span className="font-medium">
                        {
                          getVaultDomainSettingsResponse?.vaultDomainSettings
                            ?.pendingDomain
                        }
                      </span>{' '}
                      as your Vault domain.
                    </p>
                  )}
                </ConsoleCardDescription>
              </ConsoleCardDetails>

              {!customVaultDomainActive && <EnableCustomVaultDomainButton />}
            </ConsoleCardHeader>
            <ConsoleCardTableContent>
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Status</TableHead>
                    <TableHead>Type</TableHead>
                    <TableHead>Name</TableHead>
                    <TableHead>Value</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {getVaultDomainSettingsResponse?.vaultDomainSettings?.vaultDomainRecords?.map(
                    (record, i) => <DNSRecordRows key={i} record={record} />,
                  )}
                </TableBody>
              </Table>
            </ConsoleCardTableContent>
          </ConsoleCard>
          <ConsoleCard className="mt-8">
            <ConsoleCardHeader>
              <ConsoleCardDetails>
                <ConsoleCardTitle className="flex items-center">
                  <span>Email Send-From Records</span>

                  {customEmailSendFromDomainActive ? (
                    <Badge
                      variant="outline"
                      className="ml-4 bg-green-50 text-green-700 border-green-200"
                    >
                      <span className="mr-1 h-2 w-2 rounded-full bg-green-500 inline-block" />
                      Live
                    </Badge>
                  ) : (
                    <Badge className="ml-4" variant="outline">
                      Optional
                    </Badge>
                  )}
                </ConsoleCardTitle>
                <ConsoleCardDescription>
                  {customEmailSendFromDomainActive ? (
                    <p>
                      You need to keep these DNS records in place so that
                      Tesseral can continue to send emails from{' '}
                      <span className="font-medium">{`noreply@mail.${getVaultDomainSettingsResponse?.vaultDomainSettings?.pendingDomain}`}</span>
                      .
                    </p>
                  ) : (
                    <p>
                      You can optionally configure the emails Tesseral sends to
                      your end users to come from{' '}
                      <span className="font-medium">{`noreply@mail.${getVaultDomainSettingsResponse?.vaultDomainSettings?.pendingDomain}`}</span>
                      , instead of the Tesseral-provided{' '}
                      <span className="font-medium">
                        noreply@mail.tesseral.app
                      </span>
                      .
                    </p>
                  )}
                </ConsoleCardDescription>
              </ConsoleCardDetails>

              {!customEmailSendFromDomainActive && (
                <EnableEmailSendFromDomainButton />
              )}
            </ConsoleCardHeader>
            <ConsoleCardTableContent>
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Status</TableHead>
                    <TableHead>Type</TableHead>
                    <TableHead>Name</TableHead>
                    <TableHead>Value</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {getVaultDomainSettingsResponse?.vaultDomainSettings?.emailSendFromRecords?.map(
                    (record, i) => <DNSRecordRows key={i} record={record} />,
                  )}
                </TableBody>
              </Table>
            </ConsoleCardTableContent>
          </ConsoleCard>
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
        <TableCell>
          {record.correct && (
            <StatusIndicator variant="success">Configured</StatusIndicator>
          )}
          {noValue && (
            <StatusIndicator variant="pending">No record</StatusIndicator>
          )}
          {(tooManyValues || incorrectValue) && (
            <StatusIndicator variant="error">Misconfigured</StatusIndicator>
          )}
        </TableCell>
        <TableCell>{record.type}</TableCell>
        <TableCell>{record.name}</TableCell>
        <TableCell>{record.wantValue}</TableCell>
      </TableRow>

      {incorrectValue && (
        <TableRow className="bg-red-50/50 hover:bg-red-50/50">
          <TableCell colSpan={4}>
            <Alert variant="destructive" className="bg-white">
              <CircleXIcon className="w-5 h-5 text-red-500" />
              <AlertTitle>
                <span className="font-mono">{record.name}</span> is
                misconfigured
              </AlertTitle>
              <AlertDescription>
                <p className="mt-2">This record has the wrong value.</p>

                <Table>
                  <TableBody>
                    <TableRow className="border-destructive/25 hover:bg-white">
                      <TableCell>Expected</TableCell>
                      <TableCell className="font-mono">
                        {record.wantValue}
                      </TableCell>
                    </TableRow>
                    <TableRow className="border-destructive/25 hover:bg-white">
                      <TableCell>Actual</TableCell>
                      <TableCell className="font-mono">
                        {record.actualValues[0]}
                      </TableCell>
                    </TableRow>
                  </TableBody>
                </Table>
                <p className="mt-2">
                  It will take at least {record.actualTtlSeconds} seconds for
                  any change you make here to propagate, because that's the
                  time-to-live (TTL) you configured on this incorrect record.
                </p>
              </AlertDescription>
            </Alert>
          </TableCell>
        </TableRow>
      )}

      {tooManyValues && (
        <TableRow className="bg-red-50/50 hover:bg-red-50/50">
          <TableCell colSpan={4}>
            <Alert variant="destructive" className="bg-white">
              <CircleXIcon className="w-5 h-5 text-red-500" />
              <AlertTitle>
                <span className="font-mono">{record.name}</span> is
                misconfigured
              </AlertTitle>
              <AlertDescription>
                <p className="mt-2">
                  This record has too many values. Delete the following records:
                </p>

                <Table>
                  <TableHeader className="border-destructive/25 border-b">
                    <TableRow className="hover:bg-white">
                      <TableHead className="text-destructive">Type</TableHead>
                      <TableHead className="text-destructive">Name</TableHead>
                      <TableHead className="text-destructive">Value</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {record.actualValues
                      ?.filter((v) => v !== record.wantValue)
                      .map((v, i) => (
                        <TableRow
                          key={i}
                          className="border-destructive/25 hover:bg-white"
                        >
                          <TableCell>{record.type}</TableCell>
                          <TableCell>{record.name}</TableCell>
                          <TableCell>{v}</TableCell>
                        </TableRow>
                      ))}
                  </TableBody>
                </Table>
                <p className="mt-2">
                  It will take at least {record.actualTtlSeconds} seconds for
                  any change you make here to propagate, because that's the
                  time-to-live (TTL) you configured on these records.
                </p>
              </AlertDescription>
            </Alert>
          </TableCell>
        </TableRow>
      )}
    </>
  );
};

const EnableCustomVaultDomainButton = () => {
  const { refetch } = useQuery(getProject, {});
  const { data: getVaultDomainSettingsResponse } = useQuery(
    getVaultDomainSettings,
  );

  const enableCustomVaultDomainMutation = useMutation(enableCustomVaultDomain);
  const handleSubmit = async () => {
    await enableCustomVaultDomainMutation.mutateAsync({});
    await refetch();
    toast.success('Custom Vault Domain enabled');
  };

  return (
    <>
      {getVaultDomainSettingsResponse?.vaultDomainSettings
        ?.pendingVaultDomainReady ? (
        <Button variant="outline" onClick={handleSubmit}>
          Enable Custom Vault Domain
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
                Enable Custom Vault Domain
              </Button>
            </TooltipTrigger>
            <TooltipContent className="max-w-96">
              Your DNS records need to be correct and widely propagated before
              you can enable this.
            </TooltipContent>
          </Tooltip>
        </TooltipProvider>
      )}
    </>
  );
};

const EnableEmailSendFromDomainButton = () => {
  const { refetch } = useQuery(getProject, {});
  const { data: getVaultDomainSettingsResponse } = useQuery(
    getVaultDomainSettings,
  );

  const enableEmailSendFromDomainMutation = useMutation(
    enableEmailSendFromDomain,
  );
  const handleSubmit = async () => {
    await enableEmailSendFromDomainMutation.mutateAsync({});
    await refetch();
    toast.success('Custom Email Send-From Domain enabled');
  };

  return (
    <>
      {getVaultDomainSettingsResponse?.vaultDomainSettings
        ?.pendingSendFromDomainReady ? (
        <Button variant="outline" onClick={handleSubmit}>
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
              Your DNS records need to be correct and widely propagated before
              you can enable this.
            </TooltipContent>
          </Tooltip>
        </TooltipProvider>
      )}
    </>
  );
};

const schema = z.object({
  pendingDomain: z.string(),
});

const EditCustomAuthDomainButton = () => {
  const { data: getProjectEntitlementsResponse } = useQuery(
    getProjectEntitlements,
    {},
  );

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
      <AlertDialogTrigger
        disabled={!getProjectEntitlementsResponse?.entitledCustomVaultDomains}
      >
        <Button
          variant="outline"
          disabled={!getProjectEntitlementsResponse?.entitledCustomVaultDomains}
        >
          Edit
        </Button>
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
