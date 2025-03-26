import { useMutation } from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import { LoaderCircleIcon } from "lucide-react";
import React, { useState } from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { LoginFlowCard } from "../../components/login/LoginFlowCard";
import { Button } from "../../components/ui/button";
import { CardContent, CardHeader, CardTitle } from "../../components/ui/card";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "../../components/ui/form";
import { Input } from "../../components/ui/input";
import { createOrganization } from "../../gen/tesseral/intermediate/v1/intermediate-IntermediateService_connectquery";
import { useRedirectNextLoginFlowPage } from "../../hooks/use-redirect-next-login-flow-page";
import { useDarkMode } from "../../lib/dark-mode";

const schema = z.object({
  displayName: z.string().nonempty(),
});

export function CreateOrganizationPage() {
  const darkMode = useDarkMode();
  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      displayName: "",
    },
  });

  const [submitting, setSubmitting] = useState(false);
  const { mutateAsync: createOrganizationAsync } =
    useMutation(createOrganization);
  const redirectNextLoginFlowPage = useRedirectNextLoginFlowPage();

  async function handleSubmit(values: z.infer<typeof schema>) {
    setSubmitting(true);

    await createOrganizationAsync({
      displayName: values.displayName,
    });

    redirectNextLoginFlowPage();
  }

  return (
    <LoginFlowCard>
      <CardHeader>
        <CardTitle>Create new organization</CardTitle>
      </CardHeader>
      <CardContent>
        <Form {...form}>
          <form className="mt-2" onSubmit={form.handleSubmit(handleSubmit)}>
            <FormField
              control={form.control}
              name="displayName"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Organization Name</FormLabel>
                  <FormControl>
                    <Input placeholder="Example Corporation" {...field} />
                  </FormControl>
                  <FormDescription>
                    Usually, you want to name this after the company this work
                    is for. You can change this later.
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />

            <Button
              type="submit"
              className="mt-4 w-full"
              variant={darkMode ? "outline" : "default"}
              disabled={submitting}
            >
              {submitting && (
                <LoaderCircleIcon className="h-4 w-4 animate-spin" />
              )}
              Create Organization
            </Button>
          </form>
        </Form>
      </CardContent>
    </LoginFlowCard>
  );
}
