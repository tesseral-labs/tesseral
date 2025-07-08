import { useQuery } from "@connectrpc/connect-query";
import { LoaderCircle } from "lucide-react";
import React from "react";
import { Link, useParams } from "react-router";

import { ValueCopier } from "@/components/core/ValueCopier";
import { Button } from "@/components/ui/button";
import { DialogFooter } from "@/components/ui/dialog";
import { Label } from "@/components/ui/label";
import { getSAMLConnection } from "@/gen/tesseral/frontend/v1/frontend-FrontendService_connectquery";

import { SetupWizardVideo } from "../SetupWizardVideo";

export function ConfigureGoogleSamlApplication() {
  const { samlConnectionId } = useParams();

  const { data: getSamlConnectionResponse, isLoading } = useQuery(
    getSAMLConnection,
    {
      id: samlConnectionId,
    },
  );

  const samlConnection = getSamlConnectionResponse?.samlConnection;

  return (
    <>
      <div className="text-sm space-y-4">
        <SetupWizardVideo src="/videos/saml-setup-wizard/google/configure.gif" />

        <p className="font-medium">Configure your Google SAML application:</p>
        <ol className="list-decimal pl-6">
          <li>
            In the Google Admin console, go to "Apps &gt; Web and mobile apps".
          </li>
          <li>Click on your SAML app, then click on "Configure SAML".</li>
          <li>
            Enter the required information from your SAML connection settings.
          </li>
        </ol>
      </div>

      {isLoading ? (
        <div className="flex items-center justify-center">
          <LoaderCircle className="animate-spin text-muted-foreground/50" />
        </div>
      ) : (
        <div className="space-y-6 w-full">
          {samlConnection && (
            <>
              <div className="space-y-2 w-full">
                <Label>Assertion Consumer Service (ACS) URL</Label>
                <ValueCopier
                  value={samlConnection.spAcsUrl}
                  label="ACS URL"
                  maxLength={50}
                />
              </div>
              <div className="space-y-2 w-full">
                <Label>SP Entity ID</Label>
                <ValueCopier
                  value={samlConnection.spEntityId}
                  label="SP Entity ID"
                  maxLength={50}
                />
              </div>
            </>
          )}
        </div>
      )}

      <DialogFooter>
        <Link
          to={`/organization/saml-connections/${samlConnectionId}/setup/google/users`}
        >
          <Button size="sm">Continue</Button>
        </Link>
      </DialogFooter>
    </>
  );
}
