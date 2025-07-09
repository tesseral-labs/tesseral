import { useMutation, useQuery } from "@connectrpc/connect-query";
import React, { useEffect, useState } from "react";
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
import { Organization } from "@/gen/tesseral/intermediate/v1/intermediate_pb";
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

  const [validOrgs, setValidOrgs] = useState<Organization[] | undefined>();

  useEffect(() => {
    (async () => {
      if (!listOrganizationsResponse) {
        return;
      }

      let organizations = listOrganizationsResponse.organizations;
      if (!projectSettings.selfServeCreateUsers) {
        organizations = organizations.filter((org) => org.userExists);
      }

      setValidOrgs(organizations);

      if (organizations.length === 0) {
        if (!projectSettings.selfServeCreateOrganizations) {
          return;
        }

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
        {validOrgs !== undefined && validOrgs.length === 0 ? (
          <>
            <div className="text-muted-foreground text-sm">
              No organizations match your login credentials.
            </div>

            <Button className="mt-4 w-full" variant="outline" asChild>
              <Link to="/login">Back to Login Page</Link>
            </Button>
          </>
        ) : (
          <>
            <div className="space-y-2">
              {validOrgs?.map((org) => (
                <Button
                  key={org.id}
                  className="w-full"
                  variant="outline"
                  asChild
                >
                  <Link to={`/organizations/${org.id}/login`}>
                    {org.displayName}
                  </Link>
                </Button>
              ))}
            </div>

            {projectSettings.selfServeCreateOrganizations && (
              <>
                <div className="block relative w-full cursor-default my-6">
                  <div className="absolute inset-0 flex items-center border-muted-foreground">
                    <span className="w-full border-t"></span>
                  </div>
                  <div className="relative flex justify-center text-xs uppercase">
                    <span className="bg-card px-2 text-muted-foreground">
                      or
                    </span>
                  </div>
                </div>

                <Button className="w-full" variant="outline" asChild>
                  <Link to="/create-organization">
                    Create a new organization
                  </Link>
                </Button>
              </>
            )}
          </>
        )}
      </CardContent>
    </LoginFlowCard>
  );
}
