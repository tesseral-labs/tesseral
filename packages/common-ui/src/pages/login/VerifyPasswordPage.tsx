import { useMutation } from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import { LoaderCircleIcon } from "lucide-react";
import React, { useState } from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { LoginFlowCard } from "../../components/login/LoginFlowCard";
import { Button } from "../../components/ui/button";
import {
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "../../components/ui/card";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "../../components/ui/form";
import { Input } from "../../components/ui/input";
import { verifyPassword } from "../../gen/tesseral/intermediate/v1/intermediate-IntermediateService_connectquery";
import { useRedirectNextLoginFlowPage } from "../../hooks/use-redirect-next-login-flow-page";
import { useDarkMode } from "../../lib/dark-mode";

const schema = z.object({
  password: z.string().nonempty(),
});

export function VerifyPasswordPage() {
  const darkMode = useDarkMode();
  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      password: "",
    },
  });

  const [submitting, setSubmitting] = useState(false);
  const { mutateAsync: verifyPasswordAsync } = useMutation(verifyPassword);
  const redirectNextLoginFlowPage = useRedirectNextLoginFlowPage();

  async function handleSubmit(values: z.infer<typeof schema>) {
    setSubmitting(true);

    await verifyPasswordAsync({
      password: values.password,
    });

    redirectNextLoginFlowPage();
  }

  return (
    <LoginFlowCard>
      <CardHeader>
        <CardTitle>Verify your password</CardTitle>
        <CardDescription>
          Enter your password below to continue logging in.
        </CardDescription>
      </CardHeader>
      <CardContent>
        <Form {...form}>
          <form className="mt-2" onSubmit={form.handleSubmit(handleSubmit)}>
            <FormField
              control={form.control}
              name="password"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Password</FormLabel>
                  <FormControl>
                    <Input type="password" {...field} />
                  </FormControl>
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
              Verify Password
            </Button>
          </form>
        </Form>
      </CardContent>
    </LoginFlowCard>
  );
}
