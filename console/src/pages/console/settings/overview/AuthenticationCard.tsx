import { useQuery } from "@connectrpc/connect-query";
import { Shield } from "lucide-react";
import React from "react";
import { Link } from "react-router";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { getProject } from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";

export function AuthenticationCard() {
  const { data: getProjectResponse } = useQuery(getProject);

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <Shield />
          Authentication
        </CardTitle>
        <CardDescription>
          Configure what authentication methods are available to your customers,
          including SAML, OAuth, and Multi-Factor Authentication (MFA).
        </CardDescription>
      </CardHeader>
      <CardContent className="flex-grow">
        <div className="space-y-4">
          <div className="space-y-2">
            <div>
              <div className="font-semibold">
                Enabled authentication methods
              </div>
              <div className="text-xs text-muted-foreground">
                The following authentication methods are configured and can be
                enabled for Organizations in your Project.
              </div>
            </div>
            <div className="flex flex-wrap gap-2">
              {getProjectResponse?.project?.logInWithEmail && (
                <Badge variant="outline">Email Magic Link</Badge>
              )}
              {getProjectResponse?.project?.logInWithPassword && (
                <Badge variant="outline">Password</Badge>
              )}
              {getProjectResponse?.project?.logInWithGoogle && (
                <Badge variant="outline">Google OAuth</Badge>
              )}
              {getProjectResponse?.project?.logInWithMicrosoft && (
                <Badge variant="outline">Microsoft OAuth</Badge>
              )}
              {getProjectResponse?.project?.logInWithGithub && (
                <Badge variant="outline">GitHub OAuth</Badge>
              )}
              {getProjectResponse?.project?.logInWithSaml && (
                <Badge variant="outline">SAML SSO</Badge>
              )}
            </div>
          </div>
          <div className="space-y-2">
            <div>
              <div className="font-semibold">
                Multi-factor authentication (MFA)
              </div>
              <div className="text-xs text-muted-foreground">
                The following multi-factor authentication methods are configured
                and can be enabled for Organizations in your Project.
              </div>
            </div>
            <div className="flex flex-wrap gap-2">
              {getProjectResponse?.project?.logInWithAuthenticatorApp && (
                <Badge variant="outline">Authenticator Apps (TOTP)</Badge>
              )}
              {getProjectResponse?.project?.logInWithPasskey && (
                <Badge variant="outline">Passkeys</Badge>
              )}
              {!getProjectResponse?.project?.logInWithAuthenticatorApp &&
                !getProjectResponse?.project?.logInWithPasskey && (
                  <Badge variant="outline" className="text-muted-foreground">
                    No MFA methods configured
                  </Badge>
                )}
            </div>
          </div>
        </div>
      </CardContent>
      <CardFooter className="mt-4">
        <Link className="w-full" to="/settings/authentication">
          <Button className="w-full" variant="outline">
            Configure Authentication
          </Button>
        </Link>
      </CardFooter>
    </Card>
  );
}
