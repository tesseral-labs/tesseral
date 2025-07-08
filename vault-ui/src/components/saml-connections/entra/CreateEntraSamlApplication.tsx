import React from "react";
import { Link, useParams } from "react-router";

import { Button } from "@/components/ui/button";
import { DialogFooter } from "@/components/ui/dialog";

import SetupWizardVideo from "../SetupWizardVideo";

export function CreateEntraSamlApplication() {
  const { samlConnectionId } = useParams();

  return (
    <>
      <div className="space-y-4 text-sm">
        <SetupWizardVideo src="/videos/saml-setup-wizard/entra/create.gif" />
        <p className="font-medium">Create a new Entra SAML application:</p>
        <ol className="list-decimal list-inside space-y-2">
          <li>
            Go to{" "}
            <Link to="https://entra.microsoft.com">entra.microsoft.com</Link>.
          </li>
          <li>
            In the sidebar, click "Applications &gt; Enterprise Applications"
          </li>
          <li>Click on the "New application"</li>
          <li>Click "Create your own application"</li>
          <li>Enter a name into "What's the name of your app?"</li>
          <li>
            Keep "Integrate any other application you don't find in the gallery
            (Non-gallery)" checked.
          </li>
          <li>Click "Create"</li>
        </ol>
      </div>

      <DialogFooter>
        <Link
          to={`/organization/saml-connections/${samlConnectionId}/setup/entra/identifier`}
        >
          <Button type="button" size="sm">
            Continue
          </Button>
        </Link>
      </DialogFooter>
    </>
  );
}
