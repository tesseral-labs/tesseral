import { timestampDate } from "@bufbuild/protobuf/wkt";
import { useInfiniteQuery, useMutation } from "@connectrpc/connect-query";
import { LoaderCircle, Plus, Trash, TriangleAlert } from "lucide-react";
import { DateTime } from "luxon";
import React, { useState } from "react";
import { toast } from "sonner";

import { ValueCopier } from "@/components/core/ValueCopier";
import { TableSkeleton } from "@/components/skeletons/TableSkeleton";
import {
  AlertDialog,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardAction,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  deleteMyPasskey,
  getPasskeyOptions,
  listMyPasskeys,
  registerPasskey,
} from "@/gen/tesseral/frontend/v1/frontend-FrontendService_connectquery";
import { Passkey } from "@/gen/tesseral/frontend/v1/models_pb";
import { AAGUIDS } from "@/lib/passkeys";
import { base64urlEncode } from "@/lib/utils";

export function UserPasskeysCard() {
  const {
    data: listMyPasskeysResponses,
    fetchNextPage,
    hasNextPage,
    isFetchingNextPage,
    isLoading,
  } = useInfiniteQuery(
    listMyPasskeys,
    {
      pageToken: "",
    },
    {
      pageParamKey: "pageToken",
      getNextPageParam: (lastPage) => lastPage.nextPageToken || undefined,
    },
  );

  const passkeys =
    listMyPasskeysResponses?.pages.flatMap((page) => page.passkeys) || [];

  return (
    <Card>
      <CardHeader>
        <CardTitle>Passkeys</CardTitle>
        <CardDescription>
          Manage the passkeys associated with your account. Passkeys are enabled
          by your organization as a valid Multi-factor Authentication (MFA)
          method.
        </CardDescription>
        <CardAction>
          <RegisterPasskeyButton />
        </CardAction>
      </CardHeader>
      <CardContent>
        {isLoading ? (
          <TableSkeleton columns={6} />
        ) : (
          <>
            {passkeys.length === 0 ? (
              <div className="text-center text-muted-foreground text-sm pt-8">
                No passkeys registered.
              </div>
            ) : (
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Passkey</TableHead>
                    <TableHead>Status</TableHead>
                    <TableHead>Public Key</TableHead>
                    <TableHead>Created</TableHead>
                    <TableHead>Updated</TableHead>
                    <TableHead className="text-right">Actions</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {passkeys.map((passkey) => (
                    <TableRow key={passkey.id}>
                      <TableCell className="space-y-2">
                        <span className="block font-semibold">
                          {passkey.aaguid && AAGUIDS[passkey.aaguid]
                            ? AAGUIDS[passkey.aaguid]
                            : "Unknown Vendor"}
                        </span>
                        <ValueCopier value={passkey.id} label="Passkey ID" />
                      </TableCell>
                      <TableCell>
                        {passkey.disabled ? (
                          <Badge variant="secondary">Disabled</Badge>
                        ) : (
                          <Badge>Active</Badge>
                        )}
                      </TableCell>
                      <TableCell>
                        {passkey.publicKeyPkix && (
                          <a
                            className="font-medium underline underline-offset-2 decoration-muted-foreground/40"
                            download={`Public Key ${passkey.id}.pem`}
                            href={`data:text/plain;base64,${btoa(passkey.publicKeyPkix)}`}
                          >
                            Download (.pem)
                          </a>
                        )}
                      </TableCell>
                      <TableCell>
                        {passkey.createTime &&
                          DateTime.fromJSDate(
                            timestampDate(passkey.createTime),
                          ).toRelative()}
                      </TableCell>
                      <TableCell>
                        {passkey.updateTime &&
                          DateTime.fromJSDate(
                            timestampDate(passkey.updateTime),
                          ).toRelative()}
                      </TableCell>
                      <TableCell className="text-right">
                        <DeletePasskeyButton passkey={passkey} />
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            )}
          </>
        )}
      </CardContent>
      {hasNextPage && (
        <CardFooter className="justify-center">
          <Button
            variant="outline"
            size="sm"
            onClick={() => fetchNextPage()}
            disabled={isFetchingNextPage}
          >
            Load more
          </Button>
        </CardFooter>
      )}
    </Card>
  );
}

function DeletePasskeyButton({ passkey }: { passkey: Passkey }) {
  const { refetch } = useInfiniteQuery(
    listMyPasskeys,
    {
      pageToken: "",
    },
    {
      pageParamKey: "pageToken",
      getNextPageParam: (lastPage) => lastPage.nextPageToken || undefined,
    },
  );
  const deletePasskeyMutation = useMutation(deleteMyPasskey);

  const [deleteOpen, setDeleteOpen] = useState(false);

  async function handleDeletePasskey() {
    try {
      await deletePasskeyMutation.mutateAsync({ id: passkey.id });
      await refetch();
      setDeleteOpen(false);
      toast.success("Passkey deleted successfully");
    } catch {
      toast.error("Failed to delete passkey. Please try again.");
    }
  }

  return (
    <AlertDialog open={deleteOpen} onOpenChange={setDeleteOpen}>
      <AlertDialogTrigger asChild>
        <Button variant="outline" size="sm">
          <Trash />
          Delete Passkey
        </Button>
      </AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle className="flex items-center gap-2">
            <TriangleAlert />
            Are you sure?
          </AlertDialogTitle>
          <AlertDialogDescription>
            This will permanently delete this passkey. This cannot be undone.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <Button variant="outline" onClick={() => setDeleteOpen(false)}>
            Cancel
          </Button>
          <Button variant="destructive" onClick={handleDeletePasskey}>
            Delete Passkey
          </Button>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
}

function RegisterPasskeyButton() {
  const { refetch } = useInfiniteQuery(
    listMyPasskeys,
    {
      pageToken: "",
    },
    {
      pageParamKey: "pageToken",
      getNextPageParam: (lastPage) => lastPage.nextPageToken,
    },
  );
  const getPasskeyOptionsMutation = useMutation(getPasskeyOptions);
  const registerPasskeyMutation = useMutation(registerPasskey);

  const [open, setOpen] = useState(false);

  async function handleRegisterPasskey() {
    const passkeyOptions = await getPasskeyOptionsMutation.mutateAsync({});
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

    await registerPasskeyMutation.mutateAsync({
      rpId: passkeyOptions.rpId,
      attestationObject: base64urlEncode(
        (credential.response as AuthenticatorAttestationResponse)
          .attestationObject,
      ),
    });

    await refetch();
    setOpen(false);
    toast.success("Passkey registered");
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button onClick={() => handleRegisterPasskey()} size="sm">
          <Plus />
          Register Passkey
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Registering Passkey</DialogTitle>
          <DialogDescription>
            Please follow the instructions on your device to register a new
            passkey.
          </DialogDescription>
        </DialogHeader>
        <div className="justify-center items-center p-y-16">
          <LoaderCircle className="animate-spin" />
        </div>
      </DialogContent>
    </Dialog>
  );
}
