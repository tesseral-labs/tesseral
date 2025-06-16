import { useMutation, useQuery } from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import { ArrowRight, Crown, Key } from "lucide-react";
import React, { MouseEvent, useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { toast } from "sonner";
import { z } from "zod";

import { Badge } from "@/components/ui/badge";
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
import { Input } from "@/components/ui/input";
import { Switch } from "@/components/ui/switch";
import {
  createStripeCheckoutLink,
  getProject,
  getProjectEntitlements,
  updateProject,
} from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";

export function ManagedApiKeySettingsCard() {
  const { data: getProjectEntitlementsResponse } = useQuery(
    getProjectEntitlements,
  );
  const { data: getProjectResponse } = useQuery(getProject);
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
      {getProjectEntitlementsResponse &&
        getProjectEntitlementsResponse?.entitledBackendApiKeys && (
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Key />
                Managed API Keys
              </CardTitle>
              <CardDescription>
                Configure Organizations' ability to authenticate to your service
                using managed API keys.
              </CardDescription>
            </CardHeader>
            <CardContent className="flex-grow">
              <div className="space-y-4">
                <div className="flex justify-between gap-4">
                  <div className="font-semibold text-sm">Managed API Keys</div>
                  <div>
                    {getProjectResponse?.project?.apiKeysEnabled ? (
                      <Badge>Enabled</Badge>
                    ) : (
                      <Badge variant="secondary">Disabled</Badge>
                    )}
                  </div>
                </div>
                <div className="flex justify-between gap-4">
                  <div className="font-semibold text-sm">API Key Prefix</div>
                  <div className="font-mono text-xs p-1 rounded bg-muted text-muted-foreground">
                    {getProjectResponse?.project?.apiKeySecretTokenPrefix ||
                      "—"}
                  </div>
                </div>
              </div>
            </CardContent>
            <CardFooter className="mt-4">
              <ConfigureManagedApiKeysButton />
            </CardFooter>
          </Card>
        )}
    </>
  );
}

const schema = z.object({
  apiKeysEnabled: z.boolean(),
  apiKeySecretTokenPrefix: z.string(),
});

export function ConfigureManagedApiKeysButton() {
  const { data: getProjectResponse, refetch } = useQuery(getProject);
  const updateProjectMutation = useMutation(updateProject);

  const [open, setOpen] = useState(false);

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      apiKeysEnabled: getProjectResponse?.project?.apiKeysEnabled ?? false,
      apiKeySecretTokenPrefix:
        getProjectResponse?.project?.apiKeySecretTokenPrefix ?? "",
    },
  });

  function handleCancel(e: MouseEvent<HTMLButtonElement>) {
    e.preventDefault();
    e.stopPropagation();
    setOpen(false);
    return false;
  }

  async function handleSubmit(data: z.infer<typeof schema>) {
    await updateProjectMutation.mutateAsync({
      project: {
        apiKeysEnabled: data.apiKeysEnabled,
        apiKeySecretTokenPrefix: data.apiKeySecretTokenPrefix,
      },
    });
    await refetch();
    toast.success("Managed API Key settings updated successfully");
  }

  useEffect(() => {
    if (getProjectResponse) {
      form.reset({
        apiKeysEnabled: getProjectResponse.project?.apiKeysEnabled ?? false,
        apiKeySecretTokenPrefix:
          getProjectResponse.project?.apiKeySecretTokenPrefix ?? "",
      });
    }
  }, [getProjectResponse, form]);

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button className="w-full" variant="outline">
          Configure Managed API Keys
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Configure Managed API Keys</DialogTitle>
          <DialogDescription>
            Configure settings for managed API keys to allow customers to
            authenticate to your service without a session.
          </DialogDescription>
        </DialogHeader>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(handleSubmit)}>
            <div className="space-y-6">
              <FormField
                control={form.control}
                name="apiKeysEnabled"
                render={({ field }) => (
                  <FormItem className="flex items-center justify-between gap-4">
                    <div className="space-y-2 max-w-3xl">
                      <FormLabel>Managed API Keys Enabled</FormLabel>
                      <FormDescription>
                        Whether or not Organizations are allowed to create API
                        Keys.
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
                name="apiKeySecretTokenPrefix"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>API Key Secret Token Prefix</FormLabel>
                    <FormDescription>
                      Prefix for the API key secret token. This can be used to
                      customize the secret token your customers receive when
                      generating Managed API Keys. For example{" "}
                      <span className="text-mono text-xs px-2 py-1 bg-muted text-muted-foreground">
                        acme_sk_...
                      </span>{" "}
                      This is required if Managed API Keys are enabled.
                    </FormDescription>
                    <FormControl className="mt-2">
                      <Input
                        className="max-w-lg"
                        placeholder="acme_sk_"
                        {...field}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>
            <DialogFooter className="mt-8">
              <Button variant="outline" onClick={handleCancel}>
                Cancel
              </Button>
              <Button disabled={!form.formState.isDirty} type="submit">
                Save
              </Button>
              <FormMessage />
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}
