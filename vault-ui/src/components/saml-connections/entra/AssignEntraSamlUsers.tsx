import React from "react";
import { Link, useParams } from "react-router";

import { Button } from "@/components/ui/button";
import { DialogFooter } from "@/components/ui/dialog";

import { SetupWizardVideo } from "../SetupWizardVideo";

export function AssignEntraSamlUsers() {
  const { samlConnectionId } = useParams();

  return (
    <>
      <div className="space-y-4 text-sm">
        <SetupWizardVideo src="/videos/saml-setup-wizard/entra/users.gif" />

        <p className="font-medium">Assign users to the new app.</p>
        <p>
          If you're familiar with Entra application user assignments, use
          whatever process you normally use.
        </p>
        <p>Otherwise, the most straightforward process is to:</p>
        <ol className="list-decimal list-inside space-y-2">
          <li>Click on "Users and groups" in the application sidebar</li>
          <li>Click on "Add user/group"</li>
          <li>Under "Users", click on "None Selected"</li>
          <li>
            Check the checkbox next to each of the users you want to assign to
            the application. If you intend to test the application yourself,
            remember to include yourself
          </li>
          <li>Click "Select" at the bottom</li>
          <li>Click "Assign"</li>
        </ol>
      </div>
      <DialogFooter>
        <Link to={`/organization/saml-connections/${samlConnectionId}/test`}>
          <Button type="button" size="sm">
            Continue
          </Button>
        </Link>
      </DialogFooter>
    </>
  );
}
