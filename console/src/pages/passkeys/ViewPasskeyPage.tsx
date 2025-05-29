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
  ConsoleCardDescription,
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
import React, { useMemo, useState } from 'react';
import { useNavigate, useParams } from 'react-router';
import { useMutation, useQuery } from '@connectrpc/connect-query';
import {
  deletePasskey,
  getOrganization,
  getPasskey,
  getUser,
  updatePasskey,
} from '@/gen/tesseral/backend/v1/backend-BackendService_connectquery';
import { DateTime } from 'luxon';
import { timestampDate } from '@bufbuild/protobuf/wkt';
import { toast } from 'sonner';
import {
  AlertDialog,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/components/ui/alert-dialog';
import { Button } from '@/components/ui/button';
import { TabBar } from '@/components/ui/tab-bar';

export const ViewPasskeyPage = () => {
  const { organizationId, userId, passkeyId } = useParams();
  const { data: getOrganizationResponse } = useQuery(getOrganization, {
    id: organizationId,
  });
  const { data: getUserResponse } = useQuery(getUser, {
    id: userId,
  });
  const { data: getPasskeyResponse } = useQuery(getPasskey, {
    id: passkeyId,
  });

  const credentialId = useMemo(() => {
    if (!getPasskeyResponse?.passkey?.credentialId) {
      return;
    }

    return Array.from(getPasskeyResponse.passkey.credentialId)
      .map((byte) => byte.toString(16).padStart(2, '0'))
      .join('');
  }, [getPasskeyResponse?.passkey?.credentialId]);

  return (
    <>
      <PageHeader>
        <PageTitle>Passkey</PageTitle>
        <PageCodeSubtitle>{passkeyId}</PageCodeSubtitle>
        <PageDescription>
          A Passkey is a secondary authentication method tied a User, such a
          security key or Touch ID.
        </PageDescription>
      </PageHeader>

      <PageContent>
        <ConsoleCard className="my-8">
          <ConsoleCardHeader>
            <ConsoleCardTitle>General settings</ConsoleCardTitle>
            <ConsoleCardDescription>
              Basic settings for this passkey.
            </ConsoleCardDescription>
          </ConsoleCardHeader>
          <ConsoleCardContent>
            <DetailsGrid>
              <DetailsGridColumn>
                <DetailsGridEntry>
                  <DetailsGridKey>Status</DetailsGridKey>
                  <DetailsGridValue>
                    {getPasskeyResponse?.passkey?.disabled
                      ? 'Disabled'
                      : 'Active'}
                  </DetailsGridValue>
                </DetailsGridEntry>
                <DetailsGridEntry>
                  <DetailsGridKey>Vendor</DetailsGridKey>
                  <DetailsGridValue>
                    {getPasskeyResponse?.passkey?.aaguid &&
                    AAGUIDS[getPasskeyResponse.passkey.aaguid]
                      ? AAGUIDS[getPasskeyResponse.passkey.aaguid]
                      : 'Other'}
                  </DetailsGridValue>
                </DetailsGridEntry>
                <DetailsGridEntry>
                  <DetailsGridKey>AAGUID</DetailsGridKey>
                  <DetailsGridValue>
                    {getPasskeyResponse?.passkey?.aaguid}
                  </DetailsGridValue>
                </DetailsGridEntry>
              </DetailsGridColumn>
              <DetailsGridColumn>
                <DetailsGridEntry>
                  <DetailsGridKey>Public Key</DetailsGridKey>
                  <DetailsGridValue>
                    {getPasskeyResponse?.passkey?.publicKeyPkix && (
                      <a
                        className="font-medium underline underline-offset-2 decoration-muted-foreground/40"
                        download={`Public Key ${passkeyId}.pem`}
                        href={`data:text/plain;base64,${btoa(getPasskeyResponse.passkey.publicKeyPkix)}`}
                      >
                        Download (.pem)
                      </a>
                    )}
                  </DetailsGridValue>
                </DetailsGridEntry>
                <DetailsGridEntry>
                  <DetailsGridKey>Credential ID</DetailsGridKey>
                  <DetailsGridValue>{credentialId}</DetailsGridValue>
                </DetailsGridEntry>
              </DetailsGridColumn>
              <DetailsGridColumn>
                <DetailsGridEntry>
                  <DetailsGridKey>Created</DetailsGridKey>
                  <DetailsGridValue>
                    {getPasskeyResponse?.passkey?.createTime &&
                      DateTime.fromJSDate(
                        timestampDate(getPasskeyResponse.passkey.createTime),
                      ).toRelative()}
                  </DetailsGridValue>
                </DetailsGridEntry>
                <DetailsGridEntry>
                  <DetailsGridKey>Updated</DetailsGridKey>
                  <DetailsGridValue>
                    {getPasskeyResponse?.passkey?.updateTime &&
                      DateTime.fromJSDate(
                        timestampDate(getPasskeyResponse.passkey.updateTime),
                      ).toRelative()}
                  </DetailsGridValue>
                </DetailsGridEntry>
              </DetailsGridColumn>
            </DetailsGrid>
          </ConsoleCardContent>
        </ConsoleCard>

        <DangerZoneCard />
      </PageContent>
    </>
  );
};

const DangerZoneCard = () => {
  const { organizationId, userId, passkeyId } = useParams();
  const [confirmDeleteOpen, setConfirmDeleteOpen] = useState(false);
  const { data: getPasskeyResponse, refetch } = useQuery(getPasskey, {
    id: passkeyId,
  });
  const { data: getUserResponse } = useQuery(getUser, {
    id: userId,
  });

  const updatePasskeyMutation = useMutation(updatePasskey);
  const handleDisable = async () => {
    await updatePasskeyMutation.mutateAsync({
      id: passkeyId,
      passkey: {
        disabled: true,
      },
    });

    await refetch();
    toast.success('Passkey disabled');
  };

  const handleEnable = async () => {
    await updatePasskeyMutation.mutateAsync({
      id: passkeyId,
      passkey: {
        disabled: false,
      },
    });

    await refetch();
    toast.success('Passkey enabled');
  };

  const handleDelete = () => {
    setConfirmDeleteOpen(true);
  };

  const deletePasskeyMutation = useMutation(deletePasskey);
  const navigate = useNavigate();
  const handleConfirmDelete = async () => {
    await deletePasskeyMutation.mutateAsync({
      id: passkeyId,
    });

    toast.success('Passkey deleted');
    navigate(`/organizations/${organizationId}/users/${userId}`);
  };

  return (
    <>
      <AlertDialog open={confirmDeleteOpen} onOpenChange={setConfirmDeleteOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Delete Passkey?</AlertDialogTitle>
            <AlertDialogDescription>
              Deleting a passkey cannot be undone.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <Button variant="destructive" onClick={handleConfirmDelete}>
              Permanently Delete Passkey
            </Button>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      <ConsoleCard className="border-destructive">
        <ConsoleCardHeader>
          <ConsoleCardTitle>Danger Zone</ConsoleCardTitle>
        </ConsoleCardHeader>

        <ConsoleCardContent className="space-y-4">
          {getPasskeyResponse?.passkey?.disabled ? (
            <div className="flex justify-between items-center">
              <div>
                <div className="text-sm font-semibold">Enable Passkey</div>
                <p className="text-sm">
                  Enable this passkey.{' '}
                  <span className="font-medium">
                    {getUserResponse?.user?.email}
                  </span>{' '}
                  will be required to authenticate with this passkey (or another
                  active passkey) when logging in.
                </p>
              </div>

              <Button variant="destructive" onClick={handleEnable}>
                Enable Passkey
              </Button>
            </div>
          ) : (
            <div className="flex justify-between items-center">
              <div>
                <div className="text-sm font-semibold">Enable Passkey</div>
                <p className="text-sm">
                  Disable this passkey.{' '}
                  <span className="font-medium">
                    {getUserResponse?.user?.email}
                  </span>{' '}
                  will not be required to authenticate with this passkey when
                  logging in.
                </p>
              </div>

              <Button variant="destructive" onClick={handleDisable}>
                Disable Passkey
              </Button>
            </div>
          )}
          <div className="flex justify-between items-center">
            <div>
              <div className="text-sm font-semibold">Delete Passkey</div>
              <p className="text-sm">
                Delete this passkey. This cannot be undone.
              </p>
            </div>

            <Button variant="destructive" onClick={handleDelete}>
              Delete Passkey
            </Button>
          </div>
        </ConsoleCardContent>
      </ConsoleCard>
    </>
  );
};

// curl https://raw.githubusercontent.com/passkeydeveloper/passkey-authenticator-aaguids/refs/heads/main/aaguid.json | jq 'map_values(.name)'
const AAGUIDS: Record<string, string> = {
  'ea9b8d66-4d01-1d21-3ce4-b6b48cb575d4': 'Google Password Manager',
  'adce0002-35bc-c60a-648b-0b25f1f05503': 'Chrome on Mac',
  '08987058-cadc-4b81-b6e1-30de50dcbe96': 'Windows Hello',
  '9ddd1817-af5a-4672-a2b9-3e3dd95000a9': 'Windows Hello',
  '6028b017-b1d4-4c02-b4b3-afcdafc96bb2': 'Windows Hello',
  'dd4ec289-e01d-41c9-bb89-70fa845d4bf2': 'iCloud Keychain (Managed)',
  '531126d6-e717-415c-9320-3d9aa6981239': 'Dashlane',
  'bada5566-a7aa-401f-bd96-45619a55120d': '1Password',
  'b84e4048-15dc-4dd0-8640-f4f60813c8af': 'NordPass',
  '0ea242b4-43c4-4a1b-8b17-dd6d0b6baec6': 'Keeper',
  '891494da-2c90-4d31-a9cd-4eab0aed1309': 'Sésame',
  'f3809540-7f14-49c1-a8b3-8f813b225541': 'Enpass',
  'b5397666-4885-aa6b-cebf-e52262a439a2': 'Chromium Browser',
  '771b48fd-d3d4-4f74-9232-fc157ab0507a': 'Edge on Mac',
  '39a5647e-1853-446c-a1f6-a79bae9f5bc7': 'IDmelon',
  'd548826e-79b4-db40-a3d8-11116f7e8349': 'Bitwarden',
  'fbfc3007-154e-4ecc-8c0b-6e020557d7bd': 'iCloud Keychain',
  '53414d53-554e-4700-0000-000000000000': 'Samsung Pass',
  '66a0ccb3-bd6a-191f-ee06-e375c50b9846': 'Thales Bio iOS SDK',
  '8836336a-f590-0921-301d-46427531eee6': 'Thales Bio Android SDK',
  'cd69adb5-3c7a-deb9-3177-6800ea6cb72a': 'Thales PIN Android SDK',
  '17290f1e-c212-34d0-1423-365d729f09d9': 'Thales PIN iOS SDK',
  '50726f74-6f6e-5061-7373-50726f746f6e': 'Proton Pass',
  'fdb141b2-5d84-443e-8a35-4698c205a502': 'KeePassXC',
  'cc45f64e-52a2-451b-831a-4edd8022a202': 'ToothPic Passkey Provider',
  'bfc748bb-3429-4faa-b9f9-7cfa9f3b76d0': 'iPasswords',
  'b35a26b2-8f6e-4697-ab1d-d44db4da28c6': 'Zoho Vault',
  'b78a0a55-6ef8-d246-a042-ba0f6d55050c': 'LastPass',
  'de503f9c-21a4-4f76-b4b7-558eb55c6f89': 'Devolutions',
  '22248c4c-7a12-46e2-9a41-44291b373a4d': 'LogMeOnce',
};
