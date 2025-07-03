import { useQuery } from "@connectrpc/connect-query";
import React from "react";
import { Link, useParams } from "react-router";

import { ValueCopier } from "@/components/core/ValueCopier";
import { Button } from "@/components/ui/button";
import { DialogFooter } from "@/components/ui/dialog";
import { getSAMLConnection } from "@/gen/tesseral/frontend/v1/frontend-FrontendService_connectquery";

export function ConfigureOktaSamlApplication() {
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
          src="/videos/saml-setup-wizard/okta/configure.gif"
        />

        <p className="font-medium">Create your Okta SAML application:</p>
        <ol className="list-decimal list-inside space-y-2">
          <li>
            Set the "Single sign-on URL" to:
            <ValueCopier
              value={samlConnection?.spAcsUrl || ""}
              label="Single sign-on URL"
              maxLength={50}
            />
          </li>
          <li>
            Make sure "Use this for Recipient URL and Destination URL" stays
            checked.
          </li>
          <li>
            Set the "Audience URI (SP Entity ID)" to:
            <ValueCopier
              value={samlConnection?.spEntityId || ""}
              label="Audience URI"
              maxLength={50}
            />
          </li>
          <li>Click "Next"</li>
          <li>Select "This is an internal app that we have created".</li>
          <li>Click "Finish".</li>
        </ol>
      </div>

      <DialogFooter>
        <Link
          to={`/organization/saml-connections/${samlConnectionId}/setup/okta/metadata`}
        >
          <Button size="sm">Continue</Button>
        </Link>
      </DialogFooter>
    </>
  );
}
