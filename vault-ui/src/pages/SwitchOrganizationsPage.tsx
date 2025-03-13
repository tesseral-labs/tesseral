import React, { useEffect } from 'react';
import { LoginViews } from '@/lib/views';
import { isValidPrimaryAuthFactor } from '@/lib/auth-factors';
import { parseErrorMessage } from '@/lib/errors';
import { toast } from 'sonner';
import { useIntermediateExchangeAndRedirect } from '@/hooks/use-intermediate-exchange-and-redirect';
import { Organization } from '@/gen/tesseral/intermediate/v1/intermediate_pb';
import { useMutation, useQuery } from '@connectrpc/connect-query';
import {
  exchangeSessionForIntermediateSession,
  setOrganization,
  whoami,
} from '@/gen/tesseral/intermediate/v1/intermediate-IntermediateService_connectquery';
import { useParams } from 'react-router';

export const SwitchOrganizationsPage = () => {
  // largely copied from ChooseOrganization
  const intermediateExchangeAndRedirect = useIntermediateExchangeAndRedirect();
  const { organizationId } = useParams();
  const exchangeSessionForIntermediateSessionMutation = useMutation(
    exchangeSessionForIntermediateSession,
  );
  const setOrganizationMutation = useMutation(setOrganization);
  useEffect(() => {
    void (async () => {
      await exchangeSessionForIntermediateSessionMutation.mutateAsync({});
      await setOrganizationMutation.mutateAsync({
        organizationId,
      });

      // This code is wrong -- it doesn't handle additional login steps -- but
      // fixing that would be more difficult than just rewriting the login flow
      // frontend anyway.
      //
      // To that end, just immediately do the exchange and redirect.
      intermediateExchangeAndRedirect();
    })();
  }, [organizationId]);

  return <h1>switch organizations</h1>;
};
