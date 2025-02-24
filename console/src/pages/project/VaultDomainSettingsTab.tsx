import React from 'react';
import { getVaultDomainSettings } from '@/gen/tesseral/backend/v1/backend-BackendService_connectquery';
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
import { useQuery } from '@connectrpc/connect-query';
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

export const VaultDomainSettingsTab = () => {
  let { data: getVaultDomainSettingsResponse } = useQuery(
    getVaultDomainSettings,
  );
  getVaultDomainSettingsResponse = {
    vaultDomainSettings: {
      pendingDomain: 'vault1337.ucarion.com',
      currentDomain: 'project-4st5ccpz7bb29ho1hxeln03rx.laresset-dev1.app',
      dnsRecords: [
        {
          type: 'CNAME',
          name: 'vault1337.ucarion.com',
          wantValue: 'vault-cname.laresset-dns-dev1.com',
        },
        {
          type: 'TXT',
          name: '_tesseral_project_verification.vault1337.ucarion.com',
          wantValue: 'project_4st5ccpz7bb29ho1hxeln03rx',
        },
        {
          type: 'MX',
          name: 'mail.vault1337.ucarion.com',
          wantValue: '10 feedback-smtp.us-west-2.amazonses.com',
        },
        {
          type: 'TXT',
          name: 'mail.vault1337.ucarion.com',
          wantValue: 'v=spf1 include:amazonses.com ~all',
        },
        {
          type: 'CNAME',
          name: 'lmi5bww65bbdqt3zl3uppvaeqsm2hjit._domainkey.vault1337.ucarion.com',
          wantValue: 'lmi5bww65bbdqt3zl3uppvaeqsm2hjit.dkim.amazonses.com',
          actualValues: [
            'lmi5bww65bbdqt3zl3uppvaeqsm2hjit.dkim.amazonses.com.',
          ],
          actualTtlSeconds: 300,
        },
        {
          type: 'CNAME',
          name: 'kfvmppssbf3ttbjbfnqwoypn5ergkctl._domainkey.vault1337.ucarion.com',
          wantValue: 'kfvmppssbf3ttbjbfnqwoypn5ergkctl.dkim.amazonses.com',
        },
        {
          type: 'CNAME',
          name: 'f2ije72pbotsduhxniq2hkm3ujvvpetx._domainkey.vault1337.ucarion.com',
          wantValue: 'f2ije72pbotsduhxniq2hkm3ujvvpetx.dkim.amazonses.com',
        },
      ],
    },
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
                  {
                    getVaultDomainSettingsResponse?.vaultDomainSettings
                      ?.currentDomain
                  }
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
        <Card>
          <CardHeader>
            <CardTitle>DNS Records</CardTitle>
            <CardDescription>
              You need to add the following DNS records before you can use{' '}
              <span className="font-medium">
                {
                  getVaultDomainSettingsResponse?.vaultDomainSettings
                    ?.pendingDomain
                }
              </span>{' '}
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
                </TableRow>
              </TableHeader>
              <TableBody>
                {getVaultDomainSettingsResponse?.vaultDomainSettings?.dnsRecords?.map(
                  (record, i) => <DNSRecordRows key={i} record={record} />,
                )}
              </TableBody>
            </Table>
          </CardContent>
        </Card>
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
