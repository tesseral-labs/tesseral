import { useMutation, useQuery } from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import { ArrowRight, Crown, Key, LoaderCircle } from "lucide-react";
import React from "react";
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
  getOrganization,
  getProject,
  getProjectEntitlements,
  updateOrganization,
} from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";
import { useHandleUpgrade } from "@/hooks/use-handle-upgrade";

import { ListOrganizationApiKeysCard } from "./api-keys/ListOrganizationApiKeysCard";

const schema = z.object({
  apiKeysEnabled: z.boolean(),
});

export function OrganizationApiKeysTab() {
  const handleUpgrade = useHandleUpgrade();
  const { organizationId } = useParams();

  const { data: getOrganizationResponse, refetch } = useQuery(getOrganization, {
    id: organizationId,
  });
  const { data: getProjectResponse } = useQuery(getProject);
  const {
    data: getProjectEntitlementsResponse,
    isLoading: isLoadingEntitlements,
  } = useQuery(getProjectEntitlements);
  const updateOrganizationMutation = useMutation(updateOrganization);

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      apiKeysEnabled:
        getOrganizationResponse?.organization?.apiKeysEnabled || false,
    },
  });

  async function handleSubmit(data: z.infer<typeof schema>) {
    await updateOrganizationMutation.mutateAsync({
      id: organizationId,
      organization: {
        apiKeysEnabled: data.apiKeysEnabled,
      },
    });
    await refetch();
    form.reset(data);
    toast.success("Managed API Keys configuration updated successfully.");
  }

  return (
    <>
      {!isLoadingEntitlements &&
        !getProjectEntitlementsResponse?.entitledBackendApiKeys && (
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
            <Form {...form}>
              <form onSubmit={form.handleSubmit(handleSubmit)}>
                <Card>
                  <CardHeader>
                    <CardTitle className="flex items-center gap-2">
                      <Key />
                      Managed API Keys
                    </CardTitle>
                    <CardDescription>
                      Managed API Keys allow your customers to authenticate to
                      your service without a session.
                    </CardDescription>
                  </CardHeader>
                  <CardContent>
                    <div className="space-y-4">
                      {getProjectResponse?.project?.apiKeysEnabled ? (
                        <FormField
                          control={form.control}
                          name="apiKeysEnabled"
                          render={({ field }) => (
                            <FormItem className="flex items-center justify-between gap-4">
                              <div>
                                <FormLabel>Managed API Keys enabled</FormLabel>
                                <FormDescription>
                                  Allows Organizations to created Managed API
                                  Keys to authenticate to your service.
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
                      ) : (
                        <div className="flex justify-between items-center gap-4">
                          <div>
                            <div className="font-semibold text-sm">
                              Managed API Keys
                            </div>
                            <div className="text-xs text-muted-foreground">
                              Managed API Keys are not available for your
                              Project.
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
                    <Button
                      className="w-full"
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
                  </CardFooter>
                </Card>
              </form>
            </Form>
          </div>

          {getOrganizationResponse?.organization?.apiKeysEnabled && (
            <ListOrganizationApiKeysCard />
          )}
        </div>
      )}
    </>
  );
}
