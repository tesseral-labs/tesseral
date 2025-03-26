import { useMutation, useQuery } from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import { LoaderCircleIcon } from "lucide-react";
import React, { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { useNavigate } from "react-router";
import { useSearchParams } from "react-router-dom";
import { toast } from "sonner";
import { z } from "zod";

import { LoginFlowCard } from "@/components/login/LoginFlowCard";
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion";
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
  issueEmailVerificationChallenge,
  verifyEmailChallenge,
  whoami,
} from "@/gen/tesseral/intermediate/v1/intermediate-IntermediateService_connectquery";

const schema = z.object({
  emailVerificationChallengeCode: z
    .string()
    .startsWith("email_verification_challenge_code_"),
});

export function VerifyEmailPage() {
  const { data: whoamiResponse } = useQuery(whoami);

  const issueEmailVerificationChallengeMutation = useMutation(
    issueEmailVerificationChallenge,
  );
  const [hasResent, setHasResent] = useState(false);

  async function handleResend() {
    await issueEmailVerificationChallengeMutation.mutateAsync({
      email: whoamiResponse?.intermediateSession?.email,
    });

    toast.success("New verification link sent");
    setHasResent(true);
  }

  useEffect(() => {
    // allow another send after 10 seconds
    setTimeout(() => {
      setHasResent(false);
    }, 10000);
  }, [hasResent]);

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      emailVerificationChallengeCode: "",
    },
  });

  const [submitting, setSubmitting] = useState(false);
  const { mutateAsync: verifyEmailChallengeAsync } =
    useMutation(verifyEmailChallenge);
  const navigate = useNavigate();

  async function handleSubmit(values: z.infer<typeof schema>) {
    setSubmitting(true);

    await verifyEmailChallengeAsync({
      code: values.emailVerificationChallengeCode,
    });

    navigate("/choose-organization");
  }

  const [searchParams] = useSearchParams();
  useEffect(() => {
    (async () => {
      const code = searchParams.get("code");
      if (code) {
        await verifyEmailChallengeAsync({
          code,
        });

        navigate("/choose-organization");
      }
    })();
  }, [searchParams, verifyEmailChallengeAsync, navigate]);

  return (
    <LoginFlowCard>
      <CardHeader>
        <CardTitle>Check your email</CardTitle>
        <CardDescription>
          We've sent an email verification link to{" "}
          <span className="font-medium">
            {whoamiResponse?.intermediateSession?.email}
          </span>
          .
        </CardDescription>
      </CardHeader>

      <CardContent>
        <p className="text-sm text-muted-foreground">
          Didn't receive an email?
        </p>

        <Button
          className="mt-4 w-full"
          variant="outline"
          disabled={hasResent}
          onClick={handleResend}
        >
          {hasResent
            ? "Email verification resent!"
            : "Resend verification link"}
        </Button>

        <div className="block relative w-full cursor-default my-2 mt-4">
          <div className="absolute inset-0 flex items-center border-muted-foreground">
            <span className="w-full border-t"></span>
          </div>
          <div className="relative flex justify-center text-xs uppercase">
            <span className="bg-card px-2 text-muted-foreground">or</span>
          </div>
        </div>

        <Accordion type="single" collapsible>
          <AccordionItem className="border-b-0" value="advanced">
            <AccordionTrigger className="text-sm">
              Enter verification code manually
            </AccordionTrigger>
            <AccordionContent>
              <Form {...form}>
                <form
                  className="mt-2"
                  onSubmit={form.handleSubmit(handleSubmit)}
                >
                  <FormField
                    control={form.control}
                    name="emailVerificationChallengeCode"
                    render={({ field }) => (
                      <FormItem className="px-1">
                        <FormLabel>Email Verification Challenge Code</FormLabel>
                        <FormControl>
                          <Input
                            placeholder="email_verification_challenge_code_..."
                            {...field}
                          />
                        </FormControl>
                        <FormDescription>
                          Paste the full verification code from the email you
                          received.
                        </FormDescription>
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
                    Verify Email
                  </Button>
                </form>
              </Form>
            </AccordionContent>
          </AccordionItem>
        </Accordion>
      </CardContent>
    </LoginFlowCard>
  );
}
