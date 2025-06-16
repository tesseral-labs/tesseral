import { useMutation, useQuery } from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import { Shield } from "lucide-react";
import React, { MouseEvent, useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { Link, useParams } from "react-router";
import { toast } from "sonner";
import { z } from "zod";

import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
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
} from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";
import { updateOrganization } from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";

export function OrganizationSamlCard() {
  const { organizationId } = useParams();
  const { data: getOrganizationResponse } = useQuery(getOrganization, {
    id: organizationId,
  });
  const { data: getProjectResponse } = useQuery(getProject);

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <Shield />
          SAML SSO
        </CardTitle>
        <CardDescription>
          Configure SAML authentication for{" "}
          <span className="font-semibold">
            {getOrganizationResponse?.organization?.displayName}
          </span>
          .
        </CardDescription>
      </CardHeader>
      <CardContent>
        <div className="space-y-4">
          {getProjectResponse?.project?.logInWithSaml ? (
            <>
              <div className="flex justify-between items-center gap-4">
                <div>
                  <div className="font-semibold text-sm">Log in with SAML</div>
                  <div className="text-xs text-muted-foreground">
                    Allows users to log into this organization with SAML-based
                    identity providers.
                  </div>
                </div>
                <Switch
                  checked={getOrganizationResponse?.organization?.logInWithSaml}
                  disabled
                />
              </div>

              <ConfigureOrganizationSaml />
            </>
          ) : (
            <>
              <div className="text-sm text-muted-foreground">
                SAML authentication is not enabled for this project. Please
                enable SAML at the project level to configure it for this
                organization.
              </div>

              <Link to="/settings/authentication/saml">
                <Button className="w-full" variant="outline">
                  Manage Project SAML Settings
                </Button>
              </Link>
            </>
          )}
        </div>
      </CardContent>
    </Card>
  );
}

const schema = z.object({
  logInWithSaml: z.boolean().default(false),
});

function ConfigureOrganizationSaml() {
  const { organizationId } = useParams();
  const { data: getOrganizationResponse, refetch } = useQuery(getOrganization, {
    id: organizationId,
  });
  const updateOrganizationMutation = useMutation(updateOrganization);

  const [open, setOpen] = useState(false);

  const form = useForm({
    resolver: zodResolver(schema),
    defaultValues: {
      logInWithSaml:
        getOrganizationResponse?.organization?.logInWithSaml || false,
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
        logInWithSaml: data.logInWithSaml,
      },
    });
    await refetch();
    form.reset();
    setOpen(false);
    toast.success("SAML configuration updated successfully.");
  }

  useEffect(() => {
    if (getOrganizationResponse) {
      form.reset({
        logInWithSaml: getOrganizationResponse.organization?.logInWithSaml,
      });
    }
  }, [getOrganizationResponse, form]);

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button variant="outline" className="w-full">
          Configure SAML
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Configure SAML</DialogTitle>
          <DialogDescription>
            Configure SAML authentication for{" "}
            <span className="font-semibold">
              {getOrganizationResponse?.organization?.displayName}
            </span>
            . This will allow users to log in using SAML-based identity
            providers.
          </DialogDescription>
        </DialogHeader>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(handleSubmit)}>
            <FormField
              control={form.control}
              name="logInWithSaml"
              render={({ field }) => (
                <FormItem className="flex items-center justify-between gap-4">
                  <div>
                    <FormLabel>Log in with SAML</FormLabel>
                    <FormDescription>
                      Allows users to log into this organization with SAML-based
                      identity providers.
                    </FormDescription>
                    <FormMessage />
                  </div>
                  <Switch
                    id="logInWithSaml"
                    checked={field.value}
                    onCheckedChange={field.onChange}
                  />
                </FormItem>
              )}
            />

            <DialogFooter className="mt-4 justify-end">
              <Button variant="outline" onClick={handleCancel}>
                Cancel
              </Button>
              <Button type="submit" disabled={!form.formState.isDirty}>
                Save Changes
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}
