import React, { Dispatch, FC, SetStateAction, useState } from 'react';
import { Title } from '@/components/Title';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import {
  exchangeIntermediateSessionForSession,
  verifyPassword,
} from '@/gen/openauth/intermediate/v1/intermediate-IntermediateService_connectquery';
import { useMutation } from '@connectrpc/connect-query';
import { useNavigate } from 'react-router';
import { setAccessToken, setRefreshToken } from '@/auth';
import { Input } from '@/components/ui/input';
import { AuthType, useAuthType, useIntermediateOrganization } from '@/lib/auth';
import { LoginLayouts, LoginViews } from '@/lib/views';
import { useLayout } from '@/lib/settings';
import { cn } from '@/lib/utils';
import { toast } from 'sonner';
import { parseErrorMessage } from '@/lib/errors';
import Loader from '@/components/ui/loader';

interface VerifyPasswordProps {
  setView: Dispatch<SetStateAction<LoginViews>>;
}

const VerifyPassword: FC<VerifyPasswordProps> = ({ setView }) => {
  const authType = useAuthType();
  const layout = useLayout();
  const organization = useIntermediateOrganization();
  const navigate = useNavigate();

  const [password, setPassword] = useState<string>('');
  const [submitting, setSubmitting] = useState<boolean>(false);

  const exchangeIntermediateSessionForSessionMutation = useMutation(
    exchangeIntermediateSessionForSession,
  );
  const verifyPasswordMutation = useMutation(verifyPassword);

  const deriveNextView = (): LoginViews | undefined => {
    const hasMultipleSecondFactors =
      organization?.userHasAuthenticatorApp && organization?.userHasPasskey;
    const hasSecondFactor =
      organization?.userHasAuthenticatorApp || organization?.userHasPasskey;

    if (organization?.requireMfa) {
      if (hasMultipleSecondFactors || !hasSecondFactor) {
        return LoginViews.ChooseAdditionalFactor;
      } else if (organization?.userHasPasskey) {
        return LoginViews.VerifyPasskey;
      } else if (organization?.userHasAuthenticatorApp) {
        return LoginViews.VerifyAuthenticatorApp;
      }
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitting(true);

    try {
      await verifyPasswordMutation.mutateAsync({
        password,
        organizationId: organization?.id,
      });

      const nextView = deriveNextView();

      if (nextView) {
        setView(nextView);
        return;
      }

      const { accessToken, refreshToken } =
        await exchangeIntermediateSessionForSessionMutation.mutateAsync({
          organizationId: organization?.id,
        });

      setAccessToken(accessToken);
      setRefreshToken(refreshToken);

      navigate('/settings');
      setSubmitting(false);
    } catch (error) {
      setSubmitting(false);
      const message = parseErrorMessage(error);
      toast.error('Could not verify password', {
        description: message,
      });
    }
  };

  return (
    <>
      <Title title="Verify password" />

      <Card
        className={cn(
          'w-full max-w-sm',
          layout !== LoginLayouts.Centered && 'shadow-none border-0',
        )}
      >
        <CardHeader>
          <CardTitle className="text-center">Password Verification</CardTitle>
          <p className="text-sm text-center mt-2 text-gray-500">
            Please enter your password to continue{' '}
            {authType === AuthType.SignUp ? 'signing' : 'logging'} in.
          </p>
        </CardHeader>
        <CardContent className="flex flex-col items-center justify-center w-full">
          <form
            className="flex flex-col items-center w-full"
            onSubmit={handleSubmit}
          >
            <Input
              className="w-full mb-2"
              id="password"
              placeholder="Enter your password"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
            />
            <Button
              className="w-full"
              disabled={password.length < 1 || submitting}
              type="submit"
            >
              {submitting && <Loader />}
              Continue
            </Button>
          </form>
        </CardContent>
      </Card>
    </>
  );
};

export default VerifyPassword;
