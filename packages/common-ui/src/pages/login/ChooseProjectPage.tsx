import React, { FC, useEffect, useState } from "react";
import { useNavigate } from "react-router";
import {
  exchangeIntermediateSessionForSession,
  listOrganizations,
  setOrganization,
} from "../../gen/tesseral/intermediate/v1/intermediate-IntermediateService_connectquery";
import { useMutation, useQuery } from "@connectrpc/connect-query";
import { Organization } from "../../gen/tesseral/intermediate/v1/intermediate_pb";
import { Title } from "../../components/Title";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "../../components/ui/card";
import { Button } from "../../components/ui/button";
import TextDivider from "../../components/ui/text-divider";
import { parseErrorMessage } from "../../lib/errors";
import { toast } from "sonner";
import Loader from "../../components/ui/loader";

interface ChooseProjectPageProps {
  setAccessToken: (accessToken: string) => void;
  setRefreshToken: (refreshToken: string) => void;
}

export const ChooseProjectPage: FC<ChooseProjectPageProps> = ({
  setAccessToken,
  setRefreshToken,
}) => {
  const navigate = useNavigate();

  const [submitting, setSubmitting] = useState<boolean>(false);

  const { data: listOrganizationsResponse } = useQuery(listOrganizations);
  const exchangeIntermediateSessionForSessionMutation = useMutation(
    exchangeIntermediateSessionForSession
  );
  const setOrganizationMutation = useMutation(setOrganization);

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
  const deriveNextPage = (organization: Organization): string | undefined => {
    if (organization.logInWithPassword) {
      if (organization.userHasPassword) {
        return "/verify-password";
      }

      return "/register-password";
    } else if (
      organization.userHasAuthenticatorApp &&
      !organization.userHasPasskey
    ) {
      return "/verify-authenticator-app";
    } else if (
      organization.userHasPasskey &&
      !organization.userHasAuthenticatorApp
    ) {
      return "/verify-passkey";
    } else if (organization.requireMfa) {
      // this is the case where the organization has multiple secondary factors and requires mfa
      return "/authenticate-another-way";
    }
  };

  const handleOrganizationClick = async (organization: Organization) => {
    setSubmitting(true);
    try {
      await setOrganizationMutation.mutateAsync({
        organizationId: organization.id,
      });

      const nextPage = deriveNextPage(organization);

      if (nextPage) {
        navigate(nextPage);
        return;
      }

      const { accessToken, refreshToken } =
        await exchangeIntermediateSessionForSessionMutation.mutateAsync({});

      setAccessToken(accessToken);
      setRefreshToken(refreshToken);

      setSubmitting(false);
      navigate("/settings");
    } catch (error) {
      setSubmitting(false);
      const message = parseErrorMessage(error);
      toast.error("Could not select Project", {
        description: message,
      });
    }
  };

  useEffect(() => {
    if (listOrganizationsResponse?.organizations?.length === 0) {
      navigate("/create-project");
    }
  }, [listOrganizationsResponse]);

  return (
    <>
      <Title title="Choose an Project" />

      <Card className="w-full max-w-sm">
        <CardHeader>
          <CardTitle>
            {(listOrganizationsResponse?.organizations?.length || 0) > 0
              ? "Choose a Project"
              : "Create a Project"}
          </CardTitle>
        </CardHeader>
        <CardContent className="flex flex-col items-center justify-center w-full">
          {(listOrganizationsResponse?.organizations?.length || 0) > 0 && (
            <ul className="w-full p-0">
              {listOrganizationsResponse?.organizations?.map((organization) => (
                <li key={organization.id}>
                  <Button
                    className="w-full"
                    disabled={submitting}
                    onClick={() => handleOrganizationClick(organization)}
                    variant="outline"
                  >
                    {submitting && <Loader />}
                    {organization.displayName}
                  </Button>
                </li>
              ))}
            </ul>
          )}
          <div className="w-full">
            {(listOrganizationsResponse?.organizations?.length || 0) > 0 && (
              <TextDivider className="w-full">Or you can</TextDivider>
            )}

            <Button
              className="w-full"
              onClick={() => navigate("/create-project")}
            >
              Create a Project
            </Button>
          </div>
        </CardContent>
      </Card>
    </>
  );
};
