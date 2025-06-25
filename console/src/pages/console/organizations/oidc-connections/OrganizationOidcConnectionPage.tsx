import { timestampDate } from "@bufbuild/protobuf/wkt";
import { useMutation, useQuery } from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import { ArrowLeft, LoaderCircle, Trash, TriangleAlert } from "lucide-react";
import { DateTime } from "luxon";
import React, { useCallback, useEffect, useState } from "react";
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
  deleteOIDCConnection,
  getOIDCConnection,
  updateOIDCConnection,
} from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";
import { NotFound } from "@/pages/NotFoundPage";

const schema = z.object({
  configurationUrl: z.string().url("Must be a valid URL"),
  issuer: z.string().url("Must be a valid URL"),
  clientId: z.string().nonempty("Client ID is required"),
  clientSecret: z.string().nonempty("Client Secret is required"),
  primary: z.boolean(),
});

// eslint-disable-next-line @typescript-eslint/no-explicit-any
function useDebounce<T extends any[]>(
  callback: (...args: T) => void,
  delay: number,
) {
  // The useCallback hook infers the type of the debounced callback.
  const debouncedCallback = useCallback(
    (...args: T) => {
      // setTimeout returns a timer ID, which can be a number or a NodeJS.Timeout object.
      // ReturnType<typeof setTimeout> correctly types it for both environments.
      const handler = setTimeout(() => {
        callback(...args);
      }, delay);

      // The cleanup function clears the timeout.
      return () => {
        clearTimeout(handler);
      };
    },
    [callback, delay],
  );

  return debouncedCallback;
}

export function OrganizationOidcConnectionPage() {
  const { organizationId, oidcConnectionId } = useParams();

  const {
    data: getOidcConnectionResponse,
    isError,
    isLoading,
    refetch,
  } = useQuery(
    getOIDCConnection,
    {
      id: oidcConnectionId,
    },
    {
      retry: 3,
    },
  );
  const updateOidcConnectionMutation = useMutation(updateOIDCConnection);

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      configurationUrl:
        getOidcConnectionResponse?.oidcConnection?.configurationUrl || "",
      issuer: getOidcConnectionResponse?.oidcConnection?.issuer || "",
      clientId: getOidcConnectionResponse?.oidcConnection?.clientId || "",
      clientSecret:
        getOidcConnectionResponse?.oidcConnection?.clientSecret || "",
      primary: getOidcConnectionResponse?.oidcConnection?.primary || false,
    },
  });

  async function handleSubmit(data: z.infer<typeof schema>) {
    await updateOidcConnectionMutation.mutateAsync({
      id: oidcConnectionId,
      oidcConnection: {
        configurationUrl: data.configurationUrl,
        issuer: data.issuer,
        clientId: data.clientId,
        clientSecret: data.clientSecret,
        primary: data.primary,
      },
    });
    form.reset(data);
    await refetch();
    toast.success("OIDC Connection updated successfully");
  }

  useEffect(() => {
    if (getOidcConnectionResponse) {
      form.reset({
        configurationUrl:
          getOidcConnectionResponse.oidcConnection?.configurationUrl || "",
        issuer: getOidcConnectionResponse.oidcConnection?.issuer || "",
        clientId: getOidcConnectionResponse.oidcConnection?.clientId || "",
        clientSecret:
          getOidcConnectionResponse.oidcConnection?.clientSecret || "",
        primary: getOidcConnectionResponse.oidcConnection?.primary || false,
      });
    }
  }, [getOidcConnectionResponse, form]);

  const configurationUrl = form.watch("configurationUrl");
  const fetchConfiguration = useCallback(
    async (url: string) => {
      const urlData = URL.parse(url);
      if (!urlData) {
        return;
      }

      // Parse out client ID for known OIDC providers
      let clientId: string | null = null;
      // Okta
      clientId = urlData.searchParams.get("client_id");
      // Entra
      const entraRegex =
        /^https:\/\/login\.microsoftonline\.com\/(.*?)\/v2\.0\/\.well-known\/openid-configuration$/;
      const entraMatch = url.match(entraRegex);
      if (entraMatch && entraMatch[1]) {
        clientId = entraMatch[1];
      }

      try {
        const response = await fetch(url);
        if (!response.ok) {
          throw new Error("Network response was not ok");
        }
        const data: { issuer: string } = await response.json();

        if (data.issuer && !form.getValues("issuer")) {
          form.setValue("issuer", data.issuer, { shouldValidate: true });
        }
        if (clientId && !form.getValues("clientId")) {
          form.setValue("clientId", clientId, { shouldValidate: true });
        }
        form.clearErrors("configurationUrl");
      } catch (error) {
        form.setError("configurationUrl", {
          type: "manual",
          message: "Failed to fetch OIDC configuration. Please check the URL.",
        });
        console.error("Failed to fetch configuration:", error);
      }
    },
    [form],
  );

  const debouncedFetch = useDebounce(fetchConfiguration, 500);

  useEffect(() => {
    debouncedFetch(configurationUrl);
  }, [configurationUrl, debouncedFetch]);

  return (
    <>
      {isLoading ? (
        <PageLoading />
      ) : isError ? (
        <NotFound />
      ) : (
        <PageContent>
          <Title title={`OIDC Connection ${oidcConnectionId}`} />

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
              <h1 className="text-2xl font-semibold">OIDC Connection</h1>
              <ValueCopier
                value={getOidcConnectionResponse?.oidcConnection?.id || ""}
                label="OIDC Connection ID"
              />
              <div className="flex flex-wrap mt-2 gap-2 text-muted-foreground/50">
                <Badge className="border-0" variant="outline">
                  Created{" "}
                  {getOidcConnectionResponse?.oidcConnection?.createTime &&
                    DateTime.fromJSDate(
                      timestampDate(
                        getOidcConnectionResponse.oidcConnection.createTime,
                      ),
                    ).toRelative()}
                </Badge>
                <div>•</div>
                <Badge className="border-0" variant="outline">
                  Updated{" "}
                  {getOidcConnectionResponse?.oidcConnection?.updateTime &&
                    DateTime.fromJSDate(
                      timestampDate(
                        getOidcConnectionResponse.oidcConnection.updateTime,
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
                <div className="font-semibold">Redirect URL</div>
                <ValueCopier
                  value={
                    getOidcConnectionResponse?.oidcConnection?.redirectUri || ""
                  }
                  label="Redirect URL"
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
                        updateOidcConnectionMutation.isPending
                      }
                    >
                      {updateOidcConnectionMutation.isPending && (
                        <LoaderCircle className="animate-spin" />
                      )}
                      {updateOidcConnectionMutation.isPending
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
                            A primary OIDC connection gets used by default
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
                    name="configurationUrl"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>OIDC Configuration URL</FormLabel>
                        <FormDescription>
                          The OIDC Configuration URL, as configured in the
                          customer's Identity Provider.
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
                    name="issuer"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>OIDC Issuer</FormLabel>
                        <FormDescription>
                          The OIDC Issuer, as configured in the customer's
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
                    name="clientId"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>OIDC Client ID</FormLabel>
                        <FormDescription>
                          The OIDC Client ID, as configured in the customer's
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
                    name="clientSecret"
                    render={({
                      field: { onChange },
                    }: {
                      field: { onChange: (value: string) => void };
                    }) => (
                      <FormItem>
                        <FormLabel>OIDC Client Secret</FormLabel>
                        <FormDescription>
                          The OIDC Client Secret, as configured in the
                          customer's Identity Provider.
                        </FormDescription>
                        <FormMessage />
                        <FormControl>
                          <Input
                            type="password"
                            className="max-w-xl"
                            onChange={(e) => onChange(e.target.value)}
                          />
                        </FormControl>
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
  const { organizationId, oidcConnectionId } = useParams();
  const navigate = useNavigate();

  const deleteOidcConnectionMutation = useMutation(deleteOIDCConnection);

  const [deleteOpen, setDeleteOpen] = useState(false);

  async function handleDelete() {
    await deleteOidcConnectionMutation.mutateAsync({
      id: oidcConnectionId,
    });
    toast.success("OIDC Connection deleted successfully");
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
                <span>Delete OIDC Connection</span>
              </div>
              <div className="text-sm text-muted-foreground">
                Completely delete the OIDC Connection. This cannot be undone.
              </div>
            </div>
            <Button
              variant="destructive"
              size="sm"
              onClick={() => setDeleteOpen(true)}
            >
              Delete OIDC Connection
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
              This action will permanently delete the OIDC Connection. This
              cannot be undone.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <Button variant="outline" onClick={() => setDeleteOpen(false)}>
              Cancel
            </Button>
            <Button variant="destructive" onClick={handleDelete}>
              Delete OIDC Connection
            </Button>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </>
  );
}
