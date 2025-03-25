import { useMutation } from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import React from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { LoginFlowCard } from "@/components/login/LoginFlowCard";
import { Button } from "@/components/ui/button";
import {
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
import { verifyAuthenticatorApp } from "@/gen/tesseral/intermediate/v1/intermediate-IntermediateService_connectquery";
import { useRedirectNextLoginFlowPage } from "@/hooks/use-redirect-next-login-flow-page";

const schema = z.object({
  recoveryCode: z
    .string()
    .regex(/[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}/, {
      message: "Invalid recovery code",
    }),
});

export function VerifyAuthenticatorAppRecoveryCodePage() {
  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      recoveryCode: "",
    },
  });

  const { mutateAsync: verifyAuthenticatorAppAsync } = useMutation(
    verifyAuthenticatorApp,
  );
  const redirectNextLoginFlowPage = useRedirectNextLoginFlowPage();

  async function handleSubmit(values: z.infer<typeof schema>) {
    await verifyAuthenticatorAppAsync({
      recoveryCode: values.recoveryCode,
    });

    redirectNextLoginFlowPage();
  }

  return (
    <LoginFlowCard>
      <CardHeader>
        <CardTitle>Verify authenticator app recovery code</CardTitle>
        <CardDescription>
          To continue logging in, input one of the recovery codes for your
          authenticator app.
        </CardDescription>
      </CardHeader>
      <CardContent>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(handleSubmit)}>
            <FormField
              control={form.control}
              name="recoveryCode"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Recovery Code</FormLabel>
                  <FormControl>
                    <Input placeholder="0123-4567-89ab-cdef" {...field} />
                  </FormControl>
                  <FormDescription>
                    When you registered an authenticator app, you received a
                    list of recovery codes. Input one of those recovery codes
                    here.
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />

            <Button type="submit" className="mt-4 w-full">
              Verify authenticator app recovery code
            </Button>
          </form>
        </Form>
      </CardContent>
    </LoginFlowCard>
  );
}
