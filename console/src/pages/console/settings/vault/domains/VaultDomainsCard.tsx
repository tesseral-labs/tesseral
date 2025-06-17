import { useMutation, useQuery } from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import { ExternalLink, Globe, LoaderCircle } from "lucide-react";
import React, { MouseEvent, useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { Link } from "react-router";
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
import { Input } from "@/components/ui/input";
import {
  getProject,
  getProjectEntitlements,
  getVaultDomainSettings,
  updateVaultDomainSettings,
} from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";

export function VaultDomainsCard() {
  const { data: getProjectResponse } = useQuery(getProject);
  const { data: getVaultDomainSettingsResponse } = useQuery(
    getVaultDomainSettings,
  );

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-x-2">
          <Globe />
          <span>Vault Domains</span>
        </CardTitle>
        <CardDescription>
          Configure the domains for your Tesseral Vault.
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-6 text-sm">
        <div className="space-y-2">
          <div className="space-y-1">
            <p className="font-semibold ">Vault Domain</p>
            <p className="text-muted-foreground text-xs">
              The domain where your Vault page are hosted.
            </p>
          </div>
          <div>
            <Link
              className="inline-flex text-muted-foreground text-xs font-mono items-center gap-2 bg-muted py-1 px-2 rounded"
              to={`https://${getProjectResponse?.project?.vaultDomain}`}
              target="_blank"
            >
              <span>
                {getProjectResponse?.project?.vaultDomain || "Not configured"}
              </span>
              <ExternalLink className="h-3 w-3" />
            </Link>
          </div>
        </div>
        <div className="space-y-2">
          <div className="space-y-1">
            <p className="font-semibold ">Email Send From Domain</p>
            <p className="text-muted-foreground text-xs">
              The domain where Vault emails will be sent from.
            </p>
          </div>
          <div className="text-xs font-mono text-muted-foreground bg-muted inline py-1 px-2 rounded">
            {getProjectResponse?.project?.emailSendFromDomain ||
              "Not configured"}
          </div>
        </div>
        <div className="space-y-2">
          <div className="space-y-1">
            <p className="font-semibold ">Pending Custom Domain</p>
            <p className="text-muted-foreground text-xs">
              A custom Vault domain that is not enabled yet, if any.
            </p>
          </div>
          <div className="text-xs font-mono text-muted-foreground bg-muted inline py-1 px-2 rounded">
            {getVaultDomainSettingsResponse?.vaultDomainSettings
              ?.pendingDomain || <span>&mdash;</span>}
          </div>
        </div>
      </CardContent>
      <CardFooter className="mt-8">
        <ConfigureVaultDomainsButton />
      </CardFooter>
    </Card>
  );
}

const schema = z.object({
  pendingDomain: z.string(),
});

function ConfigureVaultDomainsButton() {
  const { data: getProjectResponse } = useQuery(getProject);
  const { data: getProjectEntitlementsResponse } = useQuery(
    getProjectEntitlements,
    {},
  );
  const { data: getVaultDomainSettingsResponse, refetch } = useQuery(
    getVaultDomainSettings,
  );
  const updateVaultDomainSettingsMutation = useMutation(
    updateVaultDomainSettings,
  );

  const [open, setOpen] = useState(false);

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      pendingDomain: "",
    },
  });

  function handleCancel(e: MouseEvent<HTMLButtonElement>) {
    e.preventDefault();
    e.stopPropagation();
    setOpen(false);
    form.reset({
      pendingDomain:
        getVaultDomainSettingsResponse?.vaultDomainSettings?.pendingDomain,
    });
  }

  async function handleSubmit(data: z.infer<typeof schema>) {
    await updateVaultDomainSettingsMutation.mutateAsync({
      vaultDomainSettings: {
        pendingDomain: data.pendingDomain,
      },
    });
    await refetch();
    setOpen(false);
  }

  useEffect(() => {
    if (
      getVaultDomainSettingsResponse?.vaultDomainSettings &&
      getProjectResponse?.project
    ) {
      form.reset({
        pendingDomain:
          getVaultDomainSettingsResponse?.vaultDomainSettings?.pendingDomain,
      });
    }
  }, [getProjectResponse, getVaultDomainSettingsResponse, form]);

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger
        asChild
        disabled={!getProjectEntitlementsResponse?.entitledCustomVaultDomains}
      >
        <Button
          className="w-full"
          variant="outline"
          disabled={!getProjectEntitlementsResponse?.entitledCustomVaultDomains}
        >
          Configure Vault Domain
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Edit Vault Domain</DialogTitle>
          <DialogDescription>
            Configure a custom domain for your Vault.
          </DialogDescription>
        </DialogHeader>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(handleSubmit)}>
            <div className="space-y-6">
              <FormField
                control={form.control}
                name="pendingDomain"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Custom Vault Domain</FormLabel>
                    <FormDescription>
                      A custom domain for your Vault. Typically, you'll use
                      "vault.company.com", where "company.com" is your company
                      domain.
                    </FormDescription>
                    <FormMessage />
                    <FormControl>
                      <Input
                        className="max-w-96"
                        placeholder="vault.company.com"
                        {...field}
                      />
                    </FormControl>
                  </FormItem>
                )}
              />
            </div>
            <DialogFooter className="mt-8">
              <Button variant="outline" onClick={handleCancel}>
                Cancel
              </Button>
              <Button
                type="submit"
                disabled={
                  !form.formState.isDirty ||
                  updateVaultDomainSettingsMutation.isPending
                }
              >
                {updateVaultDomainSettingsMutation.isPending && (
                  <LoaderCircle className="animate-spin" />
                )}
                {updateVaultDomainSettingsMutation.isPending
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
