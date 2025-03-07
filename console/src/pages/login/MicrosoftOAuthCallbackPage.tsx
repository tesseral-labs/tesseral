import React, { useEffect } from 'react';
import { Title } from '@/components/Title';
import { useNavigate, useSearchParams } from 'react-router-dom';
import {
  issueEmailVerificationChallenge,
  redeemMicrosoftOAuthCode,
  whoami,
} from '@/gen/tesseral/intermediate/v1/intermediate-IntermediateService_connectquery';
import { useMutation, useQuery } from '@connectrpc/connect-query';
import { LoginView } from '@/lib/views';
import { parseErrorMessage } from '@/lib/errors';
import { toast } from 'sonner';

const MicrosoftOAuthCallbackPage = () => {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();

  const issueEmailVerificationChallengeMutation = useMutation(
    issueEmailVerificationChallenge,
  );
  const redeemMicrosoftOAuthCodeMutation = useMutation(
    redeemMicrosoftOAuthCode,
  );
  const whoamiQuery = useQuery(whoami);

  useEffect(() => {
    void (async () => {
      const code = searchParams.get('code');
      const state = searchParams.get('state');

      if (code && state) {
        try {
          await redeemMicrosoftOAuthCodeMutation.mutateAsync({
            code,
            state,
            redirectUrl: `${window.location.origin}/microsoft-oauth-callback`,
          });

          const { data } = await whoamiQuery.refetch();
          if (!data) {
            throw new Error('No data returned from whoami query');
          }

          if (data?.intermediateSession?.emailVerified) {
            navigate(`/login?view=${LoginView.ChooseProject}`);
            return;
          }

          await issueEmailVerificationChallengeMutation.mutateAsync({
            email: data.intermediateSession!.email,
          });

          navigate(`/login?view=${LoginView.VerifyEmail}`);
        } catch (error) {
          const message = parseErrorMessage(error);
          toast.error('Failed to verify Microsoft log in', {
            description: message,
          });
        }
      }
    })();
  }, []);

  return (
    <div className="flex flex-col items-center justify-center h-screen">
      <Title title="Verifying Microsoft OAuth Credentials..." />

      <div className="space-y-4 text-center"></div>
    </div>
  );
};

export default MicrosoftOAuthCallbackPage;
