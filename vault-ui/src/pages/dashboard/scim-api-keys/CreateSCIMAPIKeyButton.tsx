import { useInfiniteQuery, useMutation } from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import { AlertDialogTrigger } from "@radix-ui/react-alert-dialog";
import React, { useState } from "react";
import { useForm } from "react-hook-form";
import { useNavigate } from "react-router";
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
} from "@/components/ui/alert-dialog";
import { Button } from "@/components/ui/button";
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
import Loader from "@/components/ui/loader";
import {
  createSCIMAPIKey,
  listSCIMAPIKeys,
} from "@/gen/tesseral/frontend/v1/frontend-FrontendService_connectquery";

const schema = z.object({
  displayName: z.string().min(1, "Display name is required"),
});

export function CreateSCIMAPIKeyButton() {
  const navigate = useNavigate();
  const [open, setOpen] = useState(false);
  const { refetch } = useInfiniteQuery(
    listSCIMAPIKeys,
    {
      pageToken: "",
    },
    {
      pageParamKey: "pageToken",
      getNextPageParam: (page) => page.nextPageToken || undefined,
    },
  );
  const createSCIMAPIKeyMutation = useMutation(createSCIMAPIKey);

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      displayName: "",
    },
  });

  async function handleSubmit(values: z.infer<typeof schema>) {
    const { scimApiKey } = await createSCIMAPIKeyMutation.mutateAsync({
      scimApiKey: {
        displayName: values.displayName,
      },
    });

    if (scimApiKey) {
      await refetch();
      setOpen(false);
      toast.success("SCIM API Key created successfully");
      form.reset();
      navigate(`/organization-settings/scim-api-keys/${scimApiKey.id}`);
    } else {
      toast.error("Failed to create SCIM API Key");
    }
  }

  return (
    <AlertDialog open={open} onOpenChange={setOpen}>
      <AlertDialogTrigger asChild>
        <Button variant="outline">Create SCIM API Key</Button>
      </AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Create SCIM API Key</AlertDialogTitle>
          <AlertDialogDescription>
            Create a new SCIM API key to allow for enterprise directory syncing.
          </AlertDialogDescription>
        </AlertDialogHeader>

        <Form {...form}>
          <form
            onSubmit={form.handleSubmit(handleSubmit)}
            className="space-y-4"
          >
            <FormField
              control={form.control}
              name="displayName"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Display name</FormLabel>
                  <FormDescription>
                    A human-friendly name for the SCIM API Key.
                  </FormDescription>
                  <FormControl>
                    <Input {...field} placeholder="e.g. My SCIM API Key" />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <AlertDialogFooter>
              <AlertDialogCancel>Cancel</AlertDialogCancel>
              <Button
                type="submit"
                disabled={createSCIMAPIKeyMutation.isPending}
              >
                {createSCIMAPIKeyMutation.isPending && <Loader />}
                Create
              </Button>
            </AlertDialogFooter>
          </form>
        </Form>
      </AlertDialogContent>
    </AlertDialog>
  );
}
