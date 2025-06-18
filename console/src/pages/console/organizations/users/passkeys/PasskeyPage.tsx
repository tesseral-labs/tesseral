import { timestampDate } from "@bufbuild/protobuf/wkt";
import { useMutation, useQuery } from "@connectrpc/connect-query";
import {
  ArrowLeft,
  ShieldBan,
  ShieldPlus,
  Trash,
  TriangleAlert,
} from "lucide-react";
import { DateTime } from "luxon";
import React, { useState } from "react";
import { Link, useNavigate, useParams } from "react-router";
import { toast } from "sonner";

import { ValueCopier } from "@/components/core/ValueCopier";
import { PageContent } from "@/components/page";
import { PageLoading } from "@/components/page/PageLoading";
import { Title } from "@/components/page/Title";
import {
  AlertDialog,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";
import {
  deletePasskey,
  getPasskey,
  getUser,
  updatePasskey,
} from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";
import { AAGUIDS } from "@/lib/passkeys";
import { NotFound } from "@/pages/NotFoundPage";

export function PasskeyPage() {
  const { organizationId, passkeyId, userId } = useParams();

  const {
    data: getPasskeyResponse,
    isError,
    isLoading,
  } = useQuery(
    getPasskey,
    {
      id: passkeyId,
    },
    {
      retry: false,
    },
  );

  return (
    <>
      {isLoading ? (
        <PageLoading />
      ) : isError ? (
        <NotFound />
      ) : (
        <PageContent>
          <Title title={`Passkey ${passkeyId}`} />

          <div>
            <Link
              to={`/organizations/${organizationId}/users/${userId}/passkeys`}
            >
              <Button variant="ghost" size="sm">
                <ArrowLeft />
                Back to Passkeys
              </Button>
            </Link>
          </div>

          <div>
            <h1 className="text-2xl font-semibold">Passkey</h1>
            <ValueCopier
              value={getPasskeyResponse?.passkey?.id || ""}
              label="Passkey ID"
            />
            <div className="flex flex-wrap mt-2 gap-2 text-muted-foreground/50">
              <Badge className="border-0" variant="outline">
                Created{" "}
                {getPasskeyResponse?.passkey?.createTime &&
                  DateTime.fromJSDate(
                    timestampDate(getPasskeyResponse.passkey.createTime),
                  ).toRelative()}
              </Badge>
              <div>•</div>
              <Badge className="border-0" variant="outline">
                Updated{" "}
                {getPasskeyResponse?.passkey?.updateTime &&
                  DateTime.fromJSDate(
                    timestampDate(getPasskeyResponse.passkey.updateTime),
                  ).toRelative()}
              </Badge>
            </div>
          </div>

          <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
            <Card className="col-span-1">
              <CardHeader>
                <CardTitle>Basic Details</CardTitle>
                <CardDescription>
                  Basic information about this Passkey.
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="flex items-center justify-between gap-4">
                  <span className="text-sm font-semibold">Vendor</span>
                  <Badge variant="outline">
                    {getPasskeyResponse?.passkey?.aaguid
                      ? AAGUIDS[getPasskeyResponse.passkey.aaguid] || "Unknown"
                      : "Unknown"}
                  </Badge>
                </div>
                <div className="flex items-center justify-between gap-4">
                  <span className="text-sm font-semibold">Status</span>
                  {getPasskeyResponse?.passkey?.disabled ? (
                    <Badge variant="secondary">Disabled</Badge>
                  ) : (
                    <Badge>Active</Badge>
                  )}
                </div>
                <div className="flex items-center justify-between gap-4">
                  <span className="text-sm font-semibold">Created</span>
                  <Badge className="border-0" variant="outline">
                    {getPasskeyResponse?.passkey?.createTime &&
                      DateTime.fromJSDate(
                        timestampDate(getPasskeyResponse.passkey.createTime),
                      ).toRelative()}
                  </Badge>
                </div>
                <div className="flex items-center justify-between gap-4">
                  <span className="text-sm font-semibold">Update</span>
                  <Badge className="border-0" variant="outline">
                    {getPasskeyResponse?.passkey?.updateTime &&
                      DateTime.fromJSDate(
                        timestampDate(getPasskeyResponse.passkey.updateTime),
                      ).toRelative()}
                  </Badge>
                </div>
              </CardContent>
            </Card>
            <Card className="col-span-1">
              <CardHeader>
                <CardTitle>Advanced Details</CardTitle>
                <CardDescription>
                  Advanced information about this Passkey typically required
                  when debugging issues with Passkeys.
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="flex items-center justify-between gap-4">
                  <span className="text-sm font-semibold">Public Key</span>
                  {getPasskeyResponse?.passkey?.publicKeyPkix && (
                    <a
                      className="font-medium underline underline-offset-2 decoration-muted-foreground/40 text-sm"
                      download={`Public Key ${passkeyId}.pem`}
                      href={`data:text/plain;base64,${btoa(getPasskeyResponse.passkey.publicKeyPkix)}`}
                    >
                      Download (.pem)
                    </a>
                  )}
                </div>
                <div className="flex items-center justify-between gap-4">
                  <span className="text-sm font-semibold">AAGUID</span>
                  <ValueCopier
                    value={getPasskeyResponse?.passkey?.aaguid || ""}
                    label="User ID"
                  />
                </div>
                <div className="flex items-center justify-between gap-4">
                  <span className="text-sm font-semibold">Credential ID</span>
                  {getPasskeyResponse?.passkey?.credentialId && (
                    <ValueCopier
                      value={Array.from(getPasskeyResponse.passkey.credentialId)
                        .map((byte) => byte.toString(16).padStart(2, "0"))
                        .join("")}
                      label="Credential ID"
                      maxLength={32}
                    />
                  )}
                </div>
              </CardContent>
            </Card>
          </div>

          <DangerZoneCard />
        </PageContent>
      )}
    </>
  );
}

function DangerZoneCard() {
  const { organizationId, passkeyId, userId } = useParams();
  const navigate = useNavigate();

  const { data: getPasskeyResponse, refetch } = useQuery(getPasskey, {
    id: passkeyId,
  });
  const { data: getUserResponse } = useQuery(getUser, {
    id: userId,
  });
  const deletePasskeyMutation = useMutation(deletePasskey);
  const updatePasskeyMutation = useMutation(updatePasskey);

  const [deleteOpen, setDeleteOpen] = useState(false);
  const [disableOpen, setDisableOpen] = useState(false);
  const [enableOpen, setEnableOpen] = useState(false);

  async function handleDelete() {
    await deletePasskeyMutation.mutateAsync({ id: passkeyId });
    toast.success("Passkey deleted successfully");
    navigate(`/organizations/${organizationId}/users/${userId}/passkeys`);
  }

  async function handleDisable() {
    await updatePasskeyMutation.mutateAsync({
      id: passkeyId,
      passkey: {
        disabled: true,
      },
    });
    await refetch();
    toast.success("Passkey disabled successfully");
    setDisableOpen(false);
  }

  async function handleEnable() {
    await updatePasskeyMutation.mutateAsync({
      id: passkeyId,
      passkey: {
        disabled: false,
      },
    });
    await refetch();
    toast.success("Passkey enabled successfully");
    setEnableOpen(false);
  }

  return (
    <>
      <Card className="bg-red-50/50 border-red-200">
        <CardHeader>
          <CardTitle className="flex items-center gap-2 text-destructive">
            <TriangleAlert className="w-4 h-4" />
            <span>Danger Zone</span>
          </CardTitle>
          <CardDescription>
            This section contains actions that can have significant
            consequences. Proceed with caution.
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-6">
          {getPasskeyResponse?.passkey?.disabled ? (
            <div className="flex items-center justify-between gap-8">
              <div className="space-y-1">
                <div className="text-sm font-semibold flex items-center gap-2">
                  <ShieldPlus className="w-4 h-4" />
                  <span>Enable Passkey</span>
                </div>
                <div className="text-sm text-muted-foreground">
                  Enable this passkey.{" "}
                  <span className="font-semibold">
                    {getUserResponse?.user?.email}
                  </span>{" "}
                  will be required to authenticate with this passkey (or another
                  active passkey) when logging in.
                </div>
              </div>
              <Button
                className="border-destructive text-destructive hover:bg-destructive hover:text-white"
                variant="outline"
                size="sm"
                onClick={() => setEnableOpen(true)}
              >
                Enable Passkey
              </Button>
            </div>
          ) : (
            <div className="flex items-center justify-between gap-8">
              <div className="space-y-1">
                <div className="text-sm font-semibold flex items-center gap-2">
                  <ShieldBan className="w-4 h-4" />
                  <span>Disable Passkey</span>
                </div>
                <div className="text-sm text-muted-foreground">
                  Disable this passkey.{" "}
                  <span className="font-semibold">
                    {getUserResponse?.user?.email}
                  </span>{" "}
                  will not be able to authenticate with this passkey until it is
                  enabled again.
                </div>
              </div>
              <Button
                className="border-destructive text-destructive hover:bg-destructive hover:text-white"
                variant="outline"
                size="sm"
                onClick={() => setDisableOpen(true)}
              >
                Disable Passkey
              </Button>
            </div>
          )}
          <Separator />
          <div className="flex items-center justify-between gap-8">
            <div className="space-y-1">
              <div className="text-sm font-semibold flex items-center gap-2">
                <Trash className="w-4 h-4" />
                <span>Delete Passkey</span>
              </div>
              <div className="text-sm text-muted-foreground">
                Permanently delete this passkey.{" "}
                <span className="font-semibold">
                  {getUserResponse?.user?.email}
                </span>{" "}
                will not be able to authenticate with this passkey again. This
                cannot be undone.
              </div>
            </div>
            <Button
              variant="destructive"
              size="sm"
              onClick={() => setDeleteOpen(true)}
            >
              Delete Passkey
            </Button>
          </div>
        </CardContent>
      </Card>

      {/* Disable Confirmation Dialog */}
      <AlertDialog open={disableOpen} onOpenChange={setDisableOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle className="flex items-center gap-2">
              <TriangleAlert />
              <span>Are you sure?</span>
            </AlertDialogTitle>
            <AlertDialogDescription>
              This will disable this Passkey, disallowing{" "}
              <span className="font-semibold">
                {getUserResponse?.user?.email}
              </span>
              to authentication with this Passkey on future logins.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <Button variant="outline" onClick={() => setDisableOpen(false)}>
              Cancel
            </Button>
            <Button variant="destructive" onClick={handleDisable}>
              Disable
            </Button>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      {/* Enable Confirmation Dialog */}
      <AlertDialog open={enableOpen} onOpenChange={setEnableOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle className="flex items-center gap-2">
              <TriangleAlert />
              <span>Are you sure?</span>
            </AlertDialogTitle>
            <AlertDialogDescription>
              This will enable this Passkey.{" "}
              <span className="font-semibold">
                {getUserResponse?.user?.email}
              </span>{" "}
              will be required to authenticate with this passkey (or another
              active passkey) when logging in.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <Button variant="outline" onClick={() => setEnableOpen(false)}>
              Cancel
            </Button>
            <Button variant="destructive" onClick={handleEnable}>
              Enable
            </Button>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      {/* Delete Confirmation Dialog */}
      <AlertDialog open={deleteOpen} onOpenChange={setDeleteOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle className="flex items-center gap-2">
              <TriangleAlert />
              <span>Are you sure?</span>
            </AlertDialogTitle>
            <AlertDialogDescription>
              This cannot be undone. This will permanently delete the{" "}
              <span className="font-semibold">
                {getPasskeyResponse?.passkey?.id}
              </span>{" "}
              Passkey for{" "}
              <span className="font-semibold">
                {getUserResponse?.user?.email}
              </span>
              .
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <Button variant="outline" onClick={() => setDeleteOpen(false)}>
              Cancel
            </Button>
            <Button variant="destructive" onClick={handleDelete}>
              Delete
            </Button>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </>
  );
}
