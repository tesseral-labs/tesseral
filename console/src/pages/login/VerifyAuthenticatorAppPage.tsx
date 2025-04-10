import { useMutation } from '@connectrpc/connect-query';
import { zodResolver } from '@hookform/resolvers/zod';
import { REGEXP_ONLY_DIGITS } from 'input-otp';
import React from 'react';
import { useForm } from 'react-hook-form';
import { Link } from 'react-router-dom';
import { z } from 'zod';

import { LoginFlowCard } from '@/components/login/LoginFlowCard';
import { Button } from '@/components/ui/button';
import {
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form';
import {
  InputOTP,
  InputOTPGroup,
  InputOTPSlot,
} from '@/components/ui/input-otp';
import { verifyAuthenticatorApp } from '@/gen/tesseral/intermediate/v1/intermediate-IntermediateService_connectquery';
import { useRedirectNextLoginFlowPage } from '@/hooks/use-redirect-next-login-flow-page';
import { Title } from '@/components/Title';

const schema = z.object({
  totpCode: z.string().length(6),
});

export function VerifyAuthenticatorAppPage() {
  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      totpCode: '',
    },
  });

  const { mutateAsync: verifyAuthenticatorAppAsync } = useMutation(
    verifyAuthenticatorApp,
  );
  const redirectNextLoginFlowPage = useRedirectNextLoginFlowPage();

  async function handleSubmit(values: z.infer<typeof schema>) {
    await verifyAuthenticatorAppAsync({
      totpCode: values.totpCode,
    });

    redirectNextLoginFlowPage();
  }

  return (
    <LoginFlowCard>
      <Title title="Verify authenticator app" />
      <CardHeader>
        <CardTitle>Verify authenticator app</CardTitle>
        <CardDescription>
          To continue logging in, input a one-time password from your
          authenticator app.
        </CardDescription>
      </CardHeader>
      <CardContent>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(handleSubmit)}>
            <FormField
              control={form.control}
              name="totpCode"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>One-Time Password</FormLabel>
                  <FormControl>
                    <InputOTP
                      pattern={REGEXP_ONLY_DIGITS}
                      maxLength={6}
                      {...field}
                    >
                      <InputOTPGroup>
                        <InputOTPSlot index={0} />
                        <InputOTPSlot index={1} />
                        <InputOTPSlot index={2} />
                        <InputOTPSlot index={3} />
                        <InputOTPSlot index={4} />
                        <InputOTPSlot index={5} />
                      </InputOTPGroup>
                    </InputOTP>
                  </FormControl>
                  <FormDescription>
                    Enter a six-digit code from your authenticator app.
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />

            <Button type="submit" className="mt-4 w-full">
              Verify authenticator app
            </Button>
          </form>
        </Form>

        <p className="mt-4 text-xs text-muted-foreground">
          Lost your authenticator app?{' '}
          <Link
            to="/verify-authenticator-app-recovery-code"
            className="text-foreground underline underline-offset-2 decoration-muted-foreground"
          >
            Use a recovery code instead.
          </Link>
        </p>
      </CardContent>
    </LoginFlowCard>
  );
}
