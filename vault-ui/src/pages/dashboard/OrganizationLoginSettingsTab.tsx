import { timestampDate } from "@bufbuild/protobuf/wkt";
import {
  useInfiniteQuery,
  useMutation,
  useQuery,
} from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import { CopyIcon } from "lucide-react";
import { DateTime } from "luxon";
import React, { useState } from "react";
import { useForm } from "react-hook-form";
import { Link } from "react-router-dom";
import { toast } from "sonner";
import { z } from "zod";

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
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Switch } from "@/components/ui/switch";
import {
  createSAMLConnection,
  deleteRole,
  deleteSAMLConnection,
  getOrganization,
  getProject,
  listRoles,
  listSAMLConnections,
  updateSAMLConnection,
  whoami,
} from "@/gen/tesseral/frontend/v1/frontend-FrontendService_connectquery";
import { Role, SAMLConnection } from "@/gen/tesseral/frontend/v1/models_pb";
import { EditAuthenticationMethodsButton } from "@/pages/dashboard/EditAuthenticationMethodsButton";
import { EditSecondaryAuthenticationSettingsButton } from "@/pages/dashboard/EditSecondaryAuthenticationSettingsButton";

export function OrganizationLoginSettingsTab() {
  return (
    <div className="space-y-4">
      <AuthenticationMethodsCard />
      <SecondaryAuthenticationMethodsCard />
      <SAMLConnectionsCard />
      <RolesCard />
    </div>
  );
}

function AuthenticationMethodsCard() {
  const { data: getProjectResponse } = useQuery(getProject);
  const { data: getOrganizationResponse } = useQuery(getOrganization);

  return (
    <Card>
      <div className="flex items-center justify-between">
        <CardHeader>
          <CardTitle>Authentication settings</CardTitle>
          <CardDescription>
            How users in this organization can authenticate.
          </CardDescription>
        </CardHeader>

        <div className="pr-6">
          <EditAuthenticationMethodsButton />
        </div>
      </div>
      <CardContent>
        <div className="space-y-4">
          {getProjectResponse?.project?.logInWithGoogle && (
            <div>
              <div className="text-sm font-medium">Log in with Google</div>
              <div className="text-sm">
                {getOrganizationResponse?.organization?.logInWithGoogle
                  ? "Enabled"
                  : "Disabled"}
              </div>
            </div>
          )}

          {getProjectResponse?.project?.logInWithMicrosoft && (
            <div>
              <div className="text-sm font-medium">Log in with Microsoft</div>
              <div className="text-sm">
                {getOrganizationResponse?.organization?.logInWithMicrosoft
                  ? "Enabled"
                  : "Disabled"}
              </div>
            </div>
          )}

          {getProjectResponse?.project?.logInWithGithub && (
            <div>
              <div className="text-sm font-medium">Log in with GitHub</div>
              <div className="text-sm">
                {getOrganizationResponse?.organization?.logInWithGithub
                  ? "Enabled"
                  : "Disabled"}
              </div>
            </div>
          )}

          {getProjectResponse?.project?.logInWithEmail && (
            <div>
              <div className="text-sm font-medium">
                Log in with Email (Magic Links)
              </div>
              <div className="text-sm">
                {getOrganizationResponse?.organization?.logInWithEmail
                  ? "Enabled"
                  : "Disabled"}
              </div>
            </div>
          )}

          {getProjectResponse?.project?.logInWithPassword && (
            <div>
              <div className="text-sm font-medium">Log in with Password</div>
              <div className="text-sm">
                {getOrganizationResponse?.organization?.logInWithPassword
                  ? "Enabled"
                  : "Disabled"}
              </div>
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  );
}

function SecondaryAuthenticationMethodsCard() {
  const { data: getProjectResponse } = useQuery(getProject);
  const { data: getOrganizationResponse } = useQuery(getOrganization);

  if (
    !(
      getProjectResponse?.project?.logInWithPasskey ||
      getProjectResponse?.project?.logInWithAuthenticatorApp
    )
  ) {
    return null;
  }

  return (
    <Card>
      <div className="flex items-center justify-between">
        <CardHeader>
          <CardTitle>Multi-factor authentication settings</CardTitle>
          <CardDescription>
            Settings related to multi-factor authentication.
          </CardDescription>
        </CardHeader>

        <div className="pr-6">
          <EditSecondaryAuthenticationSettingsButton />
        </div>
      </div>
      <CardContent>
        <div className="space-y-4">
          <div>
            <div className="text-sm font-medium">Require MFA</div>
            <div className="text-sm">
              {getOrganizationResponse?.organization?.requireMfa
                ? "MFA Required"
                : "MFA Optional"}
            </div>
          </div>

          {getProjectResponse?.project?.logInWithAuthenticatorApp && (
            <div>
              <div className="text-sm font-medium">
                Log in with Authenticator App (Secondary Factor)
              </div>
              <div className="text-sm">
                {getOrganizationResponse?.organization
                  ?.logInWithAuthenticatorApp
                  ? "Enabled"
                  : "Disabled"}
              </div>
            </div>
          )}

          {getProjectResponse?.project?.logInWithPasskey && (
            <div>
              <div className="text-sm font-medium">
                Log in with Passkey (Secondary Factor)
              </div>
              <div className="text-sm">
                {getOrganizationResponse?.organization?.logInWithPasskey
                  ? "Enabled"
                  : "Disabled"}
              </div>
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  );
}

function SAMLConnectionsCard() {
  const { data: whoamiResponse } = useQuery(whoami);

  const { data: getOrganizationResponse } = useQuery(getOrganization);
  const { data: listSAMLConnectionsResponses } = useInfiniteQuery(
    listSAMLConnections,
    {
      pageToken: "",
    },
    {
      pageParamKey: "pageToken",
      getNextPageParam: (page) => page.nextPageToken || undefined,
    },
  );

  const samlConnections = listSAMLConnectionsResponses?.pages.flatMap(
    (pages) => pages.samlConnections,
  );

  if (!getOrganizationResponse?.organization?.logInWithSaml) {
    return null;
  }

  return (
    <Card>
      <div className="flex items-center justify-between">
        <CardHeader>
          <CardTitle>SAML settings</CardTitle>
          <CardDescription>
            Configure SAML connections to your identity provider.
          </CardDescription>
        </CardHeader>

        <div className="pr-6">
          {whoamiResponse?.user?.owner && <CreateSAMLConnectionButton />}
        </div>
      </div>

      <CardContent>
        {samlConnections && samlConnections.length > 0 ? (
          <div className="space-y-4">
            {samlConnections.map((samlConnection) => (
              <SAMLConnectionRow
                key={samlConnection.id}
                samlConnection={samlConnection}
              />
            ))}
          </div>
        ) : (
          <div className="flex justify-center">
            <div className="flex flex-col items-center gap-y-2">
              <div className="text-sm">No SAML connections configured.</div>
              {whoamiResponse?.user?.owner && <CreateSAMLConnectionButton />}
            </div>
          </div>
        )}
      </CardContent>
    </Card>
  );
}

function SAMLConnectionRow({
  samlConnection,
}: {
  samlConnection: SAMLConnection;
}) {
  const { data: whoamiResponse } = useQuery(whoami);

  const { refetch: refetchListSAMLConnections } = useInfiniteQuery(
    listSAMLConnections,
    {
      pageToken: "",
    },
    {
      pageParamKey: "pageToken",
      getNextPageParam: (page) => page.nextPageToken || undefined,
    },
  );

  const { mutateAsync: deleteSAMLConnectionAsync } =
    useMutation(deleteSAMLConnection);

  async function handleDelete() {
    await deleteSAMLConnectionAsync({
      id: samlConnection.id,
    });
    await refetchListSAMLConnections();
    toast.success("SAML connection deleted");
  }

  return (
    <div className="flex items-center justify-between">
      <div>
        <div className="text-sm font-medium flex items-center gap-x-2">
          SAML Connection {samlConnection.id}
          {samlConnection.primary && <Badge variant="outline">Primary</Badge>}
        </div>
        <div className="text-sm flex gap-x-2 text-muted-foreground">
          <span>
            Created{" "}
            {DateTime.fromJSDate(
              timestampDate(samlConnection.createTime!),
            ).toRelative()}
          </span>

          {samlConnection?.idpRedirectUrl && (
            <>
              <span>&middot;</span>
              <span>{new URL(samlConnection.idpRedirectUrl).host}</span>
            </>
          )}
        </div>
      </div>

      <div className="flex items-center gap-x-2">
        {whoamiResponse?.user?.owner && (
          <ConfigureSAMLConnectionButton samlConnection={samlConnection} />
        )}

        {whoamiResponse?.user?.owner && (
          <Button onClick={handleDelete} variant="outline">
            Delete
          </Button>
        )}
      </div>
    </div>
  );
}

function CreateSAMLConnectionButton() {
  const {
    refetch: refetchListSAMLConnections,
    data: listSAMLConnectionsResponses,
  } = useInfiniteQuery(
    listSAMLConnections,
    {
      pageToken: "",
    },
    {
      pageParamKey: "pageToken",
      getNextPageParam: (page) => page.nextPageToken || undefined,
    },
  );

  const { mutateAsync: createSAMLConnectionAsync } =
    useMutation(createSAMLConnection);

  async function handleCreate() {
    await createSAMLConnectionAsync({
      samlConnection: {
        primary:
          listSAMLConnectionsResponses?.pages?.flatMap(
            (page) => page.samlConnections,
          )?.length === 0,
      },
    });
    await refetchListSAMLConnections();
    toast.success("SAML connection created");
  }

  return (
    <Button onClick={handleCreate} variant="outline">
      Create SAML connection
    </Button>
  );
}

const schema = z.object({
  primary: z.boolean(),
  idpEntityId: z.string().min(1, {
    message: "IDP Entity ID must be non-empty.",
  }),
  idpRedirectUrl: z.string().url({
    message: "IDP Redirect URL must be a valid URL.",
  }),
  idpX509Certificate: z.string().startsWith("-----BEGIN CERTIFICATE-----", {
    message: "IDP Certificate must be a PEM-encoded X.509 certificate.",
  }),
});

function ConfigureSAMLConnectionButton({
  samlConnection,
}: {
  samlConnection: SAMLConnection;
}) {
  const { refetch: refetchListSAMLConnections } = useInfiniteQuery(
    listSAMLConnections,
    {
      pageToken: "",
    },
    {
      pageParamKey: "pageToken",
      getNextPageParam: (page) => page.nextPageToken || undefined,
    },
  );

  const [open, setOpen] = useState(false);
  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      primary: samlConnection.primary,
      idpEntityId: samlConnection.idpEntityId,
      idpRedirectUrl: samlConnection.idpRedirectUrl,
      idpX509Certificate: samlConnection.idpX509Certificate,
    },
  });

  const { mutateAsync: updateSAMLConnectionAsync } =
    useMutation(updateSAMLConnection);

  async function handleSubmit(values: z.infer<typeof schema>) {
    await updateSAMLConnectionAsync({
      id: samlConnection.id,
      samlConnection: {
        primary: values.primary,
        idpEntityId: values.idpEntityId,
        idpRedirectUrl: values.idpRedirectUrl,
        idpX509Certificate: values.idpX509Certificate,
      },
    });
    await refetchListSAMLConnections();
    toast.success("SAML connection updated");
  }

  return (
    <AlertDialog open={open} onOpenChange={setOpen}>
      <AlertDialogTrigger asChild>
        <Button variant="outline">Configure</Button>
      </AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Configure SAML connection</AlertDialogTitle>
          <AlertDialogDescription>
            Configure a SAML connection to your identity provider.
          </AlertDialogDescription>
        </AlertDialogHeader>

        <div className="truncate space-y-4">
          <div>
            <div className="text-sm font-medium">
              Assertion Consumer Service (ACS) URL
            </div>
            <div className="flex items-center gap-x-2">
              <div>
                <CopyIcon
                  onClick={() => {
                    void navigator.clipboard.writeText(samlConnection.spAcsUrl);
                    toast.success("Copied ACS URL to clipboard");
                  }}
                  className="h-3 w-3 cursor-pointer text-muted-foreground"
                />
              </div>
              <span className="text-sm truncate">
                {samlConnection.spAcsUrl}
              </span>
            </div>
          </div>

          <div>
            <div className="text-sm font-medium">SP Entity ID</div>
            <div className="flex items-center gap-x-2">
              <div>
                <CopyIcon
                  onClick={() => {
                    void navigator.clipboard.writeText(
                      samlConnection.spEntityId,
                    );
                    toast.success("Copied ACS URL to clipboard");
                  }}
                  className="h-3 w-3 cursor-pointer text-muted-foreground"
                />
              </div>
              <span className="text-sm truncate">
                {samlConnection.spEntityId}
              </span>
            </div>
          </div>
        </div>

        <Form {...form}>
          <form
            onSubmit={form.handleSubmit(handleSubmit)}
            className="space-y-4"
          >
            <FormField
              control={form.control}
              name="primary"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Primary</FormLabel>
                  <FormControl>
                    <Switch
                      className="block"
                      checked={field.value}
                      onCheckedChange={field.onChange}
                    />
                  </FormControl>
                  <FormDescription>
                    The primary SAML connection will be used for SP-initiated
                    flows. All SAML connection support IDP-initiated flows.
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="idpEntityId"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>IDP Entity ID</FormLabel>
                  <FormControl>
                    <Input {...field} />
                  </FormControl>
                  <FormDescription>Also called an "Issuer".</FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="idpRedirectUrl"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>IDP Redirect URL</FormLabel>
                  <FormControl>
                    <Input {...field} />
                  </FormControl>
                  <FormDescription>
                    Also called a "Sign On URL".
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="idpX509Certificate"
              render={({
                field: { onChange },
              }: {
                field: { onChange: (value: string) => void };
              }) => (
                <FormItem>
                  <FormLabel>IDP Certificate</FormLabel>
                  <FormControl>
                    <Input
                      type="file"
                      onChange={async (e) => {
                        // File inputs are special; they are necessarily
                        // "uncontrolled", and their value is a FileList. We
                        // just copy over the file's contents to the
                        // react-form-hook state manually on input change.
                        if (e.target.files) {
                          onChange(await e.target.files[0].text());
                        }
                      }}
                    />
                  </FormControl>
                  <FormDescription>
                    IDP Certificate. This file begins with{" "}
                    <span className="font-mono">
                      -----BEGIN CERTIFICATE-----
                    </span>
                    .
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />
            <AlertDialogFooter className="mt-8">
              <AlertDialogCancel>Cancel</AlertDialogCancel>
              <Button type="submit">Update</Button>
            </AlertDialogFooter>
          </form>
        </Form>
      </AlertDialogContent>
    </AlertDialog>
  );
}

function RolesCard() {
  const { data: listRolesResponse } = useInfiniteQuery(
    listRoles,
    {
      pageToken: "",
    },
    {
      pageParamKey: "pageToken",
      getNextPageParam: (page) => page.nextPageToken || undefined,
    },
  );

  const roles = listRolesResponse?.pages?.flatMap((page) => page.roles);

  const { data: getOrganizationResponse } = useQuery(getOrganization);
  if (!getOrganizationResponse?.organization?.customRolesEnabled) {
    return null;
  }

  return (
    <Card>
      <div className="flex items-center justify-between">
        <CardHeader>
          <CardTitle>Roles</CardTitle>
          <CardDescription>
            Roles allow you to define permissions for users in your
            organization.
          </CardDescription>
        </CardHeader>

        <div className="pr-6">
          <Link to={`/organization-settings/roles/new`}>
            <Button variant="outline">Create custom role</Button>
          </Link>
        </div>
      </div>

      <CardContent>
        <div className="space-y-4">
          {roles?.map((role) => (
            <div key={role.id} className="flex items-center justify-between">
              <div>
                <div className="text-sm font-medium">
                  <div className="text-sm font-medium flex items-center gap-x-2">
                    {role.displayName}
                    {role.organizationId ? (
                      <Badge variant="outline">Custom</Badge>
                    ) : (
                      <Badge variant="outline">Built-in</Badge>
                    )}
                  </div>
                </div>
                <div className="text-sm text-muted-foreground">
                  {role.description}
                </div>
              </div>

              {role.organizationId && (
                <div className="flex gap-x-2">
                  <Link to={`/organization-settings/roles/${role.id}/edit`}>
                    <Button variant="outline">Edit</Button>
                  </Link>

                  <DeleteRoleButton role={role} />
                </div>
              )}
            </div>
          ))}
        </div>
      </CardContent>
    </Card>
  );
}

function DeleteRoleButton({ role }: { role: Role }) {
  const { mutateAsync: deleteRoleAsync } = useMutation(deleteRole);

  const { refetch: refetchListRoles } = useInfiniteQuery(
    listRoles,
    {
      pageToken: "",
    },
    {
      pageParamKey: "pageToken",
      getNextPageParam: (page) => page.nextPageToken || undefined,
    },
  );

  async function handleDelete() {
    await deleteRoleAsync({
      id: role.id,
    });
    await refetchListRoles();
    toast.success("Role deleted");
  }

  const [open, setOpen] = useState(false);
  return (
    <AlertDialog open={open} onOpenChange={setOpen}>
      <AlertDialogTrigger asChild>
        <Button variant="outline">Delete</Button>
      </AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Confirm Deletion</AlertDialogTitle>
          <AlertDialogDescription>
            Are you sure you want to delete{" "}
            <span className="font-medium">{role.displayName}</span>? This action
            cannot be undone.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel>Cancel</AlertDialogCancel>
          <Button variant="destructive" onClick={handleDelete}>
            Delete
          </Button>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
}
