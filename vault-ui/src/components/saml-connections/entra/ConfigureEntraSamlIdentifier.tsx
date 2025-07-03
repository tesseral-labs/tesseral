import { useQuery } from "@connectrpc/connect-query";
import React from "react";
import { Link, useParams } from "react-router";

import { ValueCopier } from "@/components/core/ValueCopier";
import { Button } from "@/components/ui/button";
import { DialogFooter } from "@/components/ui/dialog";
import { getSAMLConnection } from "@/gen/tesseral/frontend/v1/frontend-FrontendService_connectquery";

export function ConfigureEntraSamlIdentifier() {
  const { samlConnectionId } = useParams();
  const { data: getSamlConnectionResponse } = useQuery(getSAMLConnection, {
    id: samlConnectionId,
  });

  const samlConnection = getSamlConnectionResponse?.samlConnection;

  return (
    <>
      <div className="space-y-4 text-sm">
        <img
          className="rounded-xl max-w-full border shadow-md"
          src="/videos/saml-setup-wizard/entra/identifier.gif"
        />

        <p className="font-medium">Configure SAML Identifier (Entity ID)</p>
        <ol className="list-decimal list-inside space-y-2">
          <li>Navigate to your application if you haven't already.</li>
          <li>In the sidebar for the application, click "Single sign-on"</li>
          <li>Click on "SAML"</li>
          <li>
            Click on "Edit" icon to the right of "Basic SAML Configuration"
          </li>
          <li>
            Click "Add identifier". An input now appears under the "Identifier
            (Entity ID)" section.
          </li>
          <li>
            Enter the following value:
            <ValueCopier
              value={samlConnection?.spEntityId || ""}
              label="Entity ID"
              maxLength={50}
            />
          </li>
        </ol>
      </div>

      <DialogFooter>
        <Link
          to={`/organization/saml-connections/${samlConnectionId}/setup/entra/reply-url`}
        >
          <Button type="button" size="sm">
            Continue
          </Button>
        </Link>
      </DialogFooter>
    </>
  );
}
