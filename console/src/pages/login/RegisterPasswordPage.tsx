import { Code, ConnectError } from '@connectrpc/connect';
import { useMutation } from '@connectrpc/connect-query';
import { zodResolver } from '@hookform/resolvers/zod';
import { LoaderCircleIcon } from 'lucide-react';
import React, { useState } from 'react';
import { useForm } from 'react-hook-form';
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
import { Input } from '@/components/ui/input';
import { registerPassword } from '@/gen/tesseral/intermediate/v1/intermediate-IntermediateService_connectquery';
import { useRedirectNextLoginFlowPage } from '@/hooks/use-redirect-next-login-flow-page';
import { Title } from '@/components/Title';

const schema = z.object({
  password: z.string().nonempty(),
});

export function RegisterPasswordPage() {
  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      password: '',
    },
  });

  const [submitting, setSubmitting] = useState(false);
  const { mutateAsync: registerPasswordAsync } = useMutation(registerPassword);
  const redirectNextLoginFlowPage = useRedirectNextLoginFlowPage();

  async function handleSubmit(values: z.infer<typeof schema>) {
    setSubmitting(true);

    try {
      await registerPasswordAsync({
        password: values.password,
      });
    } catch (e) {
      if (
        e instanceof ConnectError &&
        e.code === Code.FailedPrecondition &&
        e.rawMessage === 'password_compromised'
      ) {
        form.setError('password', {
          type: 'manual',
          message:
            'This password has been reported as compromised. Please choose a different password.',
        });
        return;
      }

      throw e;
    } finally {
      setSubmitting(false);
    }

    redirectNextLoginFlowPage();
  }

  return (
    <LoginFlowCard>
      <Title title="Register Password" />
      <CardHeader>
        <CardTitle>Register password</CardTitle>
        <CardDescription>
          Register a password to continue logging in.
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
                  <FormDescription>
                    Choose a unique password you haven't used for anything else.
                  </FormDescription>
                </FormItem>
              )}
            />

            <Button type="submit" className="mt-4 w-full" disabled={submitting}>
              {submitting && (
                <LoaderCircleIcon className="h-4 w-4 animate-spin" />
              )}
              Register Password
            </Button>
          </form>
        </Form>
      </CardContent>
    </LoginFlowCard>
  );
}
