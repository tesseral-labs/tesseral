import { useMutation, useQuery } from '@connectrpc/connect-query';
import React, { useState } from 'react';
import { useForm } from 'react-hook-form';

import { GoogleIcon } from '@/components/login/GoogleIcon';
import { LoginFlowCard } from '@/components/login/LoginFlowCard';
import { MicrosoftIcon } from '@/components/login/MicrosoftIcon';
import { Button } from '@/components/ui/button';
import {
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import {
  getGoogleOAuthRedirectURL,
  getMicrosoftOAuthRedirectURL,
  issueEmailVerificationChallenge,
  listOrganizations,
  setEmailAsPrimaryLoginFactor,
  whoami,
} from '@/gen/tesseral/intermediate/v1/intermediate-IntermediateService_connectquery';
import { Title } from '@/components/Title';
import TextDivider from '@/components/ui/text-divider';
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form';
import { LoaderCircleIcon } from 'lucide-react';
import { z } from 'zod';
import { zodResolver } from '@hookform/resolvers/zod';
import { useNavigate } from 'react-router';
import { Input } from '@/components/ui/input';

export function AuthenticateAnotherWayPage() {
  const navigate = useNavigate();

  const [submitting, setSubmitting] = useState(false);

  const { data: whoamiResponse } = useQuery(whoami);
  const { data: listOrganizationsResponse } = useQuery(listOrganizations);

  const organization = listOrganizationsResponse?.organizations?.find(
    (org) => org.id === whoamiResponse?.intermediateSession?.organizationId,
  );

  const setEmailAsPrimaryLoginFactorMutation = useMutation(
    setEmailAsPrimaryLoginFactor,
  );
  const issueEmailVerificationChallengeMutation = useMutation(
    issueEmailVerificationChallenge,
  );

  const { mutateAsync: getGoogleOAuthRedirectURLAsync } = useMutation(
    getGoogleOAuthRedirectURL,
  );

  async function handleLogInWithGoogle() {
    const { url } = await getGoogleOAuthRedirectURLAsync({
      redirectUrl: `${window.location.origin}/google-oauth-callback`,
    });
    window.location.href = url;
  }

  const { mutateAsync: getMicrosoftOAuthRedirectURLAsync } = useMutation(
    getMicrosoftOAuthRedirectURL,
  );

  async function handleLogInWithMicrosoft() {
    const { url } = await getMicrosoftOAuthRedirectURLAsync({
      redirectUrl: `${window.location.origin}/microsoft-oauth-callback`,
    });
    window.location.href = url;
  }

  async function handleSubmit(values: z.infer<typeof schema>) {
    setSubmitting(true);
    await setEmailAsPrimaryLoginFactorMutation.mutateAsync({});
    await issueEmailVerificationChallengeMutation.mutateAsync({
      email: values.email,
    });

    navigate('/verify-email');
  }

  const schema = z.object({
    email: z.string().email(),
  });

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      email: '',
    },
  });

  return (
    <LoginFlowCard>
      <Title title="Authenticate another way" />
      <CardHeader>
        <CardTitle>Authenticate another way</CardTitle>
        <CardDescription>
          To continue logging in, you must choose from one of the login methods
          below.
        </CardDescription>
      </CardHeader>
      <CardContent>
        <div className="space-y-2">
          {organization?.logInWithGoogle && (
            <Button
              className="w-full"
              variant="outline"
              onClick={handleLogInWithGoogle}
            >
              <GoogleIcon />
              Log in with Google
            </Button>
          )}
          {organization?.logInWithMicrosoft && (
            <Button
              className="w-full"
              variant="outline"
              onClick={handleLogInWithMicrosoft}
            >
              <MicrosoftIcon />
              Log in with Microsoft
            </Button>
          )}
        </div>

        {organization?.logInWithEmail && (
          <>
            {(organization?.logInWithGoogle ||
              organization?.logInWithMicrosoft) && (
              <TextDivider>or</TextDivider>
            )}

            <Form {...form}>
              <form onSubmit={form.handleSubmit(handleSubmit)}>
                <FormField
                  control={form.control}
                  name="email"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Email</FormLabel>
                      <FormControl>
                        <Input placeholder="john.doe@example.com" {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <Button
                  type="submit"
                  className="mt-4 w-full"
                  disabled={submitting}
                >
                  {submitting && (
                    <LoaderCircleIcon className="h-4 w-4 animate-spin" />
                  )}
                  Log in
                </Button>
              </form>
            </Form>
          </>
        )}
      </CardContent>
    </LoginFlowCard>
  );
}
