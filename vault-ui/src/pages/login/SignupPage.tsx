import { useMutation, useQuery } from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import { LoaderCircleIcon } from "lucide-react";
import React, { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { useNavigate } from "react-router";
import { Link } from "react-router";
import { z } from "zod";

import { Title } from "@/components/core/Title";
import { UISettingsInjector } from "@/components/core/UISettingsInjector";
import { LoginFlowCard } from "@/components/login/LoginFlowCard";
import { GithubIcon } from "@/components/login/icons/GithubIcon";
import { GoogleIcon } from "@/components/login/icons/GoogleIcon";
import { MicrosoftIcon } from "@/components/login/icons/MicrosoftIcon";
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
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import {
  createIntermediateSession,
  getGithubOAuthRedirectURL,
  getGoogleOAuthRedirectURL,
  getMicrosoftOAuthRedirectURL,
  issueEmailVerificationChallenge,
  listOIDCOrganizations,
  listSAMLOrganizations,
  setEmailAsPrimaryLoginFactor,
  setPasswordAsPrimaryLoginFactor,
} from "@/gen/tesseral/intermediate/v1/intermediate-IntermediateService_connectquery";
import { useLoginPageQueryParams } from "@/hooks/use-login-page-query-params";
import {
  ProjectSettingsProvider,
  useProjectSettings,
} from "@/lib/project-settings";

export function SignupPage() {
  return (
    <ProjectSettingsProvider>
      <UISettingsInjector>
        <SignupPageInner>
          <SignupPageContents />
        </SignupPageInner>
      </UISettingsInjector>
    </ProjectSettingsProvider>
  );
}

function SignupPageInner({ children }: { children?: React.ReactNode }) {
  const { logInLayout } = useProjectSettings();

  return (
    <>
      {logInLayout === "centered" ? (
        <CenteredSignupPage>{children}</CenteredSignupPage>
      ) : (
        <SideBySideSignupPage>{children}</SideBySideSignupPage>
      )}
    </>
  );
}

function CenteredSignupPage({ children }: { children?: React.ReactNode }) {
  return (
    <div className="bg-background w-full min-h-screen mx-auto flex flex-col justify-center items-center py-8">
      <div className="max-w-sm w-full mx-auto">{children}</div>
    </div>
  );
}

function SideBySideSignupPage({ children }: { children?: React.ReactNode }) {
  return (
    <div className="bg-background w-full min-h-screen grid grid-cols-1 md:grid-cols-2 gap-0">
      <div className="bg-primary hidden md:block" />
      <div className="flex flex-col justify-center items-center p-4">
        <div className="max-w-sm w-full mx-auto">{children}</div>
      </div>
    </div>
  );
}

const schema = z.object({
  email: z.string().email(),
});

function SignupPageContents() {
  const settings = useProjectSettings();

  const createIntermediateSessionMutation = useMutation(
    createIntermediateSession,
  );

  const [
    { relayedSessionState, redirectURI, returnRelayedSessionTokenAsQueryParam },
    serializedQueryParamState,
  ] = useLoginPageQueryParams();

  async function createIntermediateSessionWithRelayedSessionState() {
    await createIntermediateSessionMutation.mutateAsync({
      relayedSessionState,
      redirectUri: redirectURI,
      returnRelayedSessionTokenAsQueryParam,
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

  const { mutateAsync: getGithubOAuthRedirectURLAsync } = useMutation(
    getGithubOAuthRedirectURL,
  );

  async function handleLogInWithGithub() {
    await createIntermediateSessionWithRelayedSessionState();
    const { url } = await getGithubOAuthRedirectURLAsync({
      redirectUrl: `${window.location.origin}/github-oauth-callback`,
    });
    window.location.href = url;
  }

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      email: "",
    },
  });

  const [submitting, setSubmitting] = useState(false);
  const setEmailAsPrimaryLoginFactorMutation = useMutation(
    setEmailAsPrimaryLoginFactor,
  );
  const issueEmailVerificationChallengeMutation = useMutation(
    issueEmailVerificationChallenge,
  );
  const { mutateAsync: setPasswordAsPrimaryLoginFactorAsync } = useMutation(
    setPasswordAsPrimaryLoginFactor,
  );

  const navigate = useNavigate();

  async function handleSubmit(values: z.infer<typeof schema>) {
    setSubmitting(true);

    await createIntermediateSessionWithRelayedSessionState();

    // prefer to sign up with "Log in with Email"
    if (settings.logInWithEmail) {
      await setEmailAsPrimaryLoginFactorMutation.mutateAsync({});
    } else if (settings.logInWithPassword) {
      await setPasswordAsPrimaryLoginFactorAsync({});
    }

    await issueEmailVerificationChallengeMutation.mutateAsync({
      email: values.email,
    });

    navigate("/verify-email");
  }

  const watchEmail = form.watch("email");
  const [debouncedEmail, setDebouncedEmail] = useState("");
  useEffect(() => {
    const interval = setInterval(() => setDebouncedEmail(watchEmail), 250);
    return () => clearInterval(interval);
  }, [watchEmail]);

  const { data: listSAMLOrganizationsResponse } = useQuery(
    listSAMLOrganizations,
    {
      email: debouncedEmail,
    },
    {
      enabled: settings.logInWithSaml && debouncedEmail.includes("@"),
    },
  );

  const { data: listOIDCOrganizationsResponse } = useQuery(
    listOIDCOrganizations,
    {
      email: debouncedEmail,
    },
    {
      enabled: settings.logInWithOidc && debouncedEmail.includes("@"),
    },
  );

  const hasAboveFoldMethod =
    settings.logInWithGoogle ||
    settings.logInWithMicrosoft ||
    settings.logInWithGithub;
  const hasBelowFoldMethod =
    settings.logInWithEmail ||
    settings.logInWithPassword ||
    settings.logInWithSaml ||
    settings.logInWithOidc;

  return (
    <LoginFlowCard>
      <Title title="Sign up" />
      <CardHeader>
        <CardTitle>Sign up for {settings.projectDisplayName}</CardTitle>
        <CardDescription>Please sign up to continue.</CardDescription>
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
              Sign up with Google
            </Button>
          )}
          {settings.logInWithMicrosoft && (
            <Button
              className="w-full"
              variant="outline"
              onClick={handleLogInWithMicrosoft}
            >
              <MicrosoftIcon />
              Sign up with Microsoft
            </Button>
          )}
          {settings.logInWithGithub && (
            <Button
              className="w-full"
              variant="outline"
              onClick={handleLogInWithGithub}
            >
              <GithubIcon />
              Sign up with GitHub
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
                Sign up
              </Button>

              {listSAMLOrganizationsResponse?.organizations?.map((org) => (
                <a
                  key={org.id}
                  href={`/api/saml/v1/${org.primarySamlConnectionId}/init`}
                >
                  <Button type="button" className="mt-4 w-full">
                    Log in with SAML ({org.displayName})
                  </Button>
                </a>
              ))}

              {listOIDCOrganizationsResponse?.organizations?.map((org) => (
                <a
                  key={org.id}
                  href={`/api/oidc/v1/${org.primaryOidcConnectionId}/init`}
                >
                  <Button type="button" className="mt-4 w-full">
                    Log in with OIDC ({org.displayName})
                  </Button>
                </a>
              ))}
            </form>
          </Form>
        )}

        <p className="mt-4 text-xs text-muted-foreground">
          Already have an account?{" "}
          <Link
            to={`/login${serializedQueryParamState}`}
            className="cursor-pointer text-foreground underline underline-offset-2 decoration-muted-foreground"
          >
            Log in.
          </Link>
        </p>
      </CardContent>
    </LoginFlowCard>
  );
}
