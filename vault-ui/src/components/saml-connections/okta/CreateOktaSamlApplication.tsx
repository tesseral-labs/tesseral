import React from "react";
import { Link, useParams } from "react-router";

import { Button } from "@/components/ui/button";
import { DialogFooter } from "@/components/ui/dialog";

export function CreateOktaSamlApplication() {
  const { samlConnectionId } = useParams();
  return (
    <>
      <div className="space-y-4 text-sm">
        <img
          className="rounded-xl max-w-full border shadow-md"
          src="/videos/saml-setup-wizard/okta/create.gif"
        />

        <p className="font-medium">Create your Okta SAML application:</p>
        <ol className="list-decimal list-inside space-y-2">
          <li>Go to Applications &gt; Applications in the sidebar.</li>
          <li>Click "Create App Integration"</li>
          <li>Choose "SAML 2.0"</li>
          <li>Click "Next"</li>
        </ol>
      </div>

      <DialogFooter>
        <Link
          to={`/organization/saml-connections/${samlConnectionId}/setup/okta/name`}
        >
          <Button size="sm">Continue</Button>
        </Link>
      </DialogFooter>
    </>
  );
}
