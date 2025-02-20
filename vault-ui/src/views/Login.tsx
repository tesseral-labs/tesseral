import React, { Dispatch, FC, SetStateAction } from 'react';

import EmailForm from '@/components/EmailForm';
import OAuthButton, { OAuthMethods } from '@/components/OAuthButton';
import { Title } from '@/components/Title';
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import TextDivider from '@/components/ui/text-divider';
import { useMutation } from '@connectrpc/connect-query';

import {
  createIntermediateSession,
  getGoogleOAuthRedirectURL,
  getMicrosoftOAuthRedirectURL,
} from '@/gen/tesseral/intermediate/v1/intermediate-IntermediateService_connectquery';
import { LoginLayouts, LoginViews } from '@/lib/views';
import useSettings, { useLayout } from '@/lib/settings';
import { cn } from '@/lib/utils';
import { parseErrorMessage } from '@/lib/errors';
import { toast } from 'sonner';
import { AuthType, useAuthType } from '@/lib/auth';
import { Link } from 'react-router-dom';

interface LoginProps {
  setView: Dispatch<SetStateAction<LoginViews>>;
}

const Login: FC<LoginProps> = ({ setView }) => {
  const authType = useAuthType();
  const layout = useLayout();
  const settings = useSettings();

  const createIntermediateSessionMutation = useMutation(
    createIntermediateSession,
  );
  const googleOAuthRedirectUrlMutation = useMutation(getGoogleOAuthRedirectURL);
  const microsoftOAuthRedirectUrlMutation = useMutation(
    getMicrosoftOAuthRedirectURL,
  );

  const handleGoogleOAuthLogin = async (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();

    try {
      // this sets a cookie that subsequent requests use
      await createIntermediateSessionMutation.mutateAsync({});
    } catch (error) {
      const message = parseErrorMessage(error);

      toast.error('Could not initialize new session', {
        description: message,
      });
    }

    try {
      const { url } = await googleOAuthRedirectUrlMutation.mutateAsync({
        redirectUrl: `${window.location.origin}/google-oauth-callback`,
      });

      window.location.href = url;
    } catch (error) {
      const message = parseErrorMessage(error);

      toast.error('Could not get Google OAuth redirect URL', {
        description: message,
      });
    }
  };

  const handleMicrosoftOAuthLogin = async (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();

    try {
      // this sets a cookie that subsequent requests use
      await createIntermediateSessionMutation.mutateAsync({});
    } catch (error) {
      const message = parseErrorMessage(error);

      toast.error('Could not initialize new session', {
        description: message,
      });
    }

    try {
      const { url } = await microsoftOAuthRedirectUrlMutation.mutateAsync({
        redirectUrl: `${window.location.origin}/microsoft-oauth-callback`,
      });

      window.location.href = url;
    } catch (error) {
      const message = parseErrorMessage(error);

      toast.error('Could not get Microsoft OAuth redirect URL', {
        description: message,
      });
    }
  };

  return (
    <>
      <Title title={authType === AuthType.SignUp ? 'Sign up' : 'Log in'} />

      <Card
        className={cn(
          'w-full max-w-sm',
          layout !== LoginLayouts.Centered && 'shadow-none border-0',
        )}
      >
        <CardHeader>
          {(settings?.logInWithGoogle || settings?.logInWithMicrosoft) && (
            <CardTitle className="text-center">
              {authType === AuthType.SignUp ? 'Sign up' : 'Log in'}
            </CardTitle>
          )}
        </CardHeader>

        <CardContent className="flex flex-col items-center w-full">
          <div className="w-full grid grid-cols-1 gap-4">
            {settings?.logInWithGoogle && (
              <OAuthButton
                className="w-full"
                method={OAuthMethods.google}
                onClick={handleGoogleOAuthLogin}
                variant="outline"
              />
            )}
            {settings?.logInWithMicrosoft && (
              <OAuthButton
                className="w-full"
                method={OAuthMethods.microsoft}
                onClick={handleMicrosoftOAuthLogin}
                variant="outline"
              />
            )}
          </div>

          {(settings?.logInWithGoogle || settings?.logInWithMicrosoft) &&
            settings?.logInWithEmail && (
              <TextDivider
                variant={layout !== LoginLayouts.Centered ? 'wider' : 'wide'}
              >
                or continue with email
              </TextDivider>
            )}

          {(settings?.logInWithEmail || settings?.logInWithSaml) && (
            <EmailForm
              skipListSAMLOrganizations={!settings?.logInWithSaml}
              setView={setView}
            />
          )}
        </CardContent>
        <CardFooter>
          <div className="text-sm text-center text-muted-foreground w-full">
            {authType === AuthType.SignUp ? (
              <>
                Already have an account?{' '}
                <Link className="underline" to="/login">
                  Log in
                </Link>
              </>
            ) : (
              <>
                Don't have an account?{' '}
                <Link className="underline" to="/signup">
                  Sign up
                </Link>
              </>
            )}
          </div>
        </CardFooter>
      </Card>
    </>
  );
};

export default Login;
