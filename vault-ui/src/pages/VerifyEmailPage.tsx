import React, { useEffect } from 'react';
import { Title } from '@/components/Title';
import { useNavigate, useSearchParams } from 'react-router-dom';
import {
  verifyEmailChallenge,
  whoami,
} from '@/gen/tesseral/intermediate/v1/intermediate-IntermediateService_connectquery';
import { useMutation, useQuery } from '@connectrpc/connect-query';
import { LoginViews } from '@/lib/views';
import { parseErrorMessage } from '@/lib/errors';
import { toast } from 'sonner';

export const VerifyEmailPage = () => {
  const navigate = useNavigate();

  const [searchParams] = useSearchParams();

  const verifyEmailChallengeMutation = useMutation(verifyEmailChallenge);
  const whoamiQuery = useQuery(whoami);

  useEffect(() => {
    void (async () => {
      const code = searchParams.get('code');

      if (code) {
        try {
          // Redeem the email verification code.
          await verifyEmailChallengeMutation.mutateAsync({
            code,
          });

          // Fetch the whoami query to determine if the user has verified their email.
          const { data: whoamiRes } = await whoamiQuery.refetch();
          if (!whoamiRes) {
            throw new Error('No data returned from whoami query');
          }

          navigate(`/login?view=${LoginViews.ChooseOrganization}`);
        } catch (error) {
          const message = parseErrorMessage(error);
          toast.error('Failed to verify email', {
            description: message,
          });
        }
      }
    })();
  }, []);

  return (
    <div className="flex flex-col items-center justify-center h-screen">
      <Title title="Verifying Email..." />
    </div>
  );
};
