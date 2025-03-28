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
      const { newUser, relayedSessionToken } =
        await exchangeIntermediateSessionForSessionAsync({});

      const preferredRedirect = newUser
        ? settings.afterSignupRedirectUri
        : settings.afterLoginRedirectUri;

      const url = new URL(preferredRedirect ?? settings.redirectUri);

      if (relayedSessionToken) {
        url.hash = `#__tesseral_${settings.projectId}_relayed_session_token=${relayedSessionToken}`;
      }

      window.location.href = url.toString();
    })();
  }, [settings, exchangeIntermediateSessionForSessionAsync]);

  return (
    <LoaderCircleIcon className="mx-auto text-muted-foreground h-4 w-4 animate-spin" />
  );
}
