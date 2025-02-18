import React, { Dispatch, FC, SetStateAction } from 'react';
import OAuthButton, { OAuthMethods } from '@/components/OAuthButton';
import { Title } from '@/components/Title';
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import { AuthType, useAuthType, useIntermediateOrganization } from '@/lib/auth';
import useSettings, { useLayout } from '@/lib/settings';
import { cn } from '@/lib/utils';
import { LoginLayouts, LoginViews } from '@/lib/views';
import { parseErrorMessage } from '@/lib/errors';
import { toast } from 'sonner';
import { useMutation } from '@connectrpc/connect-query';
import {
  getGoogleOAuthRedirectURL,
  getMicrosoftOAuthRedirectURL,
} from '@/gen/openauth/intermediate/v1/intermediate-IntermediateService_connectquery';
import TextDivider from '@/components/ui/text-divider';
import EmailForm from '@/components/EmailForm';
import { Button } from '@/components/ui/button';

interface ChooseOrganizationPrimaryLoginFactorProps {
  setView: Dispatch<SetStateAction<LoginViews>>;
}

const ChooseOrganizationPrimaryLoginFactor: FC<
  ChooseOrganizationPrimaryLoginFactorProps
> = ({ setView }) => {
  const authType = useAuthType();
  const layout = useLayout();
  const organization = useIntermediateOrganization();
  const settings = useSettings();

  const googleOAuthRedirectUrlMutation = useMutation(getGoogleOAuthRedirectURL);
  const microsoftOAuthRedirectUrlMutation = useMutation(
    getMicrosoftOAuthRedirectURL,
  );

  const handleGoogleOAuthLogin = async (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();

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
      <Title title="Choose a different login method" />

      <Card
        className={cn(
          'w-full max-w-sm',
          layout !== LoginLayouts.Centered && 'shadow-none border-0',
        )}
      >
        <CardHeader>
          <CardTitle>
            Choose a different{' '}
            {authType === AuthType.SignUp ? 'sign up' : 'log in'} method
          </CardTitle>
          <CardDescription>
            <b>{organization?.displayName ?? 'This organization'}</b> only
            supports the following{' '}
            {authType === AuthType.SignUp ? 'sign up' : 'log in'} methods.
            Please select one below.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div
            className={cn(
              'w-full grid gap-6',
              organization?.logInWithGoogle && organization?.logInWithMicrosoft
                ? 'grid-cols-2'
                : 'grid-cols-1',
            )}
          >
            {settings?.logInWithGoogle && organization?.logInWithGoogle && (
              <OAuthButton
                method={OAuthMethods.google}
                onClick={handleGoogleOAuthLogin}
                variant="outline"
              />
            )}
            {settings?.logInWithMicrosoft &&
              organization?.logInWithMicrosoft && (
                <OAuthButton
                  method={OAuthMethods.microsoft}
                  onClick={handleMicrosoftOAuthLogin}
                  variant="outline"
                />
              )}
          </div>

          {((settings?.logInWithGoogle && organization?.logInWithGoogle) ||
            (settings?.logInWithMicrosoft &&
              organization?.logInWithMicrosoft)) &&
            settings?.logInWithEmail &&
            organization?.logInWithEmail && (
              <TextDivider
                variant={layout !== LoginLayouts.Centered ? 'wider' : 'wide'}
              >
                or continue with email
              </TextDivider>
            )}

          {settings?.logInWithEmail && organization?.logInWithEmail && (
            <EmailForm
              skipIntermediateSessionCreation
              skipListSAMLOrganizations
              setView={setView}
            />
          )}

          {settings?.logInWithSaml &&
            organization?.logInWithSaml &&
            organization?.primarySamlConnectionId && (
              <>
                <TextDivider>or continue with SAML</TextDivider>

                <div className="flex flex-col items-center">
                  <a
                    href={`/api/saml/v1/${organization.primarySamlConnectionId}/init`}
                    className="w-full"
                  >
                    <Button variant="outline">Continue with SAML</Button>
                  </a>
                </div>
              </>
            )}
        </CardContent>
      </Card>
    </>
  );
};

export default ChooseOrganizationPrimaryLoginFactor;
