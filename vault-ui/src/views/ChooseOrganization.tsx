import React, { Dispatch, FC, SetStateAction, useState } from 'react';

import { Title } from '@/components/Title';
import { useMutation, useQuery } from '@connectrpc/connect-query';
import {
  listOrganizations,
  setOrganization,
  whoami,
} from '@/gen/tesseral/intermediate/v1/intermediate-IntermediateService_connectquery';
import { Organization } from '@/gen/tesseral/intermediate/v1/intermediate_pb';
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import { useNavigate } from 'react-router-dom';
import { setAccessToken, setRefreshToken } from '@/auth';
import { LoginLayouts, LoginViews } from '@/lib/views';
import { cn } from '@/lib/utils';
import { useLayout } from '@/lib/settings';
import { parseErrorMessage } from '@/lib/errors';
import { toast } from 'sonner';
import Loader from '@/components/ui/loader';
import {
  isValidPrimaryLoginFactor,
  PrimaryLoginFactor,
} from '@/lib/login-factors';
import TextDivider from '@/components/ui/text-divider';
import { Button } from '@/components/ui/button';
import {
  useIntermediateExchangeAndRedirect
} from '@/hooks/use-intermediate-exchange-and-redirect';

interface ChooseOrganizationProps {
  setIntermediateOrganization: Dispatch<
    SetStateAction<Organization | undefined>
  >;
  setView: Dispatch<SetStateAction<LoginViews>>;
}

const ChooseOrganization: FC<ChooseOrganizationProps> = ({
  setIntermediateOrganization,
  setView,
}) => {
  const layout = useLayout();
  const navigate = useNavigate();

  const [setting, setSetting] = useState<boolean>(false);

  const { data: whoamiRes } = useQuery(whoami);
  const { data: listOrganizationsResponse } = useQuery(listOrganizations);
  const setOrganizationMutation = useMutation(setOrganization);
  const intermediateExchangeAndRedirect = useIntermediateExchangeAndRedirect()

  // This function is effectively the central routing layer for Organization selection.
  // It determines the next view based on the current organization's settings and the user's current state.
  // If the organization requires MFA
  // - if the primary login factor is not valid, the user is redirected to the ChooseOrganizationPrimaryLoginFactor view
  // - if the organization has passwords enabled
  //   - if the user has a password, they are redirected to the VerifyPassword view
  //   - if the user does not have a password, they are redirected to the RegisterPassword view
  // - if the organization's only secondary factor is authenticator apps
  //   - if the user has an authenticator app registered, they are redirected to the VerifyAuthenticatorApp view
  //   - if the user does not have an authenticator app registered, they are redirected to the RegisterAuthenticatorApp view
  // - if the organization's only secondary factor is passkeys
  //   - if the user has a passkey registered, they are redirected to the VerifyPasskey view
  //   - if the user does not have a passkey registered, they are redirected to the RegisterPasskey view
  // - if the organization has multiple secondary factors
  //   - the user is redirected to the ChooseAdditionalFactor view
  const deriveNextView = (
    organization: Organization,
  ): LoginViews | undefined => {
    const primaryLoginFactor =
      whoamiRes?.intermediateSession?.primaryLoginFactor;

    if (
      primaryLoginFactor &&
      !isValidPrimaryLoginFactor(
        primaryLoginFactor as PrimaryLoginFactor,
        organization,
      )
    ) {
      return LoginViews.ChooseOrganizationPrimaryLoginFactor;
    } else {
      if (organization.logInWithPassword) {
        if (organization.userHasPassword) {
          return LoginViews.VerifyPassword;
        }

        return LoginViews.RegisterPassword;
      } else if (
        organization.userHasAuthenticatorApp &&
        !organization.userHasPasskey
      ) {
        return LoginViews.VerifyAuthenticatorApp;
      } else if (
        organization.userHasPasskey &&
        !organization.userHasAuthenticatorApp
      ) {
        return LoginViews.VerifyPasskey;
      } else if (organization.requireMfa) {
        // this is the case where the organization has multiple secondary factors and requires mfa
        return LoginViews.ChooseAdditionalFactor;
      }
    }
  };

  const handleOrganizationClick = async (organization: Organization) => {
    setSetting(true);
    const intermediateOrganization = {
      ...organization,
    };
    try {
      if (!whoamiRes?.intermediateSession) {
        throw new Error('No intermediate session found');
      }

      await setOrganizationMutation.mutateAsync({
        organizationId: organization.id,
      });

      setSetting(false);
      setIntermediateOrganization(intermediateOrganization);
    } catch (error) {
      setSetting(false);
      const message = parseErrorMessage(error);
      toast.error('Could not set organization', {
        description: message,
      });
    }

    try {
      // Check if the needs to provide additional factors
      const nextView = deriveNextView(intermediateOrganization);
      if (nextView) {
        setView(nextView);
        return;
      }

      intermediateExchangeAndRedirect();
    } catch (error) {
      setSetting(false);
      const message = parseErrorMessage(error);
      toast.error('Could not set organization', {
        description: message,
      });
    }
  };

  return (
    <>
      <Title title="Choose an Organization" />

      <Card
        className={cn(
          'w-full max-w-sm',
          layout !== LoginLayouts.Centered && 'shadow-none border-0',
        )}
      >
        <CardHeader>
          <CardTitle className="text-center">Choose an Organization</CardTitle>
        </CardHeader>
        <CardContent className="flex flex-col items-center justify-center w-full">
          <ul className="w-full p-0">
            {listOrganizationsResponse?.organizations?.map((organization) => (
              <li
                key={organization.id}
                onClick={() => handleOrganizationClick(organization)}
              >
                <Button
                  className="w-full"
                  variant="outline"
                  onClick={() => handleOrganizationClick(organization)}
                >
                  {setting && <Loader />}
                  {organization.displayName}
                </Button>
              </li>
            ))}
          </ul>
        </CardContent>
        <CardFooter>
          <div className="flex flex-col w-full">
            <TextDivider className="w-full">Or you can</TextDivider>

            <Button
              className="mt-4 w-full"
              onClick={() => setView(LoginViews.CreateOrganization)}
            >
              Create an organization
            </Button>
          </div>
        </CardFooter>
      </Card>
    </>
  );
};

export default ChooseOrganization;
