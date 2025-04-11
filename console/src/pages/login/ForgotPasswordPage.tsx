import { useMutation, useQuery } from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import React, { useEffect } from "react";
import { useForm } from "react-hook-form";
import { useNavigate } from "react-router";
import { toast } from "sonner";
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
import {
  issuePasswordResetCode,
  verifyPasswordResetCode,
  whoami,
} from "@/gen/tesseral/intermediate/v1/intermediate-IntermediateService_connectquery";
import { Title } from "@/components/Title";

const schema = z.object({
  passwordResetCode: z.string().startsWith("password_reset_code_"),
});

export function ForgotPasswordPage() {
  const { data: whoamiResponse } = useQuery(whoami);
  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      passwordResetCode: "",
    },
  });

  const { mutateAsync: issuePasswordResetCodeAsync } = useMutation(
    issuePasswordResetCode,
  );

  // issue an email on mount
  useEffect(() => {
    (async () => {
      await issuePasswordResetCodeAsync({});
    })();
  }, [issuePasswordResetCodeAsync]);

  async function handleResend() {
    await issuePasswordResetCodeAsync({});
    toast.success("New password reset code sent");
  }

  const { mutateAsync: verifyPasswordResetCodeAsync } = useMutation(
    verifyPasswordResetCode,
  );
  const navigate = useNavigate();

  async function handleSubmit(values: z.infer<typeof schema>) {
    await verifyPasswordResetCodeAsync({
      passwordResetCode: values.passwordResetCode,
    });

    navigate("/register-password");
  }

  return (
    <LoginFlowCard>
      <Title title="Forgot password" />
      <CardHeader>
        <CardTitle>Forgot password</CardTitle>
        <CardDescription>
          We've sent a password reset code to{" "}
          <span className="font-medium">
            {whoamiResponse?.intermediateSession?.email}
          </span>
          .
        </CardDescription>
      </CardHeader>

      <CardContent>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(handleSubmit)}>
            <FormField
              control={form.control}
              name="passwordResetCode"
              render={({ field }) => (
                <FormItem className="px-1">
                  <FormLabel>Password Reset Code</FormLabel>
                  <FormControl>
                    <Input placeholder="password_reset_code_..." {...field} />
                  </FormControl>
                  <FormDescription>
                    Paste the full code from the email you received.
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />

            <Button type="submit" className="mt-4 w-full">
              Confirm Password Reset Code
            </Button>
          </form>
        </Form>

        <p className="mt-4 text-xs text-muted-foreground">
          Didn't get an email?{" "}
          <span
            className="cursor-pointer text-foreground underline underline-offset-2 decoration-muted-foreground"
            onClick={handleResend}
          >
            Request another code.
          </span>
        </p>
      </CardContent>
    </LoginFlowCard>
  );
}
