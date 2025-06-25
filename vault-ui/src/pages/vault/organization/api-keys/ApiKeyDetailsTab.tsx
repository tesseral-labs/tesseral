import { useMutation, useQuery } from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import React, { useEffect } from "react";
import { useForm } from "react-hook-form";
import { useParams } from "react-router";
import { toast } from "sonner";
import { z } from "zod";

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
import {
  getAPIKey,
  updateAPIKey,
} from "@/gen/tesseral/frontend/v1/frontend-FrontendService_connectquery";

const schema = z.object({
  displayName: z.string().min(1, "Display name is required"),
});

export function ApiKeyDetailsTab() {
  const { apiKeyId } = useParams();

  const { data: getApiKeyResponse, refetch } = useQuery(getAPIKey, {
    id: apiKeyId,
  });
  const updateApiKeyMutation = useMutation(updateAPIKey);

  const apiKey = getApiKeyResponse?.apiKey;

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      displayName: "",
    },
  });

  async function handleSubmit(data: z.infer<typeof schema>) {
    try {
      await updateApiKeyMutation.mutateAsync({
        id: apiKeyId,
        apiKey: {
          displayName: data.displayName,
        },
      });
      await refetch();
      form.reset();
      toast.success("API key details updated successfully.");
    } catch {
      toast.error("Failed to update API key details. Please try again.");
    }
  }

  useEffect(() => {
    if (apiKey) {
      form.reset({
        displayName: apiKey.displayName || "",
      });
    }
  }, [apiKey, form]);

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(handleSubmit)}>
        <Card>
          <CardHeader>
            <CardTitle>API Key Details</CardTitle>
            <CardDescription></CardDescription>
            <CardAction>
              <Button
                size="sm"
                type="submit"
                disabled={
                  !form.formState.isDirty || updateApiKeyMutation.isPending
                }
              >
                Save changes
              </Button>
            </CardAction>
          </CardHeader>
          <CardContent>
            <FormField
              control={form.control}
              name="displayName"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Display Name</FormLabel>
                  <FormDescription>
                    A human-readable name for the API key.
                  </FormDescription>
                  <FormMessage />
                  <FormControl>
                    <Input
                      {...field}
                      className="max-w-lg"
                      placeholder="Enter display name"
                    />
                  </FormControl>
                </FormItem>
              )}
            />
          </CardContent>
        </Card>
      </form>
    </Form>
  );
}
