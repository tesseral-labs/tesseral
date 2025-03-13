import { useMutation } from "@connectrpc/connect-query";
import debounce from "lodash.debounce";
import React, {
  ChangeEvent,
  Dispatch,
  FormEvent,
  SetStateAction,
  useCallback,
  useEffect,
  useState,
} from "react";
import { toast } from "sonner";

import { setIntermediateSessionToken } from "@/auth";
import {
  createIntermediateSession,
  issueEmailVerificationChallenge,
  listSAMLOrganizations,
  setEmailAsPrimaryLoginFactor,
} from "@/gen/tesseral/intermediate/v1/intermediate-IntermediateService_connectquery";
import { Organization } from "@/gen/tesseral/intermediate/v1/intermediate_pb";
import { AuthType, useAuthType } from "@/lib/auth";
import { parseErrorMessage } from "@/lib/errors";
import { LoginViews } from "@/lib/views";

import { Button } from "./ui/button";
import { Input } from "./ui/input";
import { Label } from "./ui/label";
import Loader from "./ui/loader";
import TextDivider from "./ui/text-divider";

interface EmailFormProps {
  disableLogInWithEmail?: boolean;
  skipIntermediateSessionCreation?: boolean;
  skipListSAMLOrganizations?: boolean;
  setView: Dispatch<SetStateAction<LoginViews>>;
}

export function EmailForm({
  disableLogInWithEmail = false,
  skipListSAMLOrganizations = false,
  skipIntermediateSessionCreation = false,
  setView,
}: EmailFormProps) {
  const emailRegex = /^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,6}$/i;

  const authType = useAuthType();

  const createIntermediateSessionMutation = useMutation(
    createIntermediateSession,
  );
  const issueEmailVerificationChallengeMutation = useMutation(
    issueEmailVerificationChallenge,
  );
  const listSAMLOrganizationsMutation = useMutation(listSAMLOrganizations);
  const setEmailAsPrimaryLoginFactorMutation = useMutation(
    setEmailAsPrimaryLoginFactor,
  );

  const [email, setEmail] = useState<string>("");
  const [emailIsValid, setEmailIsValid] = useState<boolean>(false);
  const [samlOrganizations, setSamlOrganizations] = useState<Organization[]>(
    [],
  );
  const [submitting, setSubmitting] = useState<boolean>(false);

  // eslint-disable-next-line react-hooks/exhaustive-deps
  const fetchSamlOrganizations = useCallback(
    debounce(async () => {
      const { organizations } = await listSAMLOrganizationsMutation.mutateAsync(
        {
          email,
        },
      );

      setSamlOrganizations(organizations);
    }, 300),
    [email],
  );

  function handleEmail(e: ChangeEvent<HTMLInputElement>) {
    setEmail(e.target.value);
  }

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    if (!disableLogInWithEmail) {
      setSubmitting(true);

      try {
        if (!skipIntermediateSessionCreation) {
          // this sets a cookie that subsequent requests use
          const { intermediateSessionSecretToken } =
            await createIntermediateSessionMutation.mutateAsync({});

          // set the intermediate sessionToken
          setIntermediateSessionToken(intermediateSessionSecretToken);
        }

        await setEmailAsPrimaryLoginFactorMutation.mutateAsync({});

        await issueEmailVerificationChallengeMutation.mutateAsync({
          email,
        });

        setSubmitting(false);

        // redirect to challenge page
        setView(LoginViews.VerifyEmail);
      } catch (error) {
        setSubmitting(false);
        const message = parseErrorMessage(error);
        toast.error("Could not initiate login", {
          description: message,
        });
      }
    }
  }

  useEffect(() => {
    void (async () => {
      const valid = emailRegex.test(email);
      setEmailIsValid(valid);

      if (valid && !skipListSAMLOrganizations) {
        await fetchSamlOrganizations();
      }
    })();
  }, [email]); // eslint-disable-line react-hooks/exhaustive-deps

  return (
    <>
      <form
        className="flex flex-col justify-center w-full"
        onSubmit={handleSubmit}
      >
        <div className="grid gap-2">
          <Label htmlFor="email">Email</Label>
          <Input
            className="w-full mb-2"
            id="email"
            type="email"
            onChange={handleEmail}
            placeholder="jane.doe@email.com"
            value={email}
          />
        </div>

        {!disableLogInWithEmail && (
          <Button type="submit" disabled={!emailIsValid || submitting}>
            {submitting && <Loader />}
            {authType === AuthType.SignUp ? "Sign up" : "Log in"}
          </Button>
        )}
      </form>

      {samlOrganizations && samlOrganizations.length > 0 && (
        <>
          <TextDivider>or continue with SAML</TextDivider>

          {samlOrganizations.map((organization) => (
            <div key={organization.id} className="flex flex-col items-center">
              <label
                className="text-center uppercase text-foreground font-semibold text-sm mb-3 tracking-wide w-full"
                htmlFor="email"
              >
                Continue with SAML
              </label>
              <a
                href={`/api/saml/v1/${organization.primarySamlConnectionId}/init`}
                className="w-[clamp(240px,50%,100%)]"
              >
                <Button variant="outline">{organization.displayName}</Button>
              </a>
            </div>
          ))}
        </>
      )}
    </>
  );
}
