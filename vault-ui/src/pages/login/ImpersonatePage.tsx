import { useMutation, useQuery } from '@connectrpc/connect-query';
import {
  getSettings,
  redeemUserImpersonationToken,
} from '@/gen/tesseral/intermediate/v1/intermediate-IntermediateService_connectquery';
import React, { useEffect } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { setAccessToken, setRefreshToken } from '@/auth';
import { useIntermediateExchangeAndRedirect } from '@/hooks/use-intermediate-exchange-and-redirect';

export const ImpersonatePage = () => {
  const redeemUserImpersonationTokenMutation = useMutation(
    redeemUserImpersonationToken,
  );
  const [searchParams] = useSearchParams();
  const secretUserImpersonationToken = searchParams.get(
    'secret-user-impersonation-token',
  );
  const { refetch: refetchProjectSettings } = useQuery(getSettings);

  useEffect(() => {
    if (!secretUserImpersonationToken) {
      return;
    }

    void (async () => {
      const { accessToken, refreshToken } =
        await redeemUserImpersonationTokenMutation.mutateAsync({
          secretUserImpersonationToken,
        });

      setAccessToken(accessToken);
      setRefreshToken(refreshToken);

      const { data: getSettingsResponse } = await refetchProjectSettings();
      window.location.href =
        getSettingsResponse!.settings!.afterLoginRedirectUri ||
        getSettingsResponse!.settings!.redirectUri;
    })();
  }, [secretUserImpersonationToken]);

  return <></>;
};
