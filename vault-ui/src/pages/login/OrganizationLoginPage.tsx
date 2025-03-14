import { useMutation, useQuery } from "@connectrpc/connect-query";
import { LoaderCircleIcon } from "lucide-react";
import React, { useEffect } from "react";
import { useNavigate, useParams } from "react-router";

import {
  listOrganizations,
  setOrganization,
  whoami,
} from "@/gen/tesseral/intermediate/v1/intermediate-IntermediateService_connectquery";
import {
  IntermediateSession,
  Organization,
  PrimaryAuthFactor,
} from "@/gen/tesseral/intermediate/v1/intermediate_pb";

export function OrganizationLoginPage() {
  const { organizationId } = useParams();
  const { refetch: refetchWhoami } = useQuery(whoami);
  const { refetch: refetchListOrganizations } = useQuery(listOrganizations);
  const { mutateAsync: setOrganizationMutationAsync } =
    useMutation(setOrganization);
  const navigate = useNavigate();

  useEffect(() => {
    (async () => {
      const { data: listOrganizationsResponse } =
        await refetchListOrganizations();

      const { data: whoamiResponse } = await refetchWhoami();

      await setOrganizationMutationAsync({
        organizationId,
      });

      const organization = listOrganizationsResponse!.organizations.find(
        (org) => org.id === organizationId,
      )!;

      if (!isPrimaryAuthFactorAcceptable(whoamiResponse!.intermediateSession!, organization)) {
        navigate(`/authenticate-another-way`)
      }

      if (organization.userHasPassword) {
        navigate(`/verify-password`)
      }

      if (organization.userHasAuthenticatorApp || organization.userHasPasskey) {
        navigate(`/verify-secondary-factor`)
      }

      if (organization.logInWithPassword && !organization.userHasPassword) {
        navigate(`/register-password`)
      }

      if (organization.requireMfa) {
        navigate(`/register-secondary-factor`)
      }
    })();
  }, [navigate, organizationId, refetchListOrganizations, refetchWhoami, setOrganizationMutationAsync]);

  return (
    <LoaderCircleIcon className="mx-auto text-muted-foreground h-4 w-4 animate-spin" />
  );
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
