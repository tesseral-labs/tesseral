import { ConnectError } from "@connectrpc/connect";
import { useMutation, useQuery } from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import { LoaderCircleIcon } from "lucide-react";
import React, { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { useNavigate } from "react-router";
import { Link } from "react-router";
import { toast } from "sonner";
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
  verifyPassword,
} from "@/gen/tesseral/intermediate/v1/intermediate-IntermediateService_connectquery";
import { useLoginPageQueryParams } from "@/hooks/use-login-page-query-params";
import { useRedirectNextLoginFlowPage } from "@/hooks/use-redirect-next-login-flow-page";
import {
  ProjectSettingsProvider,
  useProjectSettings,
} from "@/lib/project-settings";

export function LoginPage() {
  return (
    <ProjectSettingsProvider>
      <UISettingsInjector>
        <LoginPageInner>
          <LoginPageContents />
        </LoginPageInner>
      </UISettingsInjector>
    </ProjectSettingsProvider>
  );
}

function LoginPageInner({ children }: { children?: React.ReactNode }) {
  const { logInLayout } = useProjectSettings();

  return (
    <>
      {logInLayout === "centered" ? (
        <CenteredLoginPage>{children}</CenteredLoginPage>
      ) : (
        <SideBySideLoginPage>{children}</SideBySideLoginPage>
      )}
    </>
  );
}

function CenteredLoginPage({ children }: { children?: React.ReactNode }) {
  return (
    <div className="bg-background w-full min-h-screen mx-auto flex flex-col justify-center items-center py-8">
      <div className="max-w-sm w-full mx-auto">{children}</div>
    </div>
  );
}

function SideBySideLoginPage({ children }: { children?: React.ReactNode }) {
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
  password: z.string(),
});

function LoginPageContents() {
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

  async function handleLogInWithSaml(samlConnectionId: string) {
    await createIntermediateSessionWithRelayedSessionState();
    window.location.href = `/api/saml/v1/${samlConnectionId}/init`;
  }

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      email: "",
      password: "",
    },
  });

  const [submitting, setSubmitting] = useState(false);
  const setEmailAsPrimaryLoginFactorMutation = useMutation(
    setEmailAsPrimaryLoginFactor,
  );
  const issueEmailVerificationChallengeMutation = useMutation(
    issueEmailVerificationChallenge,
  );
  const { mutateAsync: verifyPasswordAsync } = useMutation(verifyPassword);
  const { mutateAsync: setPasswordAsPrimaryLoginFactorAsync } = useMutation(
    setPasswordAsPrimaryLoginFactor,
  );
  const redirectNextLoginFlowPage = useRedirectNextLoginFlowPage();
  const navigate = useNavigate();

  async function handleLogInWithEmail(values: z.infer<typeof schema>) {
    setSubmitting(true);
    await createIntermediateSessionWithRelayedSessionState();
    await setEmailAsPrimaryLoginFactorMutation.mutateAsync({});
    await issueEmailVerificationChallengeMutation.mutateAsync({
      email: values.email,
    });

    navigate("/verify-email");
  }

  async function handleLogInWithPassword(values: z.infer<typeof schema>) {
    setSubmitting(true);
    await createIntermediateSessionWithRelayedSessionState();

    try {
      await verifyPasswordAsync({
        email: values.email,
        password: values.password,
      });
    } catch (e) {
      if (
        e instanceof ConnectError &&
        e.message === "[failed_precondition] incorrect_password"
      ) {
        form.setError("password", {
          type: "manual",
          message: "Incorrect password",
        });

        setSubmitting(false);
        return;
      }

      if (
        e instanceof ConnectError &&
        e.message === "[failed_precondition] passwords_unavailable_for_email"
      ) {
        await setPasswordAsPrimaryLoginFactorAsync({});
        await issueEmailVerificationChallengeMutation.mutateAsync({
          email: values.email,
        });

        toast.warning("To continue, you must verify your email address.");

        navigate("/verify-email");
        return;
      }

      throw e;
    }

    redirectNextLoginFlowPage();
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
    <>
      <LoginFlowCard>
        <Title title="Log in" />
        <CardHeader>
          <CardTitle>Log in to {settings.projectDisplayName}</CardTitle>
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
            {settings.logInWithGithub && (
              <Button
                className="w-full"
                variant="outline"
                onClick={handleLogInWithGithub}
              >
                <GithubIcon />
                Log in with GitHub
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
              <form
                onSubmit={form.handleSubmit(
                  settings.logInWithPassword
                    ? handleLogInWithPassword
                    : handleLogInWithEmail,
                )}
              >
                <div className="space-y-4">
                  <FormField
                    control={form.control}
                    name="email"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Email</FormLabel>
                        <FormControl>
                          <Input
                            type="email"
                            placeholder="john.doe@example.com"
                            {...field}
                          />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  {settings.logInWithPassword && (
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
                  )}
                </div>

                {settings.logInWithPassword ? (
                  <>
                    <Button
                      type="submit"
                      className="mt-4 w-full"
                      disabled={submitting}
                    >
                      Log in with Password
                    </Button>

                    {settings.logInWithEmail && (
                      <p className="text-center mt-4 text-xs text-muted-foreground">
                        or{" "}
                        <span
                          onClick={form.handleSubmit(handleLogInWithEmail)}
                          className=" cursor-pointer text-foreground underline underline-offset-2 decoration-muted-foreground"
                        >
                          Log in with Magic Link.
                        </span>
                      </p>
                    )}
                  </>
                ) : (
                  <Button
                    type="submit"
                    className="mt-4 w-full"
                    disabled={submitting}
                  >
                    {submitting && (
                      <LoaderCircleIcon className="h-4 w-4 animate-spin" />
                    )}
                    Log in with Magic Link
                  </Button>
                )}

                {listSAMLOrganizationsResponse?.organizations?.map((org) => (
                  <Button
                    key={org.id}
                    type="button"
                    className="mt-4 w-full"
                    onClick={() =>
                      handleLogInWithSaml(org.primarySamlConnectionId)
                    }
                  >
                    Log in with SAML ({org.displayName})
                  </Button>
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
        </CardContent>
      </LoginFlowCard>

      {settings.selfServeCreateUsers && (
        <p className="text-center mt-4 text-xs text-muted-foreground">
          Don't have an account?{" "}
          <Link
            to={`/signup${serializedQueryParamState}`}
            className="cursor-pointer text-foreground underline underline-offset-2 decoration-muted-foreground"
          >
            Sign up.
          </Link>
        </p>
      )}
    </>
  );
}
