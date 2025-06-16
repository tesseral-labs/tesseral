import { useMutation, useQuery } from "@connectrpc/connect-query";
import React from "react";
import { toast } from "sonner";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardAction,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Table,
  TableBody,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import {
  enableEmailSendFromDomain,
  getProject,
  getVaultDomainSettings,
} from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";

import { DNSRecordRows } from "./DNSRecordRows";

export function VaultEmailSendFromDomainRecordsCard() {
  const { data: getProjectResponse } = useQuery(getProject, {});
  const { data: getVaultDomainSettingsResponse } = useQuery(
    getVaultDomainSettings,
  );

  const customEmailSendFromDomainActive =
    getProjectResponse?.project?.emailSendFromDomain ===
    `mail.${getVaultDomainSettingsResponse?.vaultDomainSettings?.pendingDomain}`;

  return (
    <Card className="mt-8">
      <CardHeader>
        <CardTitle className="flex items-center">
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
        </CardTitle>
        <CardDescription>
          {customEmailSendFromDomainActive ? (
            <p>
              You need to keep these DNS records in place so that Tesseral can
              continue to send emails from{" "}
              <span className="font-medium">{`noreply@mail.${getVaultDomainSettingsResponse?.vaultDomainSettings?.pendingDomain}`}</span>
              .
            </p>
          ) : (
            <p>
              You can optionally configure the emails Tesseral sends to your end
              users to come from{" "}
              <span className="font-medium">{`noreply@mail.${getVaultDomainSettingsResponse?.vaultDomainSettings?.pendingDomain}`}</span>
              , instead of the Tesseral-provided{" "}
              <span className="font-medium">noreply@mail.tesseral.app</span>.
            </p>
          )}
        </CardDescription>

        {!customEmailSendFromDomainActive && (
          <CardAction>
            <EnableEmailSendFromDomainButton />
          </CardAction>
        )}
      </CardHeader>
      <CardContent>
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
      </CardContent>
    </Card>
  );
}

function EnableEmailSendFromDomainButton() {
  const { refetch } = useQuery(getProject, {});
  const { data: getVaultDomainSettingsResponse } = useQuery(
    getVaultDomainSettings,
  );

  const enableEmailSendFromDomainMutation = useMutation(
    enableEmailSendFromDomain,
  );

  async function handleSubmit() {
    await enableEmailSendFromDomainMutation.mutateAsync({});
    await refetch();
    toast.success("Custom Email Send-From Domain enabled");
  }

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
}
