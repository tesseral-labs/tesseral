import { timestampDate } from "@bufbuild/protobuf/wkt";
import { useMutation, useQuery } from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import { ArrowLeft, LoaderCircle, Trash, TriangleAlert } from "lucide-react";
import { DateTime } from "luxon";
import React, { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { Link, useNavigate, useParams } from "react-router";
import { toast } from "sonner";
import { z } from "zod";

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
  CardAction,
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
} from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";
import { NotFound } from "@/pages/NotFoundPage";

const schema = z.object({
  idpEntityId: z.string().nonempty("IDP Entity ID is required"),
  idpRedirectUrl: z.string().url("IDP Redirect URL must be a valid URL"),
  idpX509Certificate: z.string().optional(),
  primary: z.boolean(),
});

export function OrganizationSamlConnectionPage() {
  const { organizationId, samlConnectionId } = useParams();

  const {
    data: getSamlConnectionResponse,
    isError,
    isLoading,
    refetch,
  } = useQuery(
    getSAMLConnection,
    {
      id: samlConnectionId,
    },
    {
      retry: false,
    },
  );
  const updateSamlConnectionMutation = useMutation(updateSAMLConnection);

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      idpEntityId: getSamlConnectionResponse?.samlConnection?.idpEntityId || "",
      idpRedirectUrl:
        getSamlConnectionResponse?.samlConnection?.idpRedirectUrl || "",
      idpX509Certificate:
        getSamlConnectionResponse?.samlConnection?.idpX509Certificate || "",
      primary: getSamlConnectionResponse?.samlConnection?.primary || false,
    },
  });

  async function handleSubmit(data: z.infer<typeof schema>) {
    await updateSamlConnectionMutation.mutateAsync({
      id: samlConnectionId,
      samlConnection: {
        idpEntityId: data.idpEntityId,
        idpRedirectUrl: data.idpRedirectUrl,
        idpX509Certificate: data.idpX509Certificate,
        primary: data.primary,
      },
    });
    form.reset(data);
    await refetch();
    toast.success("SAML Connection updated successfully");
  }

  useEffect(() => {
    if (getSamlConnectionResponse) {
      form.reset({
        idpEntityId:
          getSamlConnectionResponse.samlConnection?.idpEntityId || "",
        idpRedirectUrl:
          getSamlConnectionResponse.samlConnection?.idpRedirectUrl || "",
        primary: getSamlConnectionResponse.samlConnection?.primary || false,
      });
    }
  }, [getSamlConnectionResponse, form]);

  return (
    <>
      {isLoading ? (
        <PageLoading />
      ) : isError ? (
        <NotFound />
      ) : (
        <PageContent>
          <Title title={`SAML Connection ${samlConnectionId}`} />

          <div>
            <Link to={`/organizations/${organizationId}/authentication`}>
              <Button variant="ghost" size="sm">
                <ArrowLeft />
                Back to Authentication
              </Button>
            </Link>
          </div>

          <div>
            <div>
              <h1 className="text-2xl font-semibold">SAML Connection</h1>
              <ValueCopier
                value={getSamlConnectionResponse?.samlConnection?.id || ""}
                label="SAML Connection ID"
              />
              <div className="flex flex-wrap mt-2 gap-2 text-muted-foreground/50">
                <Badge className="border-0" variant="outline">
                  Created{" "}
                  {getSamlConnectionResponse?.samlConnection?.createTime &&
                    DateTime.fromJSDate(
                      timestampDate(
                        getSamlConnectionResponse.samlConnection.createTime,
                      ),
                    ).toRelative()}
                </Badge>
                <div>•</div>
                <Badge className="border-0" variant="outline">
                  Updated{" "}
                  {getSamlConnectionResponse?.samlConnection?.updateTime &&
                    DateTime.fromJSDate(
                      timestampDate(
                        getSamlConnectionResponse.samlConnection.updateTime,
                      ),
                    ).toRelative()}
                </Badge>
              </div>
            </div>
          </div>

          <Card>
            <CardHeader>
              <CardTitle>Service Provider Details</CardTitle>
              <CardDescription>
                The configuration here is assigned automatically by Tesseral,
                and needs to be inputted into your customer's Identity Provider
                by their IT admin.
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="space-y-2">
                <div className="font-semibold">
                  Assertion Consumer Service (ACS) URL
                </div>
                <ValueCopier
                  value={
                    getSamlConnectionResponse?.samlConnection?.spAcsUrl || ""
                  }
                  label="ACS URL"
                />
              </div>
              <div className="space-y-2">
                <div className="font-semibold">SP Entity ID</div>
                <ValueCopier
                  value={
                    getSamlConnectionResponse?.samlConnection?.spEntityId || ""
                  }
                  label="SP Entity ID"
                />
              </div>
            </CardContent>
          </Card>

          <Form {...form}>
            <form onSubmit={form.handleSubmit(handleSubmit)}>
              <Card>
                <CardHeader>
                  <CardTitle>Identity Provider settings</CardTitle>
                  <CardDescription>
                    The configuration here needs to be copied over from the
                    customer's Identity Provider ("IDP").
                  </CardDescription>
                  <CardAction>
                    <Button
                      type="submit"
                      disabled={
                        !form.formState.isDirty ||
                        updateSamlConnectionMutation.isPending
                      }
                    >
                      {updateSamlConnectionMutation.isPending && (
                        <LoaderCircle className="animate-spin" />
                      )}
                      {updateSamlConnectionMutation.isPending
                        ? "Saving changes"
                        : "Save changes"}
                    </Button>
                  </CardAction>
                </CardHeader>
                <CardContent className="space-y-6">
                  <FormField
                    control={form.control}
                    name="primary"
                    render={({ field }) => (
                      <FormItem className="flex items-center justify-between space-x-4">
                        <div className="space-y-2">
                          <FormLabel>Primary Connection</FormLabel>
                          <FormDescription>
                            A primary SAML connection gets used by default
                            within its organization.
                          </FormDescription>
                          <FormMessage />
                        </div>
                        <FormControl>
                          <Switch
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
                          The IDP Entity ID, as configured in the customer's
                          Identity Provider.
                        </FormDescription>
                        <FormMessage />
                        <FormControl>
                          <Input className="max-w-xl" {...field} />
                        </FormControl>
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
                          The IDP Redirect URL, as configured in the customer's
                          Identity Provider.
                        </FormDescription>
                        <FormMessage />
                        <FormControl>
                          <Input className="max-w-xl" {...field} />
                        </FormControl>
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
                        <FormDescription className="max-w-xl">
                          IDP Certificate, as a PEM-encoded X.509 certificate.
                          These start with '-----BEGIN CERTIFICATE-----' and end
                          with '-----END CERTIFICATE-----'.
                        </FormDescription>
                        <FormMessage />
                        <FormControl>
                          <Input
                            className="max-w-xl"
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
                        {getSamlConnectionResponse?.samlConnection
                          ?.idpX509Certificate && (
                          <FormDescription>
                            <a
                              className="font-medium underline underline-offset-2 decoration-muted-foreground/40"
                              download={`Certificate ${samlConnectionId}.crt`}
                              href={`data:text/plain;base64,${btoa(getSamlConnectionResponse.samlConnection.idpX509Certificate)}`}
                            >
                              Download Current (.crt)
                            </a>
                          </FormDescription>
                        )}
                      </FormItem>
                    )}
                  />
                </CardContent>
              </Card>
            </form>
          </Form>

          <DangerZoneCard />
        </PageContent>
      )}
    </>
  );
}

function DangerZoneCard() {
  const { organizationId, samlConnectionId } = useParams();
  const navigate = useNavigate();

  const deleteSamlConnectionMutation = useMutation(deleteSAMLConnection);

  const [deleteOpen, setDeleteOpen] = useState(false);

  async function handleDelete() {
    await deleteSamlConnectionMutation.mutateAsync({
      id: samlConnectionId,
    });
    toast.success("SAML Connection deleted successfully");
    navigate(`/organizations/${organizationId}/authentication`);
  }

  return (
    <>
      <Card className="bg-red-50/50 border-red-200">
        <CardHeader>
          <CardTitle className="text-destructive flex items-center gap-2">
            <TriangleAlert className="w-4 h-4" />
            <span>Danger Zone</span>
          </CardTitle>
          <CardDescription>
            This section contains actions that can have significant
            consequences. Proceed with caution.
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-6">
          <div className="flex items-center justify-between gap-8 w-full lg:w-auto flex-wrap lg:flex-nowrap">
            <div className="space-y-1">
              <div className="text-sm font-semibold flex items-center gap-2">
                <Trash className="w-4 h-4" />
                <span>Delete SAML Connection</span>
              </div>
              <div className="text-sm text-muted-foreground">
                Completely delete the SAML Connection. This cannot be undone.
              </div>
            </div>
            <Button
              variant="destructive"
              size="sm"
              onClick={() => setDeleteOpen(true)}
            >
              Delete SAML Connection
            </Button>
          </div>
        </CardContent>
      </Card>

      {/* Delete Confirmation Dialog */}
      <AlertDialog open={deleteOpen} onOpenChange={setDeleteOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle className="flex items-center gap-2">
              <TriangleAlert />
              <span>Are you sure?</span>
            </AlertDialogTitle>
            <AlertDialogDescription>
              This action will permanently delete the SAML Connection. This
              cannot be undone.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <Button variant="outline" onClick={() => setDeleteOpen(false)}>
              Cancel
            </Button>
            <Button variant="destructive" onClick={handleDelete}>
              Delete SAML Connection
            </Button>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </>
  );
}
