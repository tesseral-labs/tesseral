import {
  exchangeIntermediateSessionForSession,
  getSettings,
} from '@/gen/tesseral/intermediate/v1/intermediate-IntermediateService_connectquery';
import { useMutation, useQuery } from '@connectrpc/connect-query';
import { useCallback } from 'react';

export function useIntermediateExchangeAndRedirect(): () => void {
  const { refetch } = useQuery(getSettings);
  const exchangeIntermediateSessionForSessionMutation = useMutation(
    exchangeIntermediateSessionForSession,
  );

  return useCallback(async () => {
    const { newUser } =
      await exchangeIntermediateSessionForSessionMutation.mutateAsync({});
    const { data: getSettingsResponse } = await refetch();

    const preferredRedirect = newUser
      ? getSettingsResponse!.settings!.afterSignupRedirectUri
      : getSettingsResponse!.settings!.afterLoginRedirectUri;

    window.location.href =
      preferredRedirect ?? getSettingsResponse!.settings!.redirectUri;
  }, []);
}
