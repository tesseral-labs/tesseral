import { useMutation } from '@connectrpc/connect-query';
import { redeemUserImpersonationToken } from '@/gen/tesseral/intermediate/v1/intermediate-IntermediateService_connectquery';
import React, { useEffect } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { setAccessToken, setRefreshToken } from '@/auth';

export const ImpersonatePage = () => {
  const redeemUserImpersonationTokenMutation = useMutation(
    redeemUserImpersonationToken,
  );
  const [searchParams] = useSearchParams();
  const secretUserImpersonationToken = searchParams.get(
    'secret-user-impersonation-token',
  );
  const navigate = useNavigate();

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
      navigate('/settings');
    })();
  }, [
    secretUserImpersonationToken,
    redeemUserImpersonationTokenMutation,
    navigate,
  ]);

  return <></>;
};
