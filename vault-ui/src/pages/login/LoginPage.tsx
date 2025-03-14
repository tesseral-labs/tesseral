import { useMutation } from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import React, { useCallback, useState } from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";



import { OAuthButton, OAuthMethods } from "@/components/OAuthButton";
import { GoogleIcon } from "@/components/login/GoogleIcon";
import { LoginFlowCard } from "@/components/login/LoginFlowCard";
import { MicrosoftIcon } from "@/components/login/MicrosoftIcon";
import { Button } from "@/components/ui/button";
import { CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Form, FormControl, FormDescription, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  createIntermediateSession,
  issueEmailVerificationChallenge,
  setEmailAsPrimaryLoginFactor,
} from "@/gen/tesseral/intermediate/v1/intermediate-IntermediateService_connectquery";
import { useDarkMode } from "@/lib/dark-mode";
import { ProjectSettingsProvider, useProjectSettings } from "@/lib/project-settings";
import { LoaderCircleIcon } from "lucide-react";
import { useNavigate } from "react-router";





export function LoginPage() {
  return (
    <ProjectSettingsProvider>
      <LoginPageInner>
        <LoginPageContents />
      </LoginPageInner>
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
  const settings = useProjectSettings();
  const isDarkMode = useDarkMode();

  return (
    <div className="bg-body w-screen min-h-screen mx-auto flex flex-col justify-center items-center py-8">
      <div className="mx-auto max-w-7xl sm:px-6 lg:px-8 flex justify-center">
        <div className="mb-8">
          <object
            className="max-w-[180px]"
            data={isDarkMode ? settings?.darkModeLogoUrl : settings?.logoUrl}
          />
        </div>
      </div>

      <div className="max-w-sm w-full mx-auto">{children}</div>
    </div>
  );
}

function SideBySideLoginPage({ children }: { children?: React.ReactNode }) {
  const settings = useProjectSettings();
  const isDarkMode = useDarkMode();

  return (
    <div className="bg-body w-screen min-h-screen grid grid-cols-2 gap-0">
      <div className="bg-primary" />
      <div className="flex flex-col justify-center items-center p-4">
        <div className="mx-auto max-w-7xl sm:px-6 lg:px-8 flex justify-center">
          <div className="mb-4">
            <object
              className="max-w-[180px]"
              data={isDarkMode ? settings?.darkModeLogoUrl : settings?.logoUrl}
            />
          </div>
        </div>

        {children}
      </div>
    </div>
  );
}

const schema = z.object({
  email: z.string().email(),
});

function LoginPageContents() {
  const settings = useProjectSettings();
  const darkMode = useDarkMode();

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      email: "",
    },
  });

  const [submitting, setSubmitting] = useState(false);
  const createIntermediateSessionMutation = useMutation(createIntermediateSession)
  const setEmailAsPrimaryLoginFactorMutation = useMutation(setEmailAsPrimaryLoginFactor)
  const issueEmailVerificationChallengeMutation = useMutation(issueEmailVerificationChallenge)
  const navigate = useNavigate()
  async function handleSubmit(values: z.infer<typeof schema>) {
    setSubmitting(true)
    await createIntermediateSessionMutation.mutateAsync({});
    await setEmailAsPrimaryLoginFactorMutation.mutateAsync({});
    await issueEmailVerificationChallengeMutation.mutateAsync({
      email: values.email,
    });

    navigate("/verify-email");
  }

  return (
    <LoginFlowCard>
      <CardHeader>
        <CardTitle>Log in to {settings.projectDisplayName}</CardTitle>
        <CardDescription>Please sign in to continue.</CardDescription>
      </CardHeader>
      <CardContent>
        <div className="space-y-2">
          {settings.logInWithGoogle && (
            <Button
              className="w-full"
              variant={darkMode ? "default" : "outline"}
            >
              <GoogleIcon />
              Log in with Google
            </Button>
          )}
          {settings.logInWithMicrosoft && (
            <Button
              className="w-full"
              variant={darkMode ? "default" : "outline"}
            >
              <MicrosoftIcon />
              Log in with Microsoft
            </Button>
          )}
        </div>

        {(settings.logInWithEmail || settings.logInWithSaml) && (
          <>
            <div className="block relative w-full cursor-default my-2 mt-6">
              <div className="absolute inset-0 flex items-center border-muted-foreground">
                <span className="w-full border-t"></span>
              </div>
              <div className="relative flex justify-center text-xs uppercase">
                <span className="bg-card px-2 text-muted-foreground">or</span>
              </div>
            </div>

            <Form {...form}>
              <form onSubmit={form.handleSubmit(handleSubmit)}>
                <FormField
                  control={form.control}
                  name="email"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Email</FormLabel>
                      <FormControl>
                        <Input {...field} />
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
