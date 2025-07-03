import { useMutation, useQuery } from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import { DialogDescription, DialogTrigger } from "@radix-ui/react-dialog";
import { Cookie, LoaderCircle, Settings } from "lucide-react";
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
  DialogFooter,
  DialogHeader,
  DialogTitle,
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
import {
  HoverCard,
  HoverCardContent,
  HoverCardTrigger,
} from "@/components/ui/hover-card";
import { Input } from "@/components/ui/input";
import {
  getProject,
  getProjectEntitlements,
  getVaultDomainSettings,
  updateProject,
} from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";

export function VaultCookieDomainCard() {
  const { data: getProjectResponse } = useQuery(getProject);
  const { data: getVaultDomainSettingsResponse } = useQuery(
    getVaultDomainSettings,
  );
  const { data: getProjectEntitlementsResponse } = useQuery(
    getProjectEntitlements,
    {},
  );

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <Cookie />
          <span>Vault Cookie Domain</span>
        </CardTitle>
        <CardDescription>
          Client-side JavaScript on this domain and its subdomains will have
          access to User access tokens
        </CardDescription>
      </CardHeader>
      <CardContent className="flex-grow">
        {getProjectResponse?.project?.cookieDomain ? (
          <Badge variant="outline">
            {getProjectResponse.project.cookieDomain}
          </Badge>
        ) : (
          "—"
        )}
      </CardContent>
      {getProjectEntitlementsResponse?.entitledCustomVaultDomains &&
        getVaultDomainSettingsResponse?.vaultDomainSettings?.pendingDomain && (
          <CardFooter>
            <ConfigureVaultCookieDomainButton />
          </CardFooter>
        )}
    </Card>
  );
}

const schema = z.object({
  cookieDomain: z.string().min(1, "Cookie Domain is required"),
});

function ConfigureVaultCookieDomainButton() {
  const { data: getProjectResponse, refetch } = useQuery(getProject);
  const { data: getVaultDomainSettingsResponse } = useQuery(
    getVaultDomainSettings,
  );
  const updateProjectMutation = useMutation(updateProject);

  const [open, setOpen] = useState(false);

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      cookieDomain: getProjectResponse?.project?.cookieDomain || "",
    },
  });

  function handleCancel(e: MouseEvent<HTMLButtonElement>) {
    e.preventDefault();
    e.stopPropagation();
    setOpen(false);
  }

  async function handleSubmit(data: z.infer<typeof schema>) {
    if (
      !getProjectResponse?.project?.vaultDomain?.endsWith(data.cookieDomain)
    ) {
      form.setError("cookieDomain", {
        message: `Cookie Domain must be a parent domain of the Vault domain (${getProjectResponse?.project?.vaultDomain}).`,
      });
      return;
    }

    await updateProjectMutation.mutateAsync({
      project: {
        cookieDomain: data.cookieDomain,
      },
    });
    await refetch();
    setOpen(false);
    toast.success("Vault Cookie Domain updated successfully");
  }

  useEffect(() => {
    if (getProjectResponse) {
      form.reset({
        cookieDomain: getProjectResponse.project?.cookieDomain || "",
      });
    }
  }, [getProjectResponse, form]);

  return (
    <>
      {getVaultDomainSettingsResponse && (
        <>
          {getVaultDomainSettingsResponse.vaultDomainSettings
            ?.pendingDomain && (
            <Dialog open={open} onOpenChange={setOpen}>
              <DialogTrigger asChild>
                <Button className="w-full" variant="outline">
                  <Settings />
                  Configure Vault Cookie Domain
                </Button>
              </DialogTrigger>
              <DialogContent>
                <DialogHeader>
                  <DialogTitle>Configure Vault Cookie Domain</DialogTitle>
                  <DialogDescription>
                    Set the domain for Vault cookies. This domain and its
                    subdomains will have access to User access tokens.
                  </DialogDescription>
                </DialogHeader>

                <Form {...form}>
                  <form onSubmit={form.handleSubmit(handleSubmit)}>
                    <div className="space-y-6">
                      <FormField
                        control={form.control}
                        name="cookieDomain"
                        render={({ field }) => (
                          <FormItem>
                            <FormLabel>Cookie Domain</FormLabel>
                            <FormDescription>
                              Client-side JavaScript on this domain and its
                              subdomains will have access to User access tokens.
                              You cannot modify this field until you have
                              configured a custom Vault domain.
                            </FormDescription>
                            <FormMessage />
                            <FormControl>
                              <Input
                                {...field}
                                placeholder="e.g. vault.example.com"
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
                        disabled={
                          !form.formState.isDirty ||
                          updateProjectMutation.isPending
                        }
                        type="submit"
                      >
                        {updateProjectMutation.isPending && (
                          <LoaderCircle className="animate-spin" />
                        )}
                        {updateProjectMutation.isPending
                          ? "Saving changes"
                          : "Save changes"}
                      </Button>
                    </DialogFooter>
                  </form>
                </Form>
              </DialogContent>
            </Dialog>
          )}
          {!getVaultDomainSettingsResponse.vaultDomainSettings
            ?.pendingDomain && (
            <HoverCard>
              <HoverCardTrigger className="w-full">
                <Button disabled className="w-full" variant="outline" size="sm">
                  Configure Vault Cookie Domains
                </Button>
              </HoverCardTrigger>
              <HoverCardContent className="bg-primary text-white text-xs">
                <p>
                  You cannot modify the cookie domain until you have configured
                  a custom Vault domain.
                </p>
              </HoverCardContent>
            </HoverCard>
          )}
        </>
      )}
    </>
  );
}
