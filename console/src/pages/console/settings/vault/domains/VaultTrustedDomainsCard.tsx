import { useMutation, useQuery } from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import { GlobeLock, LoaderCircle, Settings } from "lucide-react";
import React, { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { toast } from "sonner";
import { z } from "zod";

import { InputTags } from "@/components/core/InputTags";
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
import {
  getProject,
  updateProject,
} from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";

export function VaultTrustedDomainsCard() {
  const { data: getProjectResponse } = useQuery(getProject);

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <GlobeLock />
          <span>Trusted Domains</span>
        </CardTitle>
        <CardDescription>
          The domains that your app runs on, e.g. "app.company.com". The Vault
          domain is always a trusted domain.
        </CardDescription>
      </CardHeader>
      <CardContent className="flex-grow">
        {getProjectResponse && (
          <div className="flex flex-wrap gap-2">
            {(getProjectResponse.project?.trustedDomains || []).length > 0 ? (
              <>
                {getProjectResponse.project?.trustedDomains.map((domain) => (
                  <Badge key={domain} variant="outline">
                    {domain}
                  </Badge>
                ))}
              </>
            ) : (
              <div>No trusted domains</div>
            )}
          </div>
        )}
      </CardContent>
      <CardFooter>
        <ConfigureVaultTrustedDomainsButton />
      </CardFooter>
    </Card>
  );
}

const schema = z.object({
  trustedDomains: z.array(z.string()),
});

function ConfigureVaultTrustedDomainsButton() {
  const { data: getProjectResponse, refetch } = useQuery(getProject);
  const updateProjectMutation = useMutation(updateProject);

  const [open, setOpen] = useState(false);

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      trustedDomains: getProjectResponse?.project?.trustedDomains || [],
    },
  });

  function handleCancel(e: React.MouseEvent<HTMLButtonElement>) {
    e.preventDefault();
    e.stopPropagation();
    setOpen(false);
    form.reset({
      trustedDomains: getProjectResponse?.project?.trustedDomains || [],
    });
  }

  async function handleSubmit(data: z.infer<typeof schema>) {
    await updateProjectMutation.mutateAsync({
      project: {
        trustedDomains: data.trustedDomains,
      },
    });
    await refetch();
    setOpen(false);
    toast.success("Trusted domains updated successfully");
  }

  useEffect(() => {
    if (getProjectResponse?.project) {
      form.reset({
        trustedDomains: getProjectResponse.project.trustedDomains || [],
      });
    }
  }, [getProjectResponse, form]);

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button className="w-full" variant="outline">
          <Settings />
          Configure Trusted Domains
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Configure Trusted Domains</DialogTitle>
          <DialogDescription></DialogDescription>
        </DialogHeader>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(handleSubmit)}>
            <div className="space-y-8">
              <FormField
                control={form.control}
                name="trustedDomains"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Trusted Domains</FormLabel>
                    <FormDescription>
                      Add the domains that your app runs on, e.g.
                      "app.company.com". The Vault domain is always a trusted
                      domain.
                    </FormDescription>
                    <FormMessage />
                    <FormControl>
                      <InputTags
                        placeholder="app.company.com, localhost:3000"
                        {...field}
                      />
                    </FormControl>
                  </FormItem>
                )}
              />
            </div>

            <DialogFooter className="mt-8">
              <Button variant="outline" onClick={handleCancel} type="button">
                Cancel
              </Button>
              <Button
                disabled={
                  !form.formState.isDirty || updateProjectMutation.isPending
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
  );
}
