import React from "react";
import { Link, useParams } from "react-router";

import { Button } from "@/components/ui/button";
import { DialogFooter } from "@/components/ui/dialog";

import SetupWizardVideo from "../SetupWizardVideo";

export function NameOktaSamlApplication() {
  const { samlConnectionId } = useParams();
  return (
    <>
      <div className="space-y-4 text-sm">
        <SetupWizardVideo src="/videos/saml-setup-wizard/okta/name.gif" />

        <p className="font-medium">Name your Okta SAML application:</p>
        <ol className="list-decimal list-inside space-y-2">
          <li>Give your new Okta application a name.</li>
          <li>Click "Next".</li>
        </ol>
      </div>

      <DialogFooter>
        <Link
          to={`/organization/saml-connections/${samlConnectionId}/setup/okta/configure`}
        >
          <Button size="sm">Continue</Button>
        </Link>
      </DialogFooter>
    </>
  );
}
