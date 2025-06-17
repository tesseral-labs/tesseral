import { useMutation, useQuery } from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import { ArrowRight, Crown, Key, LoaderCircle } from "lucide-react";
import React, { MouseEvent, useState } from "react";
import { useForm } from "react-hook-form";
import { useParams } from "react-router";
import { toast } from "sonner";
import { z } from "zod";

import { Button } from "@/components/ui/button";
import {
  Card,
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
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Switch } from "@/components/ui/switch";
import {
  createStripeCheckoutLink,
  getOrganization,
  getProject,
  getProjectEntitlements,
  updateOrganization,
} from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";

import { ListOrganizationApiKeysCard } from "./api-keys/ListOrganizationApiKeysCard";

export function OrganizationApiKeysTab() {
  const { organizationId } = useParams();
  const { data: getOrganizationResponse } = useQuery(getOrganization, {
    id: organizationId,
  });
  const { data: getProjectResponse } = useQuery(getProject);
  const { data: getProjectEntitlementsResponse } = useQuery(
    getProjectEntitlements,
  );
  const createStripeCheckoutLinkMutation = useMutation(
    createStripeCheckoutLink,
  );

  async function handleUpgrade() {
    const { url } = await createStripeCheckoutLinkMutation.mutateAsync({});
    window.location.href = url;
  }

  return (
    <>
      {getProjectEntitlementsResponse &&
        !getProjectEntitlementsResponse.entitledBackendApiKeys && (
          <div className="bg-gradient-to-br from-violet-500 via-purple-500 to-blue-500 border-0 text-white relative overflow-hidden shadow-xl p-8 rounded-lg">
            <div className="absolute inset-0 bg-gradient-to-br from-white/10 to-transparent" />

            <div className="flex flex-wrap w-full gap-8">
              <div className="w-full space-y-4 md:flex-grow">
                <div className="flex items-center space-x-3">
                  <div className="p-2 rounded-full bg-white/20 backdrop-blur-sm">
                    <Crown className="h-6 w-6 text-white" />
                  </div>
                  <div>
                    <h3 className="font-semibold text-white">
                      Upgrade to Growth
                    </h3>
                    <p className="text-xs text-white/80">
                      Unlock advanced features
                    </p>
                  </div>
                </div>
                <p className="font-semibold text-sm">
                  Managed API Keys are available on the Growth Tier.
                </p>
                <p className="text-sm text-white/80">
                  When you upgrade, you'll also unlock custom domains, access to
                  the Tesseral API via API keys, and dedicated email support.
                </p>
              </div>

              <div className="mt-8 md:mt-auto w-full">
                <Button
                  className="bg-white text-purple-600 hover:bg-white/90 font-medium cursor-pointer"
                  onClick={handleUpgrade}
                  size="lg"
                >
                  Upgrade Now
                  <ArrowRight className="h-4 w-4 ml-2" />
                </Button>
              </div>
            </div>
          </div>
        )}

      {getProjectEntitlementsResponse?.entitledBackendApiKeys && (
        <div className="space-y-8">
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Key />
                  Managed API Keys
                </CardTitle>
                <CardDescription>
                  Managed API Keys allow your customers to authenticate to your
                  service without a session.
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {getProjectResponse?.project?.apiKeysEnabled ? (
                    <div className="flex justify-between items-center gap-4">
                      <div>
                        <div className="font-semibold text-sm">
                          Managed API Keys enabled
                        </div>
                        <div className="text-xs text-muted-foreground">
                          Allows Organizations to created Managed API Keys to
                          authenticate to your service.
                        </div>
                      </div>
                      <Switch
                        checked={
                          getOrganizationResponse?.organization?.apiKeysEnabled
                        }
                        disabled
                      />
                    </div>
                  ) : (
                    <div className="flex justify-between items-center gap-4">
                      <div>
                        <div className="font-semibold text-sm">
                          Managed API Keys
                        </div>
                        <div className="text-xs text-muted-foreground">
                          Managed API Keys are not enabled for this
                          organization.
                        </div>
                      </div>
                      <span className="text-sm font-medium text-destructive">
                        Disabled
                      </span>
                    </div>
                  )}
                </div>
              </CardContent>
              <CardFooter className="mt-4">
                <ConfigureOrganizationApiKeys />
              </CardFooter>
            </Card>
          </div>

          <ListOrganizationApiKeysCard />
        </div>
      )}
    </>
  );
}

const schema = z.object({
  apiKeysEnabled: z.boolean(),
});

function ConfigureOrganizationApiKeys() {
  const { organizationId } = useParams();

  const { data: getOrganizationResponse, refetch } = useQuery(getOrganization, {
    id: organizationId,
  });
  const updateOrganizationMutation = useMutation(updateOrganization);

  const [open, setOpen] = useState(false);

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      apiKeysEnabled:
        getOrganizationResponse?.organization?.apiKeysEnabled || false,
    },
  });

  function handleCancel(e: MouseEvent<HTMLButtonElement>) {
    e.preventDefault();
    e.stopPropagation();
    setOpen(false);
    return false;
  }

  async function handleSubmit(data: z.infer<typeof schema>) {
    await updateOrganizationMutation.mutateAsync({
      id: organizationId,
      organization: {
        apiKeysEnabled: data.apiKeysEnabled,
      },
    });
    await refetch();
    form.reset(data);
    setOpen(false);
    toast.success("Managed API Keys configuration updated successfully.");
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button className="w-full" variant="outline">
          Configure API Keys
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Configure API Keys</DialogTitle>
          <DialogDescription>
            Configure the API Keys for{" "}
            <span className="font-semibold">
              {getOrganizationResponse?.organization?.displayName}
            </span>
            .
          </DialogDescription>
        </DialogHeader>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(handleSubmit)}>
            <div className="space-y-4">
              <FormField
                control={form.control}
                name="apiKeysEnabled"
                render={({ field }) => (
                  <FormItem className="flex items-center justify-between gap-4">
                    <div>
                      <FormLabel>Managed API Keys enabled</FormLabel>
                      <FormDescription>
                        Allows Organizations to created Managed API Keys to
                        authenticate to your service.
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
            </div>

            <DialogFooter className="mt-8 justify-end gap-2">
              <Button variant="outline" onClick={handleCancel}>
                Cancel
              </Button>
              <Button
                disabled={
                  !form.formState.isDirty ||
                  updateOrganizationMutation.isPending
                }
                type="submit"
              >
                {updateOrganizationMutation.isPending && (
                  <LoaderCircle className="animate-spin" />
                )}
                {updateOrganizationMutation.isPending
                  ? "Saving changes"
                  : "Save changes"}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}
