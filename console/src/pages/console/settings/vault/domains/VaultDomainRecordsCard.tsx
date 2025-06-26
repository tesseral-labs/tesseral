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
  enableCustomVaultDomain,
  getProject,
  getVaultDomainSettings,
} from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";

import { DNSRecordRows } from "./DNSRecordRows";

export function VaultDomainRecordsCard() {
  const { data: getProjectResponse } = useQuery(getProject, {});
  const { data: getVaultDomainSettingsResponse } = useQuery(
    getVaultDomainSettings,
  );

  const customVaultDomainActive =
    getProjectResponse?.project?.vaultDomain ===
    getVaultDomainSettingsResponse?.vaultDomainSettings?.pendingDomain;

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center">
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
        </CardTitle>
        <CardDescription>
          {customVaultDomainActive ? (
            <p>
              You need to keep these DNS records in place so that{" "}
              <span className="font-medium">
                {getProjectResponse?.project?.vaultDomain}
              </span>{" "}
              continues to work as your Vault domain.
            </p>
          ) : (
            <p>
              You need to add the following DNS records before you can use{" "}
              <span className="font-medium">
                {
                  getVaultDomainSettingsResponse?.vaultDomainSettings
                    ?.pendingDomain
                }
              </span>{" "}
              as your Vault domain.
            </p>
          )}
        </CardDescription>

        {!customVaultDomainActive && (
          <CardAction>
            <EnableCustomVaultDomainButton />
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
            {getVaultDomainSettingsResponse?.vaultDomainSettings?.vaultDomainRecords?.map(
              (record, i) => <DNSRecordRows key={i} record={record} />,
            )}
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  );
}

function EnableCustomVaultDomainButton() {
  const { refetch } = useQuery(getProject, {});
  const { data: getVaultDomainSettingsResponse } = useQuery(
    getVaultDomainSettings,
  );
  const enableCustomVaultDomainMutation = useMutation(enableCustomVaultDomain);

  async function handleSubmit() {
    await enableCustomVaultDomainMutation.mutateAsync({});
    await refetch();
    toast.success("Custom Vault Domain enabled");
  }

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
}
