import { useMutation, useQuery } from "@connectrpc/connect-query";
import { LoaderCircleIcon } from "lucide-react";
import React, { useEffect } from "react";
import { useNavigate } from "react-router";
import { useSearchParams } from "react-router-dom";

import {
  issueEmailVerificationChallenge,
  redeemGithubOAuthCode,
  whoami,
} from "@/gen/tesseral/intermediate/v1/intermediate-IntermediateService_connectquery";
import { useRedirectNextLoginFlowPage } from "@/hooks/use-redirect-next-login-flow-page";

export function GithubOAuthCallbackPage() {
  const [searchParams] = useSearchParams();
  const code = searchParams.get("code");
  const state = searchParams.get("state");
  const { refetch: refetchWhoami } = useQuery(whoami);
  const navigate = useNavigate();

  const { mutateAsync: redeemGithubOAuthCodeAsync } = useMutation(
    redeemGithubOAuthCode,
  );

  const { mutateAsync: issueEmailVerificationChallengeMutationAsync } =
    useMutation(issueEmailVerificationChallenge);

  const redirectNextLoginFlowPage = useRedirectNextLoginFlowPage();

  useEffect(() => {
    (async () => {
      await redeemGithubOAuthCodeAsync({
        code: code!,
        state: state!,
        redirectUrl: `${window.location.origin}/github-oauth-callback`,
      });

      redirectNextLoginFlowPage();
    })();
  }, [
    code,
    issueEmailVerificationChallengeMutationAsync,
    navigate,
    redeemGithubOAuthCodeAsync,
    redirectNextLoginFlowPage,
    refetchWhoami,
    state,
  ]);

  return (
    <LoaderCircleIcon className="mx-auto text-muted-foreground h-4 w-4 animate-spin" />
  );
}
