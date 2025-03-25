import { useMutation, useQuery } from "@connectrpc/connect-query";
import React from "react";

import { GoogleIcon } from "@/components/login/GoogleIcon";
import { LoginFlowCard } from "@/components/login/LoginFlowCard";
import { MicrosoftIcon } from "@/components/login/MicrosoftIcon";
import { Button } from "@/components/ui/button";
import {
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  getGoogleOAuthRedirectURL,
  getMicrosoftOAuthRedirectURL,
  listOrganizations,
  whoami,
} from "@/gen/tesseral/intermediate/v1/intermediate-IntermediateService_connectquery";

export function AuthenticateAnotherWayPage() {
  const { data: whoamiResponse } = useQuery(whoami);
  const { data: listOrganizationsResponse } = useQuery(listOrganizations);

  const organization = listOrganizationsResponse?.organizations?.find(
    (org) => org.id === whoamiResponse?.intermediateSession?.organizationId,
  );

  const { mutateAsync: getGoogleOAuthRedirectURLAsync } = useMutation(
    getGoogleOAuthRedirectURL,
  );

  async function handleLogInWithGoogle() {
    const { url } = await getGoogleOAuthRedirectURLAsync({
      redirectUrl: `${window.location.origin}/google-oauth-callback`,
    });
    window.location.href = url;
  }

  const { mutateAsync: getMicrosoftOAuthRedirectURLAsync } = useMutation(
    getMicrosoftOAuthRedirectURL,
  );

  async function handleLogInWithMicrosoft() {
    const { url } = await getMicrosoftOAuthRedirectURLAsync({
      redirectUrl: `${window.location.origin}/microsoft-oauth-callback`,
    });
    window.location.href = url;
  }

  return (
    <LoginFlowCard>
      <CardHeader>
        <CardTitle>Authenticate another way</CardTitle>
        <CardDescription>
          To continue logging in, you must choose from one of the login methods
          below.
        </CardDescription>
      </CardHeader>
      <CardContent>
        <div className="space-y-2">
          {organization?.logInWithGoogle && (
            <Button
              className="w-full"
              variant="outline"
              onClick={handleLogInWithGoogle}
            >
              <GoogleIcon />
              Log in with Google
            </Button>
          )}
          {organization?.logInWithMicrosoft && (
            <Button
              className="w-full"
              variant="outline"
              onClick={handleLogInWithMicrosoft}
            >
              <MicrosoftIcon />
              Log in with Microsoft
            </Button>
          )}
        </div>
      </CardContent>
    </LoginFlowCard>
  );
}
