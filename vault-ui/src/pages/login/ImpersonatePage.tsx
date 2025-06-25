import { useMutation } from "@connectrpc/connect-query";
import { LoaderCircleIcon } from "lucide-react";
import React, { useEffect } from "react";
import { useSearchParams } from "react-router";

import { redeemUserImpersonationToken } from "@/gen/tesseral/intermediate/v1/intermediate-IntermediateService_connectquery";
import { useProjectSettings } from "@/lib/project-settings";

export function ImpersonatePage() {
  const settings = useProjectSettings();
  const [searchParams] = useSearchParams();
  const secretUserImpersonationToken = searchParams.get(
    "secret-user-impersonation-token",
  );

  const { mutateAsync: redeemUserImpersonationTokenAsync } = useMutation(
    redeemUserImpersonationToken,
  );

  useEffect(() => {
    (async () => {
      await redeemUserImpersonationTokenAsync({
        secretUserImpersonationToken: secretUserImpersonationToken!,
      });

      window.location.href = settings.redirectUri;
    })();
  }, [
    redeemUserImpersonationTokenAsync,
    secretUserImpersonationToken,
    settings.redirectUri,
  ]);

  return (
    <LoaderCircleIcon className="mx-auto text-muted-foreground h-4 w-4 animate-spin" />
  );
}
