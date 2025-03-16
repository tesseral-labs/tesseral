import { useMutation } from "@connectrpc/connect-query";
import { LoaderCircleIcon } from "lucide-react";
import React, { useEffect } from "react";
import { useParams } from "react-router";

import {
  exchangeSessionForIntermediateSession,
  setOrganization,
} from "@/gen/tesseral/intermediate/v1/intermediate-IntermediateService_connectquery";
import { useRedirectNextLoginFlowPage } from "@/hooks/use-redirect-next-login-flow-page";
import { clearAccessToken } from "@/pages/dashboard/LoggedInGate";

export function SwitchOrganizationsPage() {
  const { organizationId } = useParams();
  const { mutateAsync: exchangeSessionForIntermediateSessionAsync } =
    useMutation(exchangeSessionForIntermediateSession);
  const { mutateAsync: setOrganizationAsync } = useMutation(setOrganization);
  const redirectNextLoginFlowPage = useRedirectNextLoginFlowPage();

  useEffect(() => {
    (async () => {
      // clear our cached access token
      clearAccessToken();

      await exchangeSessionForIntermediateSessionAsync({});
      await setOrganizationAsync({ organizationId });

      redirectNextLoginFlowPage();
    })();
  }, [
    exchangeSessionForIntermediateSessionAsync,
    organizationId,
    redirectNextLoginFlowPage,
    setOrganizationAsync,
  ]);

  return (
    <LoaderCircleIcon className="mx-auto text-muted-foreground h-4 w-4 animate-spin" />
  );
}
