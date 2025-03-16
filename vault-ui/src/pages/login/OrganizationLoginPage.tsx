import { useMutation } from "@connectrpc/connect-query";
import { LoaderCircleIcon } from "lucide-react";
import React, { useEffect } from "react";
import { useParams } from "react-router";

import { setOrganization } from "@/gen/tesseral/intermediate/v1/intermediate-IntermediateService_connectquery";
import { useRedirectNextLoginFlowPage } from "@/hooks/use-redirect-next-login-flow-page";

export function OrganizationLoginPage() {
  const { organizationId } = useParams();
  const { mutateAsync: setOrganizationMutationAsync } =
    useMutation(setOrganization);
  const redirectNextLoginFlowPage = useRedirectNextLoginFlowPage();

  useEffect(() => {
    (async () => {
      await setOrganizationMutationAsync({
        organizationId,
      });

      redirectNextLoginFlowPage();
    })();
  }, [organizationId, redirectNextLoginFlowPage, setOrganizationMutationAsync]);

  return (
    <LoaderCircleIcon className="mx-auto text-muted-foreground h-4 w-4 animate-spin" />
  );
}
