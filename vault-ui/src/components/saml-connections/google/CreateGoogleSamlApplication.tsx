import React from "react";
import { Link, useParams } from "react-router";

import { Button } from "@/components/ui/button";
import { DialogFooter } from "@/components/ui/dialog";

import { SetupWizardVideo } from "../SetupWizardVideo";

export function CreateGoogleSamlApplication() {
  const { samlConnectionId } = useParams();
  return (
    <>
      <div className="space-y-4 text-sm">
        <SetupWizardVideo src="/videos/saml-setup-wizard/google/create.gif" />

        <p className="font-medium">Create a new Google SAML application:</p>
        <ol className="list-decimal list-inside space-y-2">
          <li>
            Go to{" "}
            <Link to="https://admin.google.com" target="_blank">
              admin.google.com
            </Link>
            .
          </li>
          <li>In the sidebar, click "Apps &gt; Web and mobile apps"</li>
          <li>Click on the "Add app" dropdown</li>
          <li>Click "Add custom SAML app"</li>
        </ol>
      </div>

      <DialogFooter>
        <Link
          to={`/organization/saml-connections/${samlConnectionId}/setup/google/name`}
        >
          <Button size="sm">Continue</Button>
        </Link>
      </DialogFooter>
    </>
  );
}
