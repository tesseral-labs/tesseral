import { useMutation, useQuery } from "@connectrpc/connect-query";
import React, { useEffect } from "react";
import { useNavigate } from "react-router";
import { Link } from "react-router";

import { Title } from "@/components/core/Title";
import { LoginFlowCard } from "@/components/login/LoginFlowCard";
import { Button } from "@/components/ui/button";
import { CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  createOrganization,
  listOrganizations,
  whoami,
} from "@/gen/tesseral/intermediate/v1/intermediate-IntermediateService_connectquery";
import { useRedirectNextLoginFlowPage } from "@/hooks/use-redirect-next-login-flow-page";
import { useProjectSettings } from "@/lib/project-settings";

export function ChooseOrganizationPage() {
  const projectSettings = useProjectSettings();
  const { data: listOrganizationsResponse } = useQuery(listOrganizations);
  const { refetch: refetchWhoami } = useQuery(whoami, undefined, {
    enabled: false, // disable by default, the useEffect needs the latest data
  });

  const navigate = useNavigate();
  const { mutateAsync: createOrganizationAsync } =
    useMutation(createOrganization);
  const redirectNextLoginFlowPage = useRedirectNextLoginFlowPage();

  useEffect(() => {
    (async () => {
      if (!listOrganizationsResponse) {
        return;
      }

      if (listOrganizationsResponse.organizations.length === 0) {
        if (projectSettings.autoCreateOrganizations) {
          const { data: whoamiResponse } = await refetchWhoami();

          const baseName =
            whoamiResponse?.intermediateSession?.userDisplayName ||
            whoamiResponse?.intermediateSession?.email.split("@")[0];

          await createOrganizationAsync({
            displayName: `${baseName}'s Organization`,
          });

          redirectNextLoginFlowPage();
          return;
        }

        navigate("/create-organization");
        return;
      }
    })();
  }, [
    listOrganizationsResponse,
    navigate,
    projectSettings,
    createOrganizationAsync,
    redirectNextLoginFlowPage,
    refetchWhoami,
  ]);

  return (
    <LoginFlowCard>
      <Title title="Choose an organization" />
      <CardHeader>
        <CardTitle>Choose an organization</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="space-y-2">
          {listOrganizationsResponse?.organizations?.map((org) => (
            <Button key={org.id} className="w-full" variant="outline" asChild>
              <Link to={`/organizations/${org.id}/login`}>
                {org.displayName}
              </Link>
            </Button>
          ))}
        </div>

        <div className="block relative w-full cursor-default my-6">
          <div className="absolute inset-0 flex items-center border-muted-foreground">
            <span className="w-full border-t"></span>
          </div>
          <div className="relative flex justify-center text-xs uppercase">
            <span className="bg-card px-2 text-muted-foreground">or</span>
          </div>
        </div>

        <Button className="w-full" variant="outline" asChild>
          <Link to="/create-organization">Create a new organization</Link>
        </Button>
      </CardContent>
    </LoginFlowCard>
  );
}
