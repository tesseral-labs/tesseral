import { useQuery } from "@connectrpc/connect-query";
import { Key } from "lucide-react";
import React from "react";

import { Badge } from "@/components/ui/badge";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { getProject } from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";

import { ConfigureGithubOAuthButton } from "./oauth/ConfigureGithubOauthButton";
import { ConfigureGoogleOAuthButton } from "./oauth/ConfigureGoogleOauthButton";
import { ConfigureMicrosoftOAuthButton } from "./oauth/ConfigureMicrosoftOauthButton";

export function OAuthSettingsCard() {
  const { data: getProjectResponse } = useQuery(getProject);

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <Key />
          OAuth Providers
        </CardTitle>
        <CardDescription>
          Configure which OAuth providers are to Organizations within this
          Project.
        </CardDescription>
      </CardHeader>
      <CardContent className="flex-grow">
        <div className="space-y-4 w-full">
          <div className="w-full flex justify-between gap-4">
            <div className="font-semibold text-sm space-x-6">
              <span>Google</span>
              {getProjectResponse?.project?.logInWithGoogle ? (
                <Badge>Enabled</Badge>
              ) : (
                <Badge variant="secondary">Disabled</Badge>
              )}
            </div>
            <ConfigureGoogleOAuthButton />
          </div>
          <div className="w-full flex justify-between gap-4">
            <div className="font-semibold text-sm space-x-6">
              <span>Microsoft</span>
              {getProjectResponse?.project?.logInWithMicrosoft ? (
                <Badge>Enabled</Badge>
              ) : (
                <Badge variant="secondary">Disabled</Badge>
              )}
            </div>
            <ConfigureMicrosoftOAuthButton />
          </div>
          <div className="w-full flex justify-between gap-4">
            <div className="font-semibold text-sm space-x-6">
              <span>GitHub</span>
              {getProjectResponse?.project?.logInWithGithub ? (
                <Badge>Enabled</Badge>
              ) : (
                <Badge variant="secondary">Disabled</Badge>
              )}
            </div>
            <ConfigureGithubOAuthButton />
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
