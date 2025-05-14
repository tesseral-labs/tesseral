import { useMutation, useQuery } from '@connectrpc/connect-query';
import { useCallback } from 'react';
import { useNavigate } from 'react-router';

import {
  issueEmailVerificationChallenge,
  listOrganizations,
  whoami,
} from '@/gen/tesseral/intermediate/v1/intermediate-IntermediateService_connectquery';
import {
  IntermediateSession,
  PrimaryAuthFactor,
} from '@/gen/tesseral/intermediate/v1/intermediate_pb';
import { Organization } from '@/gen/tesseral/intermediate/v1/intermediate_pb';

export function useRedirectNextLoginFlowPage(): () => void {
  // don't eagerly fetch; we won't use their initial values
  const { refetch: refetchWhoami } = useQuery(whoami, undefined, {
    enabled: false,
  });
  const { refetch: refetchListOrganizations } = useQuery(
    listOrganizations,
    undefined,
    {
      enabled: false,
    },
  );
  const { mutateAsync: issueEmailVerificationChallengeMutationAsync } =
    useMutation(issueEmailVerificationChallenge);

  const navigate = useNavigate();

  return useCallback(async () => {
    const { data: whoamiResponse } = await refetchWhoami();
    const intermediateSession = whoamiResponse!.intermediateSession!;

    if (!intermediateSession.emailVerified) {
      await issueEmailVerificationChallengeMutationAsync({
        email: whoamiResponse?.intermediateSession?.email,
      });

      navigate(`/verify-email`);
      return;
    }

    if (!intermediateSession.organizationId) {
      navigate(`/choose-organization`);
      return;
    }

    // get the latest state of the intermediate session and its organizations
    const { data: listOrganizationsResponse } =
      await refetchListOrganizations();

    const organization = listOrganizationsResponse!.organizations.find(
      (org) => org.id === intermediateSession.organizationId,
    )!;

    console.log(
      'primaryAuthFactor',
      whoamiResponse?.intermediateSession?.primaryAuthFactor,
    );

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

    // verify password if there is one registered, and it's not already verified
    if (
      organization.logInWithPassword &&
      organization.userHasPassword &&
      !intermediateSession.passwordVerified
    ) {
      navigate(`/verify-password`);
      return;
    }

    // Check if we may need to verify a secondary factor. If a user has a
    // secondary factor registered, they must always verify it.
    //
    // We can skip verifying secondary factors if one is already verified.
    if (
      !intermediateSession.passkeyVerified &&
      !intermediateSession.authenticatorAppVerified
    ) {
      // if a user has both a passkey and an authenticator app, let them choose
      // which one to verify
      if (organization.userHasPasskey && organization.userHasAuthenticatorApp) {
        navigate(`/verify-secondary-factor`);
        return;
      }

      if (organization.userHasPasskey) {
        navigate(`/verify-passkey`);
        return;
      }

      if (organization.userHasAuthenticatorApp) {
        navigate(`/verify-authenticator-app`);
        return;
      }

      // falling through here because the user does not have a registered
      // secondary factor
    }

    // register a password if the org uses them and the user + intermediate
    // session doesn't have one registered
    if (
      organization.logInWithPassword &&
      !organization.userHasPassword &&
      !intermediateSession.passwordVerified
    ) {
      navigate(`/register-password`);
      return;
    }

    // Check for needing to register a secondary factor. Users only need to
    // register these if the organization requires MFA.
    //
    // We can also skip registering secondary factors if one is already
    // verified.
    if (
      organization.requireMfa &&
      !intermediateSession.passkeyVerified &&
      !intermediateSession.authenticatorAppVerified
    ) {
      // if only one of passkey or authenticator app is configured, redirect
      // directly to that setup page
      if (
        organization.logInWithPasskey &&
        !organization.logInWithAuthenticatorApp
      ) {
        navigate(`/register-passkey`);
        return;
      }

      if (
        organization.logInWithAuthenticatorApp &&
        !organization.logInWithPasskey
      ) {
        navigate(`/register-authenticator-app`);
        return;
      }

      // let the user choose which one to register
      navigate(`/register-secondary-factor`);
      return;
    }

    // we have everything we need, finish login flow
    navigate(`/finish-login`);
  }, [
    issueEmailVerificationChallengeMutationAsync,
    navigate,
    refetchListOrganizations,
    refetchWhoami,
  ]);
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
    case PrimaryAuthFactor.GITHUB:
      return organization.logInWithGithub;
    default:
      return false;
  }
}
