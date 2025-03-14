import { useQuery } from "@connectrpc/connect-query";
import { useCallback } from "react";
import { useNavigate } from "react-router";

import {
  listOrganizations,
  whoami,
} from "@/gen/tesseral/intermediate/v1/intermediate-IntermediateService_connectquery";
import {
  IntermediateSession,
  PrimaryAuthFactor,
} from "@/gen/tesseral/intermediate/v1/intermediate_pb";
import { Organization } from "@/gen/tesseral/intermediate/v1/intermediate_pb";

export function useRedirectNextLoginFlowPage(): () => void {
  const { refetch: refetchWhoami } = useQuery(whoami);
  const { refetch: refetchListOrganizations } = useQuery(listOrganizations);
  const navigate = useNavigate();

  return useCallback(async () => {
    // get the latest state of the intermediate session and its organizations
    const { data: listOrganizationsResponse } =
      await refetchListOrganizations();
    const { data: whoamiResponse } = await refetchWhoami();

    const intermediateSession = whoamiResponse!.intermediateSession!;
    const organization = listOrganizationsResponse!.organizations.find(
      (org) => org.id === intermediateSession.organizationId,
    )!;

    // authenticate another way if our current primary auth factor is
    // unacceptable
    if (
      !isPrimaryAuthFactorAcceptable(
        whoamiResponse!.intermediateSession!,
        organization,
      )
    ) {
      navigate(`/authenticate-another-way`);
      return;
    }

    console.log("asdf1");

    // verify password if there is one registered, and it's not already verified
    if (organization.userHasPassword && !intermediateSession.passwordVerified) {
      navigate(`/verify-password`);
      return;
    }

    console.log("asdf2");

    // verify a secondary factor if there is one registered but not verified
    const needsVerifyAuthenticatorApp =
      organization.userHasAuthenticatorApp &&
      !intermediateSession.authenticatorAppVerified;
    const needsVerifyPasskey =
      organization.userHasPasskey &&
      !intermediateSession.authenticatorAppVerified;
    if (needsVerifyAuthenticatorApp || needsVerifyPasskey) {
      navigate(`/verify-secondary-factor`);
      return;
    }

    console.log("asdf3");

    // register a password if the org uses them and the user doesn't have one
    // registered
    if (organization.logInWithPassword && !organization.userHasPassword) {
      navigate(`/register-password`);
      return;
    }

    console.log("asdf4");

    // register a secondary factor if the org requires their registration (i.e.
    // requires MFA) and the user doesn't have one registered
    if (organization.requireMfa && !(organization.userHasPasskey || organization.userHasAuthenticatorApp)) {
      navigate(`/register-secondary-factor`);
      return;
    }

    console.log("asdf5");

    // we have everything we need, finish login flow
    navigate(`/finish-login`);

    console.log("asdf6");
  }, [navigate, refetchListOrganizations, refetchWhoami]);
}

function isPrimaryAuthFactorAcceptable(
  intermediateSession: IntermediateSession,
  organization: Organization,
): boolean {
  switch (intermediateSession.primaryAuthFactor) {
    case PrimaryAuthFactor.EMAIL:
      return organization.logInWithEmail;
    case PrimaryAuthFactor.GOOGLE:
      return organization.logInWithGoogle;
    case PrimaryAuthFactor.MICROSOFT:
      return organization.logInWithMicrosoft;
    default:
      return false;
  }
}
