import { useMutation, useQuery } from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import { ChevronLeft } from "lucide-react";
import React, { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { useNavigate, useParams } from "react-router";
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
  deleteSAMLConnection,
  getSAMLConnection,
  updateSAMLConnection,
} from "@/gen/tesseral/frontend/v1/frontend-FrontendService_connectquery";

export function ViewSAMLConnectionPage() {
  const { samlConnectionId } = useParams();
  const { data: getSAMLConnectionResponse } = useQuery(getSAMLConnection, {
    id: samlConnectionId,
  });

  return (
    <div className="space-y-8">
      <Link to="/organization-settings/saml-connections">
        <Button variant="ghost" size="sm">
          <ChevronLeft className="h-4 w-4" />
          Back
        </Button>
      </Link>
      <Card>
        <CardHeader className="space-y-2">
          <CardTitle>Connection Details</CardTitle>
          <CardDescription>
            The following details are required to set up your SAML connection.
            Please copy these values and provide them to your Identity Provider
            (IDP) administrator.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="text-sm">
            <div className="space-y-4">
              <div>
                <h3 className="font-semibold mb-2">
                  Assertion Consumer Service (ACS) URL
                </h3>
                <div>{getSAMLConnectionResponse?.samlConnection?.spAcsUrl}</div>
              </div>
              <div>
                <h3 className="font-semibold mb-2">SP Entity ID</h3>
                <div>
                  {getSAMLConnectionResponse?.samlConnection?.spEntityId}
                </div>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader className="flex-row items-center justify-between space-x-4">
          <div className="space-y-2">
            <CardTitle>Configuration</CardTitle>
            <CardDescription>
              The following configuration is required to set up your SAML
              connection.
            </CardDescription>
          </div>
          <EditSAMLConnnectionConfigurationButton />
        </CardHeader>
        <CardContent>
          <div className="text-sm">
            <div className="space-y-4 col-span-2">
              <div>
                <h3 className="font-semibold mb-2">Primary</h3>
                <div>
                  {getSAMLConnectionResponse?.samlConnection?.primary
                    ? "Yes"
                    : "No"}
                </div>
              </div>
              <div>
                <h3 className="font-semibold mb-2">ID Entity ID</h3>
                <div>
                  {getSAMLConnectionResponse?.samlConnection?.idpEntityId ||
                    "—"}
                </div>
              </div>
              <div>
                <h3 className="font-semibold mb-2">IDP Redirect URL</h3>
                <div>
                  {getSAMLConnectionResponse?.samlConnection?.idpRedirectUrl ||
                    "—"}
                </div>
              </div>
              <div>
                <h3 className="font-semibold mb-2">IDP Certificate</h3>
                <div>
                  {getSAMLConnectionResponse?.samlConnection
                    ?.idpX509Certificate ? (
                    <a
                      className="font-medium underline underline-offset-2 decoration-muted-foreground/40"
                      download={`Certificate ${samlConnectionId}.crt`}
                      href={`data:text/plain;base64,${btoa(getSAMLConnectionResponse.samlConnection.idpX509Certificate)}`}
                    >
                      Download (.crt)
                    </a>
                  ) : (
                    "-"
                  )}
                </div>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      <DangerZoneCard />
    </div>
  );
}

const schema = z.object({
  primary: z.boolean(),
  idpEntityId: z
    .string()
    .min(1, {
      message: "IDP Entity ID must be non-empty.",
    })
    .optional()
    .or(z.literal("")),
  idpRedirectUrl: z
    .string()
    .url({
      message: "IDP Redirect URL must be a valid URL.",
    })
    .optional()
    .or(z.literal("")),
  idpX509Certificate: z
    .string()
    .startsWith("-----BEGIN CERTIFICATE-----", {
      message: "IDP Certificate must be a PEM-encoded X.509 certificate.",
    })
    .optional()
    .or(z.literal("")),
});

function EditSAMLConnnectionConfigurationButton() {
  const { samlConnectionId } = useParams();
  const [open, setOpen] = useState(false);
  const { data: getSAMLConnectionResponse, refetch } = useQuery(
    getSAMLConnection,
    {
      id: samlConnectionId,
    },
  );
  const updateSamlConnectionMutation = useMutation(updateSAMLConnection);

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      primary: false,
      idpEntityId: "",
      idpRedirectUrl: "",
      idpX509Certificate: "",
    },
  });

  async function handleSubmit(data: z.infer<typeof schema>) {
    await updateSamlConnectionMutation.mutateAsync({
      id: samlConnectionId,
      samlConnection: {
        primary: data.primary,
        idpEntityId: data.idpEntityId || undefined,
        idpRedirectUrl: data.idpRedirectUrl || undefined,
        idpX509Certificate: data.idpX509Certificate
          ? data.idpX509Certificate.toString()
          : undefined,
      },
    });

    toast.success("SAML connection updated successfully");
    await refetch();
    setOpen(false);
  }

  useEffect(() => {
    if (getSAMLConnectionResponse?.samlConnection) {
      form.reset({
        primary: getSAMLConnectionResponse.samlConnection.primary,
        idpEntityId: getSAMLConnectionResponse.samlConnection.idpEntityId || "",
        idpRedirectUrl:
          getSAMLConnectionResponse.samlConnection.idpRedirectUrl || "",
        idpX509Certificate:
          getSAMLConnectionResponse.samlConnection.idpX509Certificate || "",
      });
    }
  }, [form, getSAMLConnectionResponse]);

  return (
    <AlertDialog open={open} onOpenChange={setOpen}>
      <AlertDialogTrigger asChild>
        <Button variant="outline">Edit Configuration</Button>
      </AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>
            Edit SAML Connection Configuration
          </AlertDialogTitle>
          <AlertDialogDescription>
            Update the configuration of your SAML connection.
          </AlertDialogDescription>
        </AlertDialogHeader>

        <Form {...form}>
          <form
            onSubmit={form.handleSubmit(handleSubmit)}
            className="space-y-4"
          >
            <div className="space-y-4">
              <FormField
                control={form.control}
                name="primary"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Primary</FormLabel>
                    <FormDescription>
                      If enabled, this SAML connection will be used as the
                      primary connection for SAML authentication.
                    </FormDescription>
                    <FormControl>
                      <Switch
                        className="block"
                        checked={field.value}
                        onCheckedChange={field.onChange}
                      />
                    </FormControl>
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="idpEntityId"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>IDP Entity ID</FormLabel>
                    <FormDescription>
                      The Entity ID of your Identity Provider (IDP).
                    </FormDescription>
                    <FormControl>
                      <Input type="text" {...field} />
                    </FormControl>
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
                    <FormDescription>
                      The redirect URL for your Identity Provider (IDP).
                    </FormDescription>
                    <FormControl>
                      <Input type="text" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="idpX509Certificate"
                render={({ field: { onChange } }) => (
                  <FormItem>
                    <FormLabel>IDP X.509 Certificate</FormLabel>
                    <FormDescription>
                      This is the certificate from your IDP and is accepted as a
                      PEM-encoded X.509 certificate. These start with
                      '-----BEGIN CERTIFICATE-----' and end with '-----END
                      CERTIFICATE-----'.
                    </FormDescription>
                    <FormControl>
                      <Input
                        type="file"
                        onChange={async (e) => {
                          // File inputs are special; they are necessarily "uncontrolled", and their value is a FileList.
                          // We just copy over the file's contents to the react-form-hook state manually on input change.
                          if (e.target.files) {
                            onChange(await e.target.files[0].text());
                          }
                        }}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>

            <AlertDialogFooter>
              <AlertDialogCancel>Cancel</AlertDialogCancel>
              <Button type="submit">Save</Button>
            </AlertDialogFooter>
          </form>
        </Form>
      </AlertDialogContent>
    </AlertDialog>
  );
}

function DangerZoneCard() {
  const { samlConnectionId } = useParams();
  const [confirmDeleteOpen, setConfirmDeleteOpen] = useState(false);

  function handleDelete() {
    setConfirmDeleteOpen(true);
  }

  const deleteSAMLConnectionMutation = useMutation(deleteSAMLConnection);
  const navigate = useNavigate();
  async function handleConfirmDelete() {
    await deleteSAMLConnectionMutation.mutateAsync({
      id: samlConnectionId,
    });

    toast.success("SAML connection deleted");
    navigate(`/organization-settings/saml-connections`);
  }

  return (
    <>
      <AlertDialog open={confirmDeleteOpen} onOpenChange={setConfirmDeleteOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Delete SAML Connection?</AlertDialogTitle>
            <AlertDialogDescription>
              Deleting a SAML connection cannot be undone.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <Button variant="destructive" onClick={handleConfirmDelete}>
              Permanently Delete SAML Connection
            </Button>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      <Card className="border-destructive">
        <CardHeader>
          <CardTitle>Danger Zone</CardTitle>
        </CardHeader>

        <CardContent>
          <div className="flex justify-between items-center gap-8 w-full lg:w-auto flex-wrap lg:flex-nowrap">
            <div>
              <div className="text-sm font-semibold">
                Delete SAML Connection
              </div>
              <p className="text-sm">
                Delete this SAML connection. This cannot be undone.
              </p>
            </div>

            <Button variant="destructive" onClick={handleDelete}>
              Delete SAML Connection
            </Button>
          </div>
        </CardContent>
      </Card>
    </>
  );
}
