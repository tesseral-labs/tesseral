import { useMutation } from '@connectrpc/connect-query';
import { zodResolver } from '@hookform/resolvers/zod';
import { LoaderCircleIcon } from 'lucide-react';
import React, { useEffect, useState } from 'react';
import { useForm } from 'react-hook-form';
import { useNavigate } from 'react-router';
import { Link, useSearchParams } from 'react-router-dom';
import { z } from 'zod';

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
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form';
import { Input } from '@/components/ui/input';
import {
  createIntermediateSession,
  getGoogleOAuthRedirectURL,
  getMicrosoftOAuthRedirectURL,
  issueEmailVerificationChallenge,
  setEmailAsPrimaryLoginFactor,
} from '@/gen/tesseral/intermediate/v1/intermediate-IntermediateService_connectquery';
import {
  ProjectSettingsProvider,
  useProjectSettings,
} from '@/lib/project-settings';
import { Title } from '@/components/Title';

export function LoginPage() {
  return (
    <ProjectSettingsProvider>
      <CenteredLoginPage>
        <LoginPageContents />
      </CenteredLoginPage>
    </ProjectSettingsProvider>
  );
}

function CenteredLoginPage({ children }: { children?: React.ReactNode }) {
  return (
    <div className="bg-zinc-950 w-screen min-h-screen mx-auto flex flex-col justify-center items-center py-8 relative">
      <div className="absolute flex justify-center items-center blur-3xl w-full z-5">
        <div className="relative rounded-full w-[750px] h-[750px] bg-indigo-600/30 blur-3xl m-auto" />
      </div>

      <div className="mx-auto max-w-7xl sm:px-6 lg:px-8 flex justify-center z-10">
        <div className="mb-8">
          <img
            className="max-w-[180px]"
            src="/images/tesseral-logo-white.svg"
            alt="Tesseral"
          />
        </div>
      </div>

      <div className="max-w-sm w-full mx-auto z-10">{children}</div>
    </div>
  );
}

const schema = z.object({
  email: z.string().email(),
});

function LoginPageContents() {
  const settings = useProjectSettings();

  const createIntermediateSessionMutation = useMutation(
    createIntermediateSession,
  );
  const [searchParams, setSearchParams] = useSearchParams();

  const [relayedSessionState, setRelayedSessionState] = useState<
    string | undefined
  >();
  useEffect(() => {
    if (relayedSessionState !== undefined) {
      return;
    }

    setRelayedSessionState(
      searchParams.get('relayed-session-state') ?? undefined,
    );

    const searchParamsCopy = new URLSearchParams(searchParams);
    searchParamsCopy.delete('relayed-session-state');
    setSearchParams(searchParamsCopy);
  }, [relayedSessionState, searchParams, setSearchParams]);

  async function createIntermediateSessionWithRelayedSessionState() {
    await createIntermediateSessionMutation.mutateAsync({
      relayedSessionState,
    });
  }

  const { mutateAsync: getGoogleOAuthRedirectURLAsync } = useMutation(
    getGoogleOAuthRedirectURL,
  );

  async function handleLogInWithGoogle() {
    await createIntermediateSessionWithRelayedSessionState();
    const { url } = await getGoogleOAuthRedirectURLAsync({
      redirectUrl: `${window.location.origin}/google-oauth-callback`,
    });
    window.location.href = url;
  }

  const { mutateAsync: getMicrosoftOAuthRedirectURLAsync } = useMutation(
    getMicrosoftOAuthRedirectURL,
  );

  async function handleLogInWithMicrosoft() {
    await createIntermediateSessionWithRelayedSessionState();
    const { url } = await getMicrosoftOAuthRedirectURLAsync({
      redirectUrl: `${window.location.origin}/microsoft-oauth-callback`,
    });
    window.location.href = url;
  }

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      email: '',
    },
  });

  const [submitting, setSubmitting] = useState(false);
  const setEmailAsPrimaryLoginFactorMutation = useMutation(
    setEmailAsPrimaryLoginFactor,
  );
  const issueEmailVerificationChallengeMutation = useMutation(
    issueEmailVerificationChallenge,
  );
  const navigate = useNavigate();

  async function handleSubmit(values: z.infer<typeof schema>) {
    setSubmitting(true);
    await createIntermediateSessionWithRelayedSessionState();
    await setEmailAsPrimaryLoginFactorMutation.mutateAsync({});
    await issueEmailVerificationChallengeMutation.mutateAsync({
      email: values.email,
    });

    navigate('/verify-email');
  }

  const hasAboveFoldMethod =
    settings.logInWithGoogle || settings.logInWithMicrosoft;
  const hasBelowFoldMethod = settings.logInWithEmail || settings.logInWithSaml;

  return (
    <LoginFlowCard>
      <Title title="Log in" />
      <CardHeader>
        <CardTitle>Log in to Tesseral</CardTitle>
        <CardDescription>Please sign in to continue.</CardDescription>
      </CardHeader>
      <CardContent>
        <div className="space-y-2">
          {settings.logInWithGoogle && (
            <Button
              className="w-full"
              variant="outline"
              onClick={handleLogInWithGoogle}
            >
              <GoogleIcon />
              Log in with Google
            </Button>
          )}
          {settings.logInWithMicrosoft && (
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

        {hasAboveFoldMethod && hasBelowFoldMethod && (
          <div className="block relative w-full cursor-default my-2 mt-6">
            <div className="absolute inset-0 flex items-center border-muted-foreground">
              <span className="w-full border-t"></span>
            </div>
            <div className="relative flex justify-center text-xs uppercase">
              <span className="bg-card px-2 text-muted-foreground">or</span>
            </div>
          </div>
        )}

        {hasBelowFoldMethod && (
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
        )}

        <p className="mt-4 text-xs text-muted-foreground">
          Don't have an account?{' '}
          <Link
            to="/signup"
            className="cursor-pointer text-foreground underline underline-offset-2 decoration-muted-foreground"
          >
            Sign up.
          </Link>
        </p>
      </CardContent>
    </LoginFlowCard>
  );
}
