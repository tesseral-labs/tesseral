import { useMutation, useQuery } from "@connectrpc/connect-query";
import { LoaderCircleIcon } from "lucide-react";
import React, { useEffect } from "react";
import { useNavigate } from "react-router";
import { useSearchParams } from "react-router-dom";

import {
  issueEmailVerificationChallenge,
  redeemMicrosoftOAuthCode,
  whoami,
} from "../../gen/tesseral/intermediate/v1/intermediate-IntermediateService_connectquery";
import { useRedirectNextLoginFlowPage } from "../../hooks/use-redirect-next-login-flow-page";

export function MicrosoftOAuthCallbackPage() {
  const [searchParams] = useSearchParams();
  const code = searchParams.get("code");
  const state = searchParams.get("state");
  const { refetch: refetchWhoami } = useQuery(whoami);
  const navigate = useNavigate();

  const { mutateAsync: redeemMicrosoftOAuthCodeAsync } = useMutation(
    redeemMicrosoftOAuthCode
  );

  const { mutateAsync: issueEmailVerificationChallengeMutationAsync } =
    useMutation(issueEmailVerificationChallenge);

  const redirectNextLoginFlowPage = useRedirectNextLoginFlowPage();

  useEffect(() => {
    (async () => {
      await redeemMicrosoftOAuthCodeAsync({
        code: code!,
        state: state!,
        redirectUrl: `${window.location.origin}/microsoft-oauth-callback`,
      });

      redirectNextLoginFlowPage();
    })();
  }, [
    code,
    issueEmailVerificationChallengeMutationAsync,
    navigate,
    redeemMicrosoftOAuthCodeAsync,
    redirectNextLoginFlowPage,
    refetchWhoami,
    state,
  ]);

  return (
    <LoaderCircleIcon className="mx-auto text-muted-foreground h-4 w-4 animate-spin" />
  );
}
