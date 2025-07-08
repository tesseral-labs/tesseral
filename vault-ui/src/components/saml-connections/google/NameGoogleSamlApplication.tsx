import React from "react";
import { Link, useParams } from "react-router";

import { Button } from "@/components/ui/button";
import { DialogFooter } from "@/components/ui/dialog";

import SetupWizardVideo from "../SetupWizardVideo";

export function NameGoogleSamlApplication() {
  const { samlConnectionId } = useParams();
  return (
    <>
      <div className="space-y-4 text-sm">
        <SetupWizardVideo src="/videos/saml-setup-wizard/google/name.gif" />

        <p className="font-medium">Name your Google SAML application:</p>
        <ol className="list-decimal list-inside space-y-2">
          <li>Give the new Google application a name.</li>
          <li>Optionally, provide a description and upload a logo.</li>
          <li>Click "Continue".</li>
        </ol>
      </div>

      <DialogFooter>
        <Link
          to={`/organization/saml-connections/${samlConnectionId}/setup/google/metadata`}
        >
          <Button size="sm">Continue</Button>
        </Link>
      </DialogFooter>
    </>
  );
}
