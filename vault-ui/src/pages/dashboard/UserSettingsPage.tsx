import { timestampDate } from "@bufbuild/protobuf/wkt";
import { useMutation, useQuery } from "@connectrpc/connect-query";
import { DateTime } from "luxon";
import React, { useState } from "react";
import { toast } from "sonner";

import {
  AlertDialog,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  deleteMyPasskey,
  getPasskeyOptions,
  listMyPasskeys,
  registerPasskey,
  whoami,
} from "@/gen/tesseral/frontend/v1/frontend-FrontendService_connectquery";
import { base64urlEncode } from "@/lib/utils";
import { RegisterAuthenticatorAppButton } from "@/pages/dashboard/RegisterAuthenticatorAppButton";

export function UserSettingsPage() {
  const { data: whoamiResponse } = useQuery(whoami);
  const { data: listMyPasskeysResponse } = useQuery(listMyPasskeys);

  return (
    <div className="space-y-4">
      <Card>
        <CardHeader>
          <CardTitle>Profile Information</CardTitle>
          <CardDescription>
            Basic information about your account.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="text-sm font-medium">Email</div>
          <div className="text-sm">{whoamiResponse?.user?.email}</div>
        </CardContent>
      </Card>
      <Card>
        <CardHeader>
          <CardTitle>Multi-Factor Authentication</CardTitle>
          <CardDescription>
            Additional layers of security for your account.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            <div className="flex justify-between items-center">
              <div>
                <div className="text-sm font-medium">Authenticator App</div>
                <div className="text-sm">
                  {whoamiResponse?.user?.hasAuthenticatorApp
                    ? "Configured"
                    : "Not Configured"}
                </div>
              </div>

              <RegisterAuthenticatorAppButton />
            </div>

            <div className="flex justify-between items-center">
              <div>
                <div className="text-sm font-medium">Passkeys</div>
                <div className="text-sm">
                  {(listMyPasskeysResponse?.passkeys?.length ?? 0) > 0
                    ? "Configured"
                    : "Not Configured"}
                </div>
              </div>

              <RegisterPasskeyButton />
            </div>

            {listMyPasskeysResponse?.passkeys?.map((passkey) => (
              <div
                key={passkey.id}
                className="flex justify-between items-center"
              >
                <div>
                  <div className="text-sm font-medium">
                    Passkey {passkey.id}
                  </div>
                  <div className="text-sm flex gap-x-2 text-muted-foreground">
                    <span>
                      Created{" "}
                      {DateTime.fromJSDate(
                        timestampDate(passkey.createTime!),
                      ).toRelative()}
                    </span>
                    <span>&middot;</span>
                    <span>{AAGUIDS[passkey.aaguid] ?? "Unknown Vendor"}</span>
                  </div>
                </div>

                <DeletePasskeyButton passkeyId={passkey.id} />
              </div>
            ))}
          </div>
        </CardContent>
      </Card>
    </div>
  );
}

function RegisterPasskeyButton() {
  const { mutateAsync: getPasskeyOptionsAsync } =
    useMutation(getPasskeyOptions);
  const { mutateAsync: registerPasskeyAsync } = useMutation(registerPasskey);
  const { refetch: refetchListMyPasskeys } = useQuery(listMyPasskeys);

  async function handleRegisterPasskey() {
    const passkeyOptions = await getPasskeyOptionsAsync({});
    const credentialOptions: PublicKeyCredentialCreationOptions = {
      challenge: new Uint8Array([0]).buffer,
      rp: {
        id: passkeyOptions.rpId,
        name: passkeyOptions.rpName,
      },
      user: {
        id: new TextEncoder().encode(passkeyOptions.userId).buffer,
        name: passkeyOptions.userDisplayName,
        displayName: passkeyOptions.userDisplayName,
      },
      pubKeyCredParams: [
        { type: "public-key", alg: -7 }, // ECDSA with SHA-256
        { type: "public-key", alg: -257 }, // RSA with SHA-256
      ],
      timeout: 60000,
      attestation: "direct",
    };

    const credential = (await navigator.credentials.create({
      publicKey: credentialOptions,
    })) as PublicKeyCredential;

    if (!credential) {
      throw new Error("No credential returned");
    }

    await registerPasskeyAsync({
      rpId: passkeyOptions.rpId,
      attestationObject: base64urlEncode(
        (credential.response as AuthenticatorAttestationResponse)
          .attestationObject,
      ),
    });

    await refetchListMyPasskeys();

    toast.success("Passkey registered");
  }

  return (
    <Button variant="outline" onClick={handleRegisterPasskey}>
      Register Passkey
    </Button>
  );
}

function DeletePasskeyButton({ passkeyId }: { passkeyId: string }) {
  const { mutateAsync: deletePasskeyAsync } = useMutation(deleteMyPasskey);
  const { refetch: refetchListMyPasskeys } = useQuery(listMyPasskeys);

  const [dialogOpen, setDialogOpen] = useState(false);

  async function handleDeletePasskey() {
    await deletePasskeyAsync({ id: passkeyId });
    await refetchListMyPasskeys();
    toast.success("Passkey deleted");
  }

  return (
    <AlertDialog open={dialogOpen} onOpenChange={setDialogOpen}>
      <AlertDialogTrigger asChild>
        <Button variant="outline">Delete</Button>
      </AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Are you sure?</AlertDialogTitle>
          <AlertDialogDescription>
            This action cannot be undone. The passkey will be permanently
            deleted.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel>Cancel</AlertDialogCancel>
          <Button variant="destructive" onClick={handleDeletePasskey}>
            Delete Passkey
          </Button>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
}

// curl https://raw.githubusercontent.com/passkeydeveloper/passkey-authenticator-aaguids/refs/heads/main/aaguid.json | jq 'map_values(.name)'
const AAGUIDS: Record<string, string> = {
  "ea9b8d66-4d01-1d21-3ce4-b6b48cb575d4": "Google Password Manager",
  "adce0002-35bc-c60a-648b-0b25f1f05503": "Chrome on Mac",
  "08987058-cadc-4b81-b6e1-30de50dcbe96": "Windows Hello",
  "9ddd1817-af5a-4672-a2b9-3e3dd95000a9": "Windows Hello",
  "6028b017-b1d4-4c02-b4b3-afcdafc96bb2": "Windows Hello",
  "dd4ec289-e01d-41c9-bb89-70fa845d4bf2": "iCloud Keychain (Managed)",
  "531126d6-e717-415c-9320-3d9aa6981239": "Dashlane",
  "bada5566-a7aa-401f-bd96-45619a55120d": "1Password",
  "b84e4048-15dc-4dd0-8640-f4f60813c8af": "NordPass",
  "0ea242b4-43c4-4a1b-8b17-dd6d0b6baec6": "Keeper",
  "891494da-2c90-4d31-a9cd-4eab0aed1309": "Sésame",
  "f3809540-7f14-49c1-a8b3-8f813b225541": "Enpass",
  "b5397666-4885-aa6b-cebf-e52262a439a2": "Chromium Browser",
  "771b48fd-d3d4-4f74-9232-fc157ab0507a": "Edge on Mac",
  "39a5647e-1853-446c-a1f6-a79bae9f5bc7": "IDmelon",
  "d548826e-79b4-db40-a3d8-11116f7e8349": "Bitwarden",
  "fbfc3007-154e-4ecc-8c0b-6e020557d7bd": "iCloud Keychain",
  "53414d53-554e-4700-0000-000000000000": "Samsung Pass",
  "66a0ccb3-bd6a-191f-ee06-e375c50b9846": "Thales Bio iOS SDK",
  "8836336a-f590-0921-301d-46427531eee6": "Thales Bio Android SDK",
  "cd69adb5-3c7a-deb9-3177-6800ea6cb72a": "Thales PIN Android SDK",
  "17290f1e-c212-34d0-1423-365d729f09d9": "Thales PIN iOS SDK",
  "50726f74-6f6e-5061-7373-50726f746f6e": "Proton Pass",
  "fdb141b2-5d84-443e-8a35-4698c205a502": "KeePassXC",
  "cc45f64e-52a2-451b-831a-4edd8022a202": "ToothPic Passkey Provider",
  "bfc748bb-3429-4faa-b9f9-7cfa9f3b76d0": "iPasswords",
  "b35a26b2-8f6e-4697-ab1d-d44db4da28c6": "Zoho Vault",
  "b78a0a55-6ef8-d246-a042-ba0f6d55050c": "LastPass",
  "de503f9c-21a4-4f76-b4b7-558eb55c6f89": "Devolutions",
  "22248c4c-7a12-46e2-9a41-44291b373a4d": "LogMeOnce",
};
