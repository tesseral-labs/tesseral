import React from "react";
import {
  getVaultDomainSettings
} from '@/gen/tesseral/backend/v1/backend-BackendService_connectquery';
import {
  Card, CardContent,
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
import { useQuery } from '@connectrpc/connect-query';
import {
  DetailsGrid,
  DetailsGridColumn,
  DetailsGridEntry, DetailsGridKey, DetailsGridValue,
} from '@/components/details-grid';
import { Badge } from '@/components/ui/badge';
import { CheckIcon, XIcon } from 'lucide-react';

export const VaultDomainSettingsTab = () => {
  let { data: getVaultDomainSettingsResponse } = useQuery(getVaultDomainSettings);
  getVaultDomainSettingsResponse = {
    "vaultDomainSettings": {
      "currentDomain": "console.tesseral.example.com",
      "pendingDomain": "vault1337.ucarion.com",
      "mainRecord": {
        "type": "CNAME",
        "name": "vault1337.ucarion.com"
      },
      "projectVerificationRecord": {
        "type": "TXT",
        "name": "_tesseral_project_verification.vault1337.ucarion.com",
        "value": "project_6ja4u7a7bslj2j9lj9bd6qawv"
      },
      "dkimRecords": [
        {
          "type": "CNAME",
          "name": "lmi5bww65bbdqt3zl3uppvaeqsm2hjit._domainkey.vault1337.ucarion.com",
          "value": "lmi5bww65bbdqt3zl3uppvaeqsm2hjit.dkim.amazonses.com"
        },
        {
          "type": "CNAME",
          "name": "kfvmppssbf3ttbjbfnqwoypn5ergkctl._domainkey.vault1337.ucarion.com",
          "value": "kfvmppssbf3ttbjbfnqwoypn5ergkctl.dkim.amazonses.com"
        },
        {
          "type": "CNAME",
          "name": "f2ije72pbotsduhxniq2hkm3ujvvpetx._domainkey.vault1337.ucarion.com",
          "value": "f2ije72pbotsduhxniq2hkm3ujvvpetx.dkim.amazonses.com"
        }
      ],
      "spfRecords": [
        {
          "type": "MX",
          "name": "mail.vault1337.ucarion.com"
        },
        {
          "type": "TXT",
          "name": "mail.vault1337.ucarion.com",
          "value": "v=spf1 include:amazonses.com ~all"
        }
      ]
    }
  } as any;

  return (
    <div className="space-y-8">
      <Card>
        <CardHeader>
          <CardTitle>Vault Domain Settings</CardTitle>
          <CardDescription>
            Configure a custom domain for your Vault.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <DetailsGrid>
            <DetailsGridColumn>
              <DetailsGridEntry>
                <DetailsGridKey>Current Domain</DetailsGridKey>
                <DetailsGridValue>
                  {getVaultDomainSettingsResponse?.vaultDomainSettings?.currentDomain}
                </DetailsGridValue>
              </DetailsGridEntry>
            </DetailsGridColumn>
            <DetailsGridColumn>
              <DetailsGridEntry>
                <DetailsGridKey>Pending Custom Domain</DetailsGridKey>
                <DetailsGridValue>
                  {getVaultDomainSettingsResponse?.vaultDomainSettings?.pendingDomain || '-'}
                </DetailsGridValue>
              </DetailsGridEntry>
            </DetailsGridColumn>
          </DetailsGrid>
        </CardContent>
      </Card>

      {getVaultDomainSettingsResponse?.vaultDomainSettings?.pendingDomain && (
        <Card>
          <CardHeader>
            <CardTitle>DNS Records</CardTitle>
            <CardDescription>
              You need to add the following DNS records before you can use <span
              className="font-medium">{getVaultDomainSettingsResponse?.vaultDomainSettings?.pendingDomain}</span>{" "}
              as your custom Vault domain.
            </CardDescription>
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
                <TableRow>
                  <TableCell>{getVaultDomainSettingsResponse?.vaultDomainSettings?.mainRecord?.type}</TableCell>
                  <TableCell>{getVaultDomainSettingsResponse?.vaultDomainSettings?.mainRecord?.name}</TableCell>
                  <TableCell className="font-mono">{getVaultDomainSettingsResponse?.vaultDomainSettings?.mainRecord?.value}</TableCell>
                  <TableCell>
                    <RecordStatus configured={!!getVaultDomainSettingsResponse?.vaultDomainSettings?.mainRecordConfigured} />
                  </TableCell>
                </TableRow>
                {getVaultDomainSettingsResponse?.vaultDomainSettings?.dkimRecords.map((record, index) => (
                  <TableRow key={`dkim-${index}`}>
                    <TableCell>{record.type}</TableCell>
                    <TableCell>{record.name}</TableCell>
                    <TableCell className="font-mono">{record.value}</TableCell>
                    <TableCell>
                      <RecordStatus configured={!!getVaultDomainSettingsResponse?.vaultDomainSettings?.dkimConfigured} />
                    </TableCell>
                  </TableRow>
                ))}
                {getVaultDomainSettingsResponse?.vaultDomainSettings?.spfRecords.map((record, index) => (
                  <TableRow key={`spf-${index}`}>
                    <TableCell>{record.type}</TableCell>
                    <TableCell>{record.name}</TableCell>
                    <TableCell className="font-mono">{record.value}</TableCell>
                    <TableCell>
                      <RecordStatus configured={!!getVaultDomainSettingsResponse?.vaultDomainSettings?.spfConfigured} />
                    </TableCell>
                  </TableRow>
                ))}
                <TableRow>
                  <TableCell>{getVaultDomainSettingsResponse?.vaultDomainSettings?.projectVerificationRecord?.type}</TableCell>
                  <TableCell>{getVaultDomainSettingsResponse?.vaultDomainSettings?.projectVerificationRecord?.name}</TableCell>
                  <TableCell className="font-mono">{getVaultDomainSettingsResponse?.vaultDomainSettings?.projectVerificationRecord?.value}</TableCell>
                  <TableCell>
                    <RecordStatus configured={!!getVaultDomainSettingsResponse?.vaultDomainSettings?.projectVerificationRecordConfigured} />
                  </TableCell>
                </TableRow>
              </TableBody>
            </Table>
          </CardContent>
        </Card>
      )}
    </div>
  )
};

const RecordStatus = ({ configured }: { configured: boolean }) => {
  return (
    <Badge variant={configured ? "default" : "destructive"} className="flex items-center gap-1">
      {configured ? (
        <>
          <CheckIcon size={14} />
          Configured
        </>
      ) : (
        <>
          <XIcon size={14} />
          Not Configured
        </>
      )}
    </Badge>
  )
}
