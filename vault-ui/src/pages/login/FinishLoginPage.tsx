import { useMutation } from "@connectrpc/connect-query";
import { LoaderCircleIcon } from "lucide-react";
import React, { useEffect } from "react";

import { exchangeIntermediateSessionForSession } from "@/gen/tesseral/intermediate/v1/intermediate-IntermediateService_connectquery";
import { useProjectSettings } from "@/lib/project-settings";

export function FinishLoginPage() {
  const settings = useProjectSettings();
  const { mutateAsync: exchangeIntermediateSessionForSessionAsync } =
    useMutation(exchangeIntermediateSessionForSession);

  useEffect(() => {
    (async () => {
      const {
        newUser,
        relayedSessionToken,
        redirectUri,
        returnRelayedSessionTokenAsQueryParam,
      } = await exchangeIntermediateSessionForSessionAsync({});

      const preferredRedirect =
        (newUser
          ? settings.afterSignupRedirectUri
          : settings.afterLoginRedirectUri) || settings.redirectUri;

      const url = new URL(redirectUri || preferredRedirect);

      if (relayedSessionToken) {
        if (returnRelayedSessionTokenAsQueryParam) {
          url.searchParams.set(
            `__tesseral_${settings.projectId}_relayed_session_token`,
            relayedSessionToken,
          );
          url.searchParams.set(
            `__tesseral_${settings.projectId}_redirect_uri`,
            preferredRedirect,
          );
        } else {
          url.hash = new URLSearchParams({
            [`__tesseral_${settings.projectId}_relayed_session_token`]:
              relayedSessionToken,
          }).toString();
        }
      }

      window.location.href = url.toString();
    })();
  }, [settings, exchangeIntermediateSessionForSessionAsync]);

  return (
    <LoaderCircleIcon className="mx-auto text-muted-foreground h-4 w-4 animate-spin" />
  );
}
